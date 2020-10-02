package wallet

import (
	"github.com/google/uuid"
	"reflect"
	"github.com/Iftikhor99/wallet/v1/pkg/types"
	"testing"
)

// func TestService_Reject_success1(t *testing.T) {
// 	s := &Service{}
// 	phone := types.Phone("+992000000001")
// 	account, err := s.RegisterAccount(phone)
	
// 	if err != nil {
// 		t.Errorf("Reject(): cant't register account, error = %v", err)
// 		return
// 	} 

// 	err = s.Deposit(account.ID, 10_000_00)
// 	if err != nil {
// 		t.Errorf("Reject(): cant't deposit account, error = %v", err)
// 		return
// 	}	

// 	payment, err := s.Pay(account.ID, 1000_00, "auto")
// 	if err != nil {
// 		t.Errorf("Reject(): cant't create payment, error = %v", err)
// 		return
// 	}

// 	err = s.Reject(payment.ID)
// 	if err != nil {
// 		t.Errorf("Reject(): cant't reject payment, error = %v", err)
// 		return
// 	}
// }

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
		
	// Tpo6yem HavtTu nnaTéx
		
	payment := payments[0]
		
	got, err := s.FindPaymentByID(payment. ID)
		
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}
		
	// CpaBHMBaem nnaTexu
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// co3paém cepsuc
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	// TipoOyem HaWTM HeECyWeECTByWuMA nnaTéex
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}
	
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
	
}

func TestService_Reject_success(t *testing.T) {

// co3paém cepsuc
s := newTestService()

_, payments, err := s.addAccount(defaultTestAccount)

if err != nil {
	t.Error(err)
	return
}

// TipoOyem OTMeHMTb nnaTéx

payment := payments[0]

err = s.Reject(payment. ID)

if err != nil {
	t.Errorf("Reject(): error = %v", err)
	return
}

savedPayment, err := s.FindPaymentByID(payment. ID)
if err != nil {
	t.Errorf("Reject(): can't find payment by id, error = %v", err)
	return
}
if savedPayment.Status != types.PaymentStatusFail {
	t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
	return
}

savedAccount, err := s.FindAccountByID(payment.AccountID)

if err != nil {
	t.Errorf("Reject(): can't find account by id, error = %v", err)
	return
}

if savedAccount.Balance != defaultTestAccount.balance {
	t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
	return
}

}

func TestService_Repeat_success(t *testing.T) {

	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
		
	// Tpo6yem HavtTu nnaTéx
		
	payment := payments[0]
		
	err = s.Repeat(payment.ID)
		
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
		
}