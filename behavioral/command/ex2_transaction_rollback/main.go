package main

import "fmt"

// command interface
type UndoableCommand interface {
	Execute() error
	Undo() error
}

// chef/receiver
type Account struct {
	Balance int
}

func (a *Account) Add(amount int) {
	a.Balance += amount
}
func (a *Account) Subtract(amount int) error {
	if a.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	a.Balance -= amount
	return nil
}

// tickets/concrete commands
type DepositCommand struct {
	account *Account
	amount  int
}

func (dc *DepositCommand) Execute() error {
	dc.account.Add(dc.amount)
	return nil
}
func (dc *DepositCommand) Undo() error {
	return dc.account.Subtract(dc.amount)
}

type WithdrawCommand struct {
	account *Account
	amount  int
}

func (wc *WithdrawCommand) Execute() error {
	return wc.account.Subtract(wc.amount)
}
func (wc *WithdrawCommand) Undo() error {
	wc.account.Add(wc.amount)
	return nil
}

// invoker/transaction manager/waiter
type TransactionManager struct {
	history []UndoableCommand
}

func (tm *TransactionManager) ExecuteCommand(uc UndoableCommand) error {
	err := uc.Execute()
	if err == nil {
		tm.history = append(tm.history, uc)
		return nil
	}
	return err
}

func (tm *TransactionManager) UndoLast() error {
	if len(tm.history) == 0 {
		return fmt.Errorf("no history")
	}
	lastIdx := len(tm.history) - 1
	cmd := tm.history[lastIdx]
	tm.history = tm.history[:lastIdx]
	return cmd.Undo()
}
