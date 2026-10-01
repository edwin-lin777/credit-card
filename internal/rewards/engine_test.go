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
