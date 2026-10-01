package rewards

import (
	"math"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

func Evaluate(c catalog.Card, p spending.Profile) int64 {
	var earned float64
	annualUsed := make([]int64, len(c.Groups))

	for m := range p.Months {
		monthlyUsed := make([]int64, len(c.Groups))

		for cat := range p.Months[m] {
			remaining := p.Months[m][cat]

			for i, g := range c.Groups {
				r, ok := g.Rates[spending.Category(cat)]
				if !ok {
					continue
				}
				if g.Cap == nil {
					earned += float64(remaining) * r
					remaining = 0
					break
				}
				if g.Cap.Period == catalog.Monthly {
					available := g.Cap.AmountCents - monthlyUsed[i]
					bonusSpend := remaining
					if bonusSpend > available {
						bonusSpend = available
					}
					earned += float64(bonusSpend) * r
					monthlyUsed[i] += bonusSpend
					remaining -= bonusSpend
					if remaining == 0 {
						break
					}
				}
				if g.Cap.Period == catalog.Annual {
					available := g.Cap.AmountCents - annualUsed[i]
					bonusSpend := remaining
					if bonusSpend > available {
						bonusSpend = available
					}
					earned += float64(bonusSpend) * r
					annualUsed[i] += bonusSpend
					remaining -= bonusSpend
					if remaining == 0 {
						break
					}

				}

			}
			earned += float64(remaining) * c.BaseRate

		}

	}
	return int64(math.Round(earned)) - c.AnnualFeeCents

}
