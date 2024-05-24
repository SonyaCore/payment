package discounts

import "time"

type Config struct {
	CreditExpiration time.Duration // Duration to expire a discount code
	CodeLength       int           // Length of generated charge or voucher codes
}
