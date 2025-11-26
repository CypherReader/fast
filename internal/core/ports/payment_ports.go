package ports

// PaymentGateway defines the interface for payment processing
type PaymentGateway interface {
	CreateCustomer(email, name string) (string, error)
	CreateSubscription(customerID, priceID string) (string, error)
	CreatePayout(amount float64, currency, destination string) (string, error)
	ConstructEvent(payload []byte, header string) (interface{}, error)
}
