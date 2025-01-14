package payment

import (
	"crowdfunding/config"
	"crowdfunding/internal/user"
	"log"
	"reflect"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

type service struct{}

func NewService() service {
	return service{}
}

func (svc service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	log.Printf("Server Key : %+v \nbertipe : %v", config.Cfg.Payment.ServerKey, reflect.TypeOf(config.Cfg.Payment.ServerKey))

	// 1. Initiate Snap client
	var client snap.Client
	client.New(config.Cfg.Payment.ServerKey, midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		log.Printf("ERROR PAYMENT : %+v", err)
		return "", err
	}

	return snapResp.RedirectURL, nil

}
