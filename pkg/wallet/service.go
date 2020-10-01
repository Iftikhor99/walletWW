package wallet

import (
	"github.com/google/uuid"
	"errors"
	"github.com/Iftikhor99/wallet/v1/pkg/types"
)

type Error string

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}

func (e Error) Error() string {
	return string(e)
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account {
		ID: 	s.nextAccountID,
		Phone: 	phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String() 
	payment := &types.Payment{
		ID:			paymentID,
		AccountID:	accountID,
		Amount:		amount,
		Category: 	category,
		Status:		types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var accountToReturn *types.Account
	var errType = ErrAccountNotFound
	for _, account := range s.accounts {
		if account.ID != accountID {
			accountToReturn = nil			
		}else{
			accountToReturn = account
			errType = nil
		}
	}
	return accountToReturn, errType
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var paymentToReturn *types.Payment
	var errType = ErrAccountNotFound
	for _, payment := range s.payments {
		if payment.ID != paymentID {
			paymentToReturn = nil			
		}else{
			paymentToReturn = payment
			errType = nil
		}
	}
	return paymentToReturn, errType
}

func (s *Service) Reject(paymentID string) error {

	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			pay.Status = types.PaymentStatusFail
			payment = pay
			break
		}
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == payment.AccountID {
			account = acc
			break
		}
	}

	if payment == nil {
		return ErrPaymentNotFound
	}
	
	account.Balance += payment.Amount
	return nil
}