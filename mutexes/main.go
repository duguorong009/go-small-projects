package main

import (
	"strconv"
	"sync"
)

var myBalance = &balance{amount: 50.00, currency: "GBP"}

type balance struct {
	amount   float64
	currency string
	mu       sync.RWMutex
}

func (b *balance) Add(i float64) {
	b.mu.Lock()
	// This is racy
	b.amount += i
	b.mu.Unlock()
}

func (b *balance) Display() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return strconv.FormatFloat(b.amount, 'f', 2, 64) + " " + b.currency
}

func main() {

}
