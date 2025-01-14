package transaction

import (
	"context"
	"crowdfunding/internal/campaign"
	"crowdfunding/internal/payment"
	"errors"
	"time"
)

type Repository interface {
	GetCampaignID(ctx context.Context, CampaignID int) (transactions []Transaction, err error)
	GetByUserID(ctx context.Context, userID int) (transactions []Transaction, err error)
	Save(ctx context.Context, model Transaction) (transaction Transaction, err error)
	Update(ctx context.Context, model Transaction) (transaction Transaction, err error)
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
