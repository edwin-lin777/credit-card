package rewards

import (
	"math"
	"sort"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

func Evaluate(c catalog.Card, p spending.Profile) int64 {
	type opportunity struct {
		group int
		cat   spending.Category
		rate  float64
	}

	var opportunities []opportunity

	for i, g := range c.Groups {
		for cat, rate := range g.Rates {
			if rate <= c.BaseRate {
				continue
			}

			opportunities = append(opportunities, opportunity{
				group: i,
				cat:   cat,
				rate:  rate,
			})
		}
	}

	sort.Slice(opportunities, func(i, j int) bool {
		if opportunities[i].rate != opportunities[j].rate {
			return opportunities[i].rate > opportunities[j].rate
		}

		if opportunities[i].group != opportunities[j].group {
			return opportunities[i].group < opportunities[j].group
		}

		return opportunities[i].cat < opportunities[j].cat
	})

	var earned float64
	annualUsed := make([]int64, len(c.Groups))

	for m := range p.Months {
		monthlyUsed := make([]int64, len(c.Groups))
		month := p.Months[m]

		for _, op := range opportunities {
			remaining := month[op.cat]

			if remaining == 0 {
				continue
			}

			g := c.Groups[op.group]

			if g.Cap == nil {
				earned += float64(remaining) * op.rate
				month[op.cat] = 0
				continue
			}

			if g.Cap.Period == catalog.Monthly {
				available := g.Cap.AmountCents - monthlyUsed[op.group]

				if available <= 0 {
					continue
				}

				bonusSpend := remaining
				if bonusSpend > available {
					bonusSpend = available
				}

				earned += float64(bonusSpend) * op.rate
				monthlyUsed[op.group] += bonusSpend
				month[op.cat] -= bonusSpend
			}

			if g.Cap.Period == catalog.Annual {
				available := g.Cap.AmountCents - annualUsed[op.group]

				if available <= 0 {
					continue
				}

				bonusSpend := remaining
				if bonusSpend > available {
					bonusSpend = available
				}

				earned += float64(bonusSpend) * op.rate
				annualUsed[op.group] += bonusSpend
				month[op.cat] -= bonusSpend
			}
		}

		for cat := range month {
			earned += float64(month[cat]) * c.BaseRate
		}
	}

	return int64(math.Round(earned)) - c.AnnualFeeCents
}
