package services

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"SotoAyam/config"
	"SotoAyam/internal/dto"
	"SotoAyam/internal/models"
	"SotoAyam/internal/repository"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var (
	ErrInvalidSignature            = errors.New("signature notifikasi midtrans tidak valid")
	ErrTransactionAlreadyProcessed = errors.New("notifikasi dengan status ini sudah pernah diproses")
)

type PaymentService interface {
	CreateSnapTransaction(
		ctx context.Context,
		order *models.Order,
	) (*dto.CreateSnapTokenResponse, error)

	HandleNotification(
		ctx context.Context,
		payload dto.MidtransNotificationPayload,
	) error
}

type paymentService struct {
	orderRepository   repository.OrderRepository
	paymentRepository repository.PaymentRepository
	snapClient        snap.Client
	serverKey         string
}

func NewPaymentService(
	orderRepository repository.OrderRepository,
	paymentRepository repository.PaymentRepository,
	cfg config.MidtransConfig,
) PaymentService {
	env := midtrans.Sandbox
	if cfg.IsProduction {
		env = midtrans.Production
	}

	var snapClient snap.Client
	snapClient.New(cfg.ServerKey, env)

	return &paymentService{
		orderRepository:   orderRepository,
		paymentRepository: paymentRepository,
		snapClient:        snapClient,
		serverKey:         cfg.ServerKey,
	}
}

func (s *paymentService) CreateSnapTransaction(
	ctx context.Context,
	order *models.Order,
) (*dto.CreateSnapTokenResponse, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.OrderCode,
			GrossAmt: int64(order.Total),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.CustomerName,
			Phone: order.CustomerPhone.String,
		},
	}

	snapResp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat transaksi midtrans: %w",
			err,
		)
	}

	if err := s.orderRepository.UpdateSnapToken(
		ctx,
		order.ID,
		snapResp.Token,
	); err != nil {
		return nil, fmt.Errorf(
			"gagal menyimpan snap token: %w",
			err,
		)
	}

	return &dto.CreateSnapTokenResponse{
		SnapToken:   snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (s *paymentService) HandleNotification(
	ctx context.Context,
	payload dto.MidtransNotificationPayload,
) error {
	if !s.isValidSignature(payload) {
		return ErrInvalidSignature
	}

	existing, err := s.paymentRepository.FindLatestByTransactionID(
		ctx,
		payload.TransactionID,
	)
	if err != nil && !errors.Is(err, repository.ErrPaymentTransactionNotFound) {
		return fmt.Errorf(
			"gagal memeriksa transaksi sebelumnya: %w",
			err,
		)
	}

	if existing != nil && existing.TransactionStatus == payload.TransactionStatus {
		return ErrTransactionAlreadyProcessed
	}

	order, err := s.orderRepository.FindByCode(
		ctx,
		payload.OrderID,
	)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		return fmt.Errorf(
			"gagal menemukan order: %w",
			err,
		)
	}

	rawResponse, _ := json.Marshal(payload)
	grossAmount, _ := strconv.ParseFloat(payload.GrossAmount, 64)

	transaction := &models.PaymentTransaction{
		OrderID:           order.ID,
		TransactionID:     payload.TransactionID,
		TransactionStatus: payload.TransactionStatus,
		RawResponse:       rawResponse,
		PaymentType: sql.NullString{
			String: payload.PaymentType,
			Valid:  payload.PaymentType != "",
		},
		GrossAmount: sql.NullFloat64{
			Float64: grossAmount,
			Valid:   true,
		},
		FraudStatus: sql.NullString{
			String: payload.FraudStatus,
			Valid:  payload.FraudStatus != "",
		},
	}

	newPaymentStatus, newOrderStatus := resolvePaymentResult(
		payload.TransactionStatus,
		payload.FraudStatus,
	)

	tx, err := s.orderRepository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)
	}
	defer tx.Rollback(ctx)

	if err := s.paymentRepository.CreateTransaction(
		ctx,
		tx,
		transaction,
	); err != nil {
		return fmt.Errorf(
			"gagal menyimpan payment transaction: %w",
			err,
		)
	}

	if newPaymentStatus != "" {
		if err := s.orderRepository.UpdatePaymentStatus(
			ctx,
			tx,
			order.ID,
			newPaymentStatus,
		); err != nil {
			return fmt.Errorf(
				"gagal memperbarui payment status: %w",
				err,
			)
		}
	}

	if newOrderStatus != "" && newOrderStatus != order.Status {
		if err := s.orderRepository.UpdateStatus(
			ctx,
			tx,
			order.ID,
			newOrderStatus,
		); err != nil {
			return fmt.Errorf(
				"gagal memperbarui status order: %w",
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"gagal commit transaction: %w",
			err,
		)
	}

	return nil
}

func (s *paymentService) isValidSignature(
	payload dto.MidtransNotificationPayload,
) bool {
	raw := payload.OrderID +
		payload.StatusCode +
		payload.GrossAmount +
		s.serverKey

	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	return expected == payload.SignatureKey
}

// resolvePaymentResult memetakan transaction_status dari Midtrans
// ke payment_status dan status order internal.
// String kosong berarti "tidak berubah".
func resolvePaymentResult(
	transactionStatus string,
	fraudStatus string,
) (paymentStatus string, orderStatus string) {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return "paid", "confirmed"
		}
		return "unpaid", ""

	case "settlement":
		return "paid", "confirmed"

	case "pending":
		return "unpaid", ""

	case "deny", "cancel", "expire":
		return "failed", "cancelled"

	case "refund", "partial_refund":
		return "refunded", ""

	default:
		return "", ""
	}
}