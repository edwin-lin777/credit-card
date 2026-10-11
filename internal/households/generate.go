package households

import (
	"math"
	"math/rand"

	"github.com/edwin-lin777/credit-card/internal/spending"
)

// Generate creates one synthetic household's spending profile.
// The same seed and household ID produce the same results.
func Generate(anchors SpendingAnchors, seed int64, householdID int64) spending.Profile {
	rng := rand.New(rand.NewSource(seed + householdID))

	var profile spending.Profile

	for cat := spending.Category(0); cat < spending.NumCategories; cat++ {
		average := anchors.AnnualAverageCents[cat]

		// Randomize annual spending between 60% and 140%
		// of the category average.
		multiplier := 0.6 + rng.Float64()*0.8
		annual := int64(math.Round(float64(average) * multiplier))

		// Give each month a slightly different spending weight.
		var weights [12]float64
		totalWeight := 0.0

		for month := 0; month < 12; month++ {
			weights[month] = 0.8 + rng.Float64()*0.4
			totalWeight += weights[month]
		}

		// Distribute annual spending across 12 months.
		remaining := annual

		for month := 0; month < 11; month++ {
			amount := int64(float64(annual) * weights[month] / totalWeight)

			profile.Months[month][cat] = amount
			remaining -= amount
		}

		// Assign remaining cents to December.
		profile.Months[11][cat] = remaining
	}

	return profile
}
