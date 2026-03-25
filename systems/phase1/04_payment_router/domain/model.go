package domain

type PaymentRequest struct {
	TransactionID     string
	UserID            string
	Amount            int64
	Region            string
	ClientIP          string
	DeviceFingerprint string
}

func (p *PaymentRequest) GetIP() string {
	return p.ClientIP
}
func (p *PaymentRequest) GetFingerprint() string {
	return p.DeviceFingerprint
}
func (p *PaymentRequest) GetAmount() int64 {
	return p.Amount
}
func (p *PaymentRequest) GetRegion() string {
	return p.Region
}
