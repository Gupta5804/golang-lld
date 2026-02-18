package main

type AccountSystem interface {
	CheckKYC(user string) error
}

type BankSystem interface {
	Debit(amount int) error
}

type BlockChainSystem interface {
	Credit(amount int) error
}

type NotificationSystem interface {
	SendReceipt(user string, amount int)
}

type WalletFacade struct {
	account      AccountSystem
	bank         BankSystem
	blockchain   BlockChainSystem
	notification NotificationSystem
}

func (w *WalletFacade) BuyCrypto(amount int,user string) error{
	if err := w.account.CheckKYC(user); err != nil{
		return err
	}
	if err := w.bank.Debit(amount); err != nil {
		return err
	}
	if err := w.blockchain.Credit(amount); err != nil{
		return err
	}
	w.notification.SendReceipt(user,amount)
	return nil
}
