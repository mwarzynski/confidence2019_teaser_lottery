package app

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidData   = errors.New("invalid data")
)

type Service struct {
	accounts map[string]*Account
	lottery  *Lottery

	mutex sync.RWMutex
}

func NewService(ctx context.Context, lotteryPeriod time.Duration) *Service {
	return &Service{
		accounts: make(map[string]*Account),
		lottery:  NewLottery(ctx, lotteryPeriod),
	}
}

func (s *Service) AccountAdd(account Account) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, found := s.accounts[account.Name]; found {
		return ErrAlreadyExists
	}
	s.accounts[account.Name] = &account
	return nil
}

func (s *Service) AccountAddAmount(name string, amount int) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	account, found := s.accounts[name]
	if !found {
		return ErrNotFound
	}
	if err := account.AddAmount(amount); err != nil {
		return err
	}
	s.accounts[name] = account
	return nil
}

func (s *Service) AccountGet(name string) (Account, bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	account, found := s.accounts[name]
	if !found {
		return Account{}, false, ErrNotFound
	}
	superUser := s.lottery.IsWinner(name) || account.IsMillionaire()
	return *account, superUser, nil
}

func (s *Service) LotteryAdd(name string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	account, found := s.accounts[name]
	if !found {
		return ErrNotFound
	}
	s.lottery.Add(*account)
	return nil
}

func (s *Service) LotteryResults() []string {
	return s.lottery.Winners()
}
