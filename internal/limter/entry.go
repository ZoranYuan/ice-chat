package limter

import (
	"sync"

	"golang.org/x/time/rate"
)

type Limter struct {
	mu      sync.Mutex
	burst   int
	limters map[uint64]*rate.Limiter
	rate    rate.Limit
}

func NewLimter(r rate.Limit, brust int) *Limter {
	return &Limter{
		burst:   brust,
		rate:    r,
		limters: make(map[uint64]*rate.Limiter),
	}
}

func (l *Limter) Allow(uid uint64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.limters[uid]
	if !ok {
		limter := rate.NewLimiter(l.rate, l.burst)
		l.limters[uid] = limter
	}

	return limiter.Allow()
}
