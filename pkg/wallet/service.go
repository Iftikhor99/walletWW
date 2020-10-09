package wallet

import (
	"strings"
	"io"
	"strconv"
	"os"
	"log"
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
	favorites []*types.Favorite
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
//ErrFavoriteNotFound for
var ErrFavoriteNotFound = errors.New("favorite not found")


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
	
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return payment, err
	}

	if account.Balance < payment.Amount {
		return payment, ErrNotEnoughBalance
	}

	account.Balance -= payment.Amount
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

//FavoritePayment for
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	var favoriteToReturn *types.Favorite
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return favoriteToReturn, err
	}
		
	
	favoriteID := uuid.New().String() 
	favorite := &types.Favorite{
		ID:			favoriteID,
		AccountID:	payment.AccountID,
		Name:		name,
		Amount:		payment.Amount,
		Category: 	payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

//PayFromFavorite for
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	
	
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}

	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}

	if favorite.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == favorite.AccountID {
			account = acc
			break
		}
	}
	
	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < favorite.Amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= favorite.Amount
	paymentID := uuid.New().String() 
	payment := &types.Payment{
		ID:			paymentID,
		AccountID:	favorite.AccountID,
		Amount:		favorite.Amount,
		Category: 	favorite.Category,
		Status:		types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}
//ExportToFile for
func (s *Service) ExportToFile(path string) error{

	
	
	fileNew, err := os.Create(path)
	if err != nil {
		log.Print(err)
		
	}

	defer func() {
		
		if cerr := fileNew.Close(); err != nil {
			log.Print(cerr)
		}
	}()
	for index, account := range s.accounts {
	//	account, err = s.FindAccountByID(int64(ind))
	// fmt.Println(newP2)
	// fmt.Println(ee3)
	 if index != 0 {
		_, err = fileNew.Write([]byte("|"))
		if err != nil {
			log.Print(err)
			
		}
		
	}

	_, err = fileNew.Write([]byte(strconv.FormatInt((account.ID), 10)))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(";"))
	if err != nil {
		log.Print(err)
		
	}
	_, err = fileNew.Write([]byte(string(account.Phone)))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(";"))
	if err != nil {
		log.Print(err)
		
	}

	_, err = fileNew.Write([]byte(strconv.FormatInt(int64(account.Balance), 10)))
	if err != nil {
		log.Print(err)
		
	}


	}

	return err

}

//ImportFromFile for
func (s *Service) ImportFromFile(path string) error{

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	log.Printf("%#v", file)

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		content = append(content, buf[:read]...)
	}
	
	data := string(content)
	newData := strings.Split(data, "|")
	log.Print(data)
	log.Print(newData)
	
	
	for ind1, stroka := range newData {
		log.Print(stroka)
		account := &types.Account {
			}	
		newData2 := strings.Split(stroka, ";")
		for ind, stroka2 := range newData2 {
			if stroka2 == "" {
				return ErrPhoneRegistered
			}
			log.Print(stroka2)
			if ind == 0{
				id, _ := strconv.ParseInt(stroka2, 10, 64)
				account.ID = id
			}
			if ind == 1{
				account.Phone = types.Phone(stroka2)
			}
			if ind == 2{
				balance, _ := strconv.ParseInt(stroka2, 10, 64)
				account.Balance = types.Money(balance)
					
			}
					
			// if (ind1 == 0) && (ind ==2) {
				log.Print(ind1)
			// 	s.accounts = append(s.accounts, account)		
			// }  

			// if (ind1 == 1) && (ind ==2) {
			// 	log.Print(account)
			// 	s.accounts = append(s.accounts, account)		
			// } 

		}
		for _, accountCheck := range s.accounts {
			if accountCheck.Phone == account.Phone {
				return ErrPhoneRegistered
			}
			if accountCheck.ID == account.ID {
				return ErrPhoneRegistered
			}
			

		}
		s.accounts = append(s.accounts, account)
	}
	for _, account := range s.accounts {
	//	if account.Phone == phone {
			log.Print(account)
	//	}
	}
	
	return err

}