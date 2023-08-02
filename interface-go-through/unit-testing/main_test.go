package main

import (
	"testing"
	"time"
)

type MockShopDB struct{}

func (m *MockShopDB) CountCustomers(_ time.Time) (int, error) {
	return 1000, nil
}

func (m *MockShopDB) CountSales(_ time.Time) (int, error) {
	return 333, nil
}

func TestCalculateSalesRate(t *testing.T) {
	// Initialize the mock.
	m := &MockShopDB{}
	// Pass the mock to the calculateSalesRate() function.
	sr, err := calculateSalesRate(m)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the return value is as expected, based on the mocked
	// inputs.
	exp := "0.33"
	if sr != exp {
		t.Fatalf("got %v; expected %v", sr, exp)
	}
}
