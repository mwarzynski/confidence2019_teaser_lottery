package app

import "github.com/pkg/errors"

const (
	MaxAmountsLen = 4
	MaxAmount     = 99
)

type Account struct {
	Name    string `json:"name"`
	Amounts []int  `json:"amounts"`
}

func (a *Account) Validate() error {
	if a.Name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func (a *Account) AddAmount(amount int) error {
	if amount < 0 || amount > MaxAmount {
		return errors.Wrapf(ErrInvalidData, "amount must be positive and less than %d: got '%d'", MaxAmount+1, amount)
	}
	if len(a.Amounts) >= MaxAmountsLen {
		return errors.Wrapf(ErrInvalidData, "reached maximum number of amounts (%d)", MaxAmountsLen)
	}
	a.Amounts = append(a.Amounts, amount)
	return nil
}

func (a *Account) IsMillionaire() bool {
	sum := 0
	for _, a := range a.Amounts {
		sum += a
	}
	return sum >= 1000000
}
