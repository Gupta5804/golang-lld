package main

import "testing"

func TestTransactionManager_Rollback(t *testing.T) {
	acc := &Account{}
	tm := &TransactionManager{}

	cmd1 := &DepositCommand{
		account : acc,
		amount: 100,
	}
	cmd2 := &WithdrawCommand{
		account: acc,
		amount: 40,
	}
	cmd3 := &WithdrawCommand{
		account: acc,
		amount: 100,
	}

	tm.ExecuteCommand(cmd1)
	if acc.Balance != 100 {
		t.Errorf("Expected balance 100, got %d",acc.Balance)
	}
	tm.ExecuteCommand(cmd2)
	if err := tm.ExecuteCommand(cmd3); err == nil{
		t.Error("Expected insufficient account error, got nil")
	}
	tm.UndoLast()
	if acc.Balance != 100{
		t.Errorf("Expected balance 100, got %d",acc.Balance)
	}
	tm.UndoLast()
	if acc.Balance != 0{
		t.Errorf("Expected balance 0, got %d",acc.Balance)
	}
	if err := tm.UndoLast(); err == nil{
		t.Error("Expected 'no history' error, got nil")
	}
}
