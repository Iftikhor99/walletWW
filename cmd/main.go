package main

import (
	"fmt"
	"github.com/Iftikhor99/wallet/v1/pkg/wallet"
)


func main() {
	svc := &wallet.Service{}
	account , err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	} 

	err = svc.Deposit(account.ID, 100_000_00)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")		
		}		
		return
	}
	fmt.Println(account.Balance)

	newP, ee2 := svc.Pay(account.ID,10_000_00,"food")

	
	fmt.Println(account.Balance)
	fmt.Println(newP)
	fmt.Println(ee2)

	newP2, ee3 := svc.FindPaymentByID(newP.ID)
	fmt.Println(newP2)
	fmt.Println(ee3)

	//ee4 := svc.Reject(newP.ID)
	//fmt.Println(account.Balance)
	//fmt.Println(ee4)

	newP3, ee5 := svc.Repeat(newP.ID)
	fmt.Println(ee5)
	fmt.Println(newP3.Amount)
	fmt.Println(account.Balance)


}