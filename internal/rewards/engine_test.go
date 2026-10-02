package rewards

import (
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

func TestEvaluate(t *testing.T) {
	card := catalog.Card{BaseRate: 0.015}
	var p spending.Profile

	for m := range p.Months {
		p.Months[m][spending.Groceries] = 40000
	}
	got := Evaluate(card, p)
	want := int64(7200)
	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d ", got, want)
	}
}
func TestEvaluateBonusCategory(t *testing.T) {
	card := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.04,
				},
			},
		},
	}

	var p spending.Profile
	for m := range p.Months {
		p.Months[m][spending.Groceries] = 40000
		p.Months[m][spending.Dining] = 20000
	}
	got := Evaluate(card, p)
	want := int64(21600)
	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d", got, want)
	}

}
func TestEvaluateSharedCapHigherRateFirst(t *testing.T) {
	card := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.02,
					spending.Bills:     0.04,
				},
				Cap: &catalog.Cap{
					AmountCents: 50000,
					Period:      catalog.Monthly,
				},
			},
		},
	}

	var p spending.Profile
	p.Months[0][spending.Groceries] = 50000
	p.Months[0][spending.Bills] = 50000

	got := Evaluate(card, p)
	want := int64(2500)

	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d", got, want)
	}
}
func TestEvaluateMonthlyCap(t *testing.T) {
	card := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.04,
				},
				Cap: &catalog.Cap{
					AmountCents: 50000,
					Period:      catalog.Monthly,
				},
			},
		},
	}

	var p spending.Profile

	for m := range p.Months {
		p.Months[m][spending.Groceries] = 60000
	}

	got := Evaluate(card, p)
	want := int64(25200)

	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d", got, want)
	}
}
func TestEvaluateSpilloverToLowerRate(t *testing.T) {
	card := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.04,
				},
				Cap: &catalog.Cap{
					AmountCents: 50000,
					Period:      catalog.Monthly,
				},
			},
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.02,
				},
			},
		},
	}

	var p spending.Profile
	p.Months[0][spending.Groceries] = 60000

	got := Evaluate(card, p)
	want := int64(2200)

	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d", got, want)
	}
}
func TestEvaluateAnnualCap(t *testing.T) {
	card := catalog.Card{
		BaseRate: 0.01,
		Groups: []catalog.Group{
			{
				Rates: map[spending.Category]float64{
					spending.Groceries: 0.04,
				},
				Cap: &catalog.Cap{
					AmountCents: 100000,
					Period:      catalog.Annual,
				},
			},
		},
	}

	var p spending.Profile

	for m := range p.Months {
		p.Months[m][spending.Groceries] = 40000
	}

	got := Evaluate(card, p)
	want := int64(7800)

	if got != want {
		t.Errorf("Evaluate() = %d cents, want %d", got, want)
	}
}
