package rewards

import (
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

func TestEvaluate(t *testing.T){
	card := catalog.Card{BaseRate: 0.015}
	var p spending.Profile

	for m := range p.Months{
		p.Months[m][spending.Groceries] = 40000
	}
	got := Evaluate(card, p)
	want := int64(7200)
	if got != want{
		t.Errorf("Evalulate() = %d cents, want %d ", got, want)
	}
}