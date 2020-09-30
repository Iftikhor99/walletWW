package types

type Money int64

type PaymentCategory string

type PaymentStatus string

const (
	StatusOk PaymentStatus = "OK"
	StatusFail PaymentStatus = "FAIL"
	StatusInProgress PaymentStatus = "INPROGRESS"
)


type Payment struct {
	ID string 
	AccountID int64
	Amount Money
	Category PaymentCategory
	Status PaymentStatus
}
type Phone string

type Account struct {
    ID int64
    Phone Phone
    Balance Money 
}