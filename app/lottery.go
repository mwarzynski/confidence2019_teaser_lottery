package app

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Lottery struct {
	accounts map[string]Account
	winners  map[string]struct{}

	mutex sync.RWMutex
}

func NewLottery(ctx context.Context, period time.Duration) *Lottery {
	l := &Lottery{
		accounts: make(map[string]Account),
		winners:  make(map[string]struct{}),
	}
	go func() {
		ticker := time.NewTicker(period)
		for {
			select {
			case <-ticker.C:
				l.evaluate()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return l
}

func (l *Lottery) Add(account Account) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.accounts[account.Name] = account
}

func (l *Lottery) IsWinner(name string) bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if _, won := l.winners[name]; won {
		return true
	}
	return false
}

func (l *Lottery) Winners() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	var ws []string
	for w := range l.winners {
		ws = append(ws, w)
	}
	return ws
}

func (l *Lottery) evaluate() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	accounts := l.accounts
	l.winners = make(map[string]struct{})
	l.accounts = make(map[string]Account)
	for name, account := range accounts {
		amounts := append(account.Amounts, randInt(999913, 3700000))
		sum := 0
		for _, a := range amounts {
			sum += a
		}
		if sum == 0x133700 {
			l.winners[name] = struct{}{}
		}
	}
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}
