package main

import (
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
)

func TestBuildWallets(t *testing.T) {
	cards := []catalog.Card{
		{ID: "a"},
		{ID: "b"},
		{ID: "c"},
		{ID: "d"},
	}

	wallets := buildWallets(cards)

	// 4 choose 2 = 6
	// 4 choose 3 = 4
	// Total = 10
	if len(wallets) != 10 {
		t.Fatalf("expected 10 wallets, got %d", len(wallets))
	}
	
}

