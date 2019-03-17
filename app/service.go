package app

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidData   = errors.New("invalid data")
)

type Service struct {
	accounts map[string]Account
	lottery  *Lottery

	deleteAccountAfter time.Duration

	mutex sync.RWMutex
}

func NewService(ctx context.Context, lotteryPeriod, deleteAccountAfter time.Duration) *Service {
	return &Service{
		accounts:           make(map[string]Account),
		lottery:            NewLottery(ctx, lotteryPeriod),
		deleteAccountAfter: deleteAccountAfter,
	}
}

func (s *Service) AccountAdd() (Account, error) {
	account := Account{
		Name:    randString(16),
		Amounts: make([]int, 0, 0),
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, found := s.accounts[account.Name]; found {
		return Account{}, ErrAlreadyExists
	}
	s.accounts[account.Name] = account
	go func() {
		<-time.After(s.deleteAccountAfter)
		s.mutex.Lock()
		defer s.mutex.Unlock()
		delete(s.accounts, account.Name)
	}()
	return account, nil
}

func (s *Service) AccountAddAmount(name string, amount int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	return account, superUser, nil
}

func (s *Service) LotteryAdd(name string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	account, found := s.accounts[name]
	if !found {
		return ErrNotFound
	}
	s.lottery.Add(account)
	return nil
}

func (s *Service) LotteryResults() []string {
	return s.lottery.Winners()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
