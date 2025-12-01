package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	owner   string
	balance float64
}

func NewBankAccount(owner string, balance float64) *BankAccount {
	return &BankAccount{
		owner:   owner,
		balance: balance,
	}
}

func (b *BankAccount) Deposit(amt float64) error {
	if amt <= 0 {
		return errors.New("Depoist must be positive")
	}
	b.balance += amt
	return nil
}

func (b *BankAccount) Withdraw(amt float64) error {
	if amt <= 0 {
		return errors.New("Depoist must be positive")
	}
	if b.balance-amt < 150 {
		return errors.New("insufficient funds")
	}
	b.balance -= amt
	return nil
}

func (b *BankAccount) GetBalance() float64 {
	return b.balance
}

func (b *BankAccount) GetOwner() string {
	return b.owner
}

func main() {
	account := NewBankAccount("Alice", 1000)

	// Deposit
	if err := account.Deposit(500); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Deposited successfully. Balance: %.2f\n", account.balance)
	}

	// Withdraw
	if err := account.Withdraw(200); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Withdrew successfully. Balance: %.2f\n", account.balance)
	}

	// Insufficient funds
	if err := account.Withdraw(2000); err != nil {
		fmt.Println("Error:", err) // Will print error
	}

	fmt.Printf("Final balance for %s: %.2f\n", account.owner, account.balance)

}
