package wallet

import (
//	"fmt"
	"github.com/google/uuid"
	"errors"
	"github.com/Iftikhor99/wallet/v1/pkg/types"
)

//Error for
type Error string

//Service for
type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}


func (e Error) Error() string {
	return string(e)
}
//ErrPhoneRegistered for
var ErrPhoneRegistered = errors.New("phone already registered")
//ErrAmountMustBePositive for
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
//ErrAccountNotFound for
var ErrAccountNotFound = errors.New("account not found")
//ErrNotEnoughBalance for
var ErrNotEnoughBalance = errors.New("not enough balance")
//ErrPaymentNotFound for
var ErrPaymentNotFound = errors.New("payment not found")

//RegisterAccount for
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

//Deposit for
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

//Pay for
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

//FindAccountByID for
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

//FindPaymentByID for
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var paymentToReturn *types.Payment
	var errType = ErrPaymentNotFound
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

//Reject for
func (s *Service) Reject(paymentID string) error {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail		
	account.Balance += payment.Amount
	return nil
}

//Repeat for
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	// var paymentToReturn = &types.Payment{}
	payment, err := s.FindPaymentByID(paymentID)
	 if err != nil {
	 	return payment, err
	 }
	
	// account, err := s.FindAccountByID(payment.AccountID)
	// if err != nil {
	// 	return paymentToReturn, err
	// }

//	if account.Balance < payment.Amount {
//		return paymentToReturn, ErrNotEnoughBalance
//	}

//	account.Balance -= payment.Amount
	paymentIDNew := uuid.New().String() 
	paymentNew := &types.Payment{
		ID:			paymentIDNew,
		AccountID:	payment.AccountID,
		Amount:		payment.Amount,
		Category: 	payment.Category,
		Status:		types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, paymentNew)
	
	return paymentNew, nil
}
