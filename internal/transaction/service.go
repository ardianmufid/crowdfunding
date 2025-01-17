package transaction

import (
	"context"
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/payment"
	"errors"
	"strconv"
	"time"
)

type Repository interface {
	GetCampaignID(ctx context.Context, CampaignID int) (transactions []Transaction, err error)
	GetByUserID(ctx context.Context, userID int) (transactions []Transaction, err error)
	Save(ctx context.Context, model Transaction) (transaction Transaction, err error)
	Update(ctx context.Context, model Transaction) (transaction Transaction, err error)
	GetByID(ctx context.Context, ID int) (Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repo Repository, campaignRepo campaign.Repository, paymentService payment.Service) service {
	return service{
		repository:         repo,
		campaignRepository: campaignRepo,
		paymentService:     paymentService,
	}
}

func (s service) GetTransactionByCampaignID(ctx context.Context, request GetCampaignTrasactionRequest) (transactions []Transaction, err error) {

	campaign, err := s.campaignRepository.FindCampaignByID(ctx, request.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != request.User.Id {
		return []Transaction{}, errors.New("failed to update campaign")
	}

	transactions, err = s.repository.GetCampaignID(ctx, request.ID)
	if err != nil {
		return []Transaction{}, err
	}

	return
}

func (s service) GetTransactionByUserID(ctx context.Context, userID int) (transactions []Transaction, err error) {

	transactions, err = s.repository.GetByUserID(ctx, userID)
	if err != nil {
		return []Transaction{}, nil
	}

	return
}

func (s service) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (transaction Transaction, err error) {

	modelTransaction := NewFromCreateTransactionRequest(request)

	transaction, err = s.repository.Save(ctx, modelTransaction)
	if err != nil {
		return
	}

	paymentTransaction := payment.Transaction{
		ID:     transaction.ID,
		Amount: transaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, request.User)
	if err != nil {
		return transaction, err
	}

	transaction.PaymentURL = &paymentURL
	transaction.UpdatedAt = time.Now()

	newTransaction, err := s.repository.Update(ctx, transaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s service) ProcessPayment(ctx context.Context, request TransactionNotificationRequest) error {

	transaction_id, _ := strconv.Atoi(request.OrderID)

	transaction, err := s.repository.GetByID(ctx, transaction_id)
	if err != nil {
		return err
	}

	if request.PaymentType == "credit_card" && request.TransactionStatus == "capture" && request.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "deny" || request.TransactionStatus == "expire" || request.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(ctx, transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindCampaignByID(ctx, transaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BeckerCount = campaign.BeckerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(ctx, campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
