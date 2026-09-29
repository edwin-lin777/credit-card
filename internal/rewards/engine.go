package rewards

import (
	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
	"math"
)

func Evaluate(c catalog.Card, p spending.Profile) int64 {
	var earned float64
	for m := range p.Months {
		for cat := range p.Months[m] {
			earned += float64(p.Months[m][cat]) * c.BaseRate

		}

	}
	return int64(math.Round(earned)) - c.AnnualFeeCents

}
