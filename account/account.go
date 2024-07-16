package account

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type Account struct {
	ID      int
	Balance float64
	Mu      sync.Mutex
}

func (a *Account) Deposit(amount float64) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	a.Balance += amount
	logOperation(a.ID, "Deposit", amount)
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if amount > a.Balance {
		return errors.New("insufficient funds")
	}
	a.Balance -= amount
	logOperation(a.ID, "Withdraw", amount)
	return nil
}

func (a *Account) GetBalance() float64 {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	logOperation(a.ID, "GetBalance", a.Balance)
	return a.Balance
}

func logOperation(accountID int, operation string, amount float64) {
	fmt.Printf("Account ID: %d, Operation: %s, Amount: %.2f, Time: %s\n", accountID, operation, amount, time.Now().Format(time.RFC3339))
}
