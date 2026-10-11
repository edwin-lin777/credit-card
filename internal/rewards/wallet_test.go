package rewards

import (
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

func TestWalletSplitsMonthlyCap(t *testing.T) {
	bonusCard := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.05,
				},
				Cap: &catalog.Cap{
					AmountCents: 50000,
					Period:      catalog.Monthly,
				},
			},
		},
	}

	baseCard := catalog.Card{BaseRate: 0.015}

	var p spending.Profile
	p.Months[0][spending.Groceries] = 80000

	got := EvaluateWallet([]catalog.Card{bonusCard, baseCard}, p)

	// $500 at 5% = $25
	// $300 at 1.5% = $4.50
	// Total = $29.50
	if got != 2950 {
		t.Fatalf("expected 2950 cents, got %d", got)
	}
}

func TestWalletSharedMonthlyCap(t *testing.T) {
	bonusCard := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.04,
					spending.Dining:    0.05,
				},
				Cap: &catalog.Cap{
					AmountCents: 50000,
					Period:      catalog.Monthly,
				},
			},
		},
	}

	baseCard := catalog.Card{BaseRate: 0.015}

	var p spending.Profile
	p.Months[0][spending.Groceries] = 40000
	p.Months[0][spending.Dining] = 40000

	got := EvaluateWallet([]catalog.Card{bonusCard, baseCard}, p)

	// Dining: $400 at 5% = $20
	// Groceries: $100 at 4% = $4
	// Remaining groceries: $300 at 1.5% = $4.50
	// Total = $28.50
	if got != 2850 {
		t.Fatalf("expected 2850 cents, got %d", got)
	}
}

func TestWalletAnnualCap(t *testing.T) {
	bonusCard := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.05,
				},
				Cap: &catalog.Cap{
					AmountCents: 60000,
					Period:      catalog.Annual,
				},
			},
		},
	}

	baseCard := catalog.Card{BaseRate: 0.015}

	var p spending.Profile
	p.Months[0][spending.Groceries] = 50000
	p.Months[1][spending.Groceries] = 50000

	got := EvaluateWallet([]catalog.Card{bonusCard, baseCard}, p)

	// January: $500 at 5% = $25
	// February: $100 at 5% = $5
	// February: $400 at 1.5% = $6
	// Total = $36
	if got != 3600 {
		t.Fatalf("expected 3600 cents, got %d", got)
	}

}
func TestWalletChargesBothFees(t *testing.T) {
	cards := []catalog.Card{
		{AnnualFeeCents: 12000},
		{AnnualFeeCents: 5000},
	}

	got := EvaluateWallet(cards, spending.Profile{})

	if got != -17000 {
		t.Fatalf("expected -17000 cents, got %d", got)
	}
}
