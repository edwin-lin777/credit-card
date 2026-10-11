package households

import (
	"math"
	"testing"

	"github.com/edwin-lin777/credit-card/internal/spending"
)

// Same inputs should always produce the same household.
func TestGenerateDeterministic(t *testing.T) {
	a := Generate(Canada2023, 42, 1)
	b := Generate(Canada2023, 42, 1)
	c := Generate(Canada2023, 42, 2)

	if a != b {
		t.Fatal("same seed and household ID produced different results")
	}

	if a == c {
		t.Fatal("different household IDs produced identical profiles")
	}
}

// Every category should stay within the configured
// 60%-140% annual spending range.
func TestGenerateSpendingBounds(t *testing.T) {
	for id := int64(0); id < 1000; id++ {
		p := Generate(Canada2023, 42, id)

		for cat := spending.Category(0); cat < spending.NumCategories; cat++ {
			var annual int64

			for month := 0; month < 12; month++ {
				amount := p.Months[month][cat]

				if amount < 0 {
					t.Fatal("negative monthly spending")
				}

				annual += amount
			}

			avg := Canada2023.AnnualAverageCents[cat]
			min := int64(math.Round(float64(avg) * 0.6))
			max := int64(math.Round(float64(avg) * 1.4))

			if annual < min || annual > max {
				t.Fatalf("household %d, category %d: spending outside expected range", id, cat)
			}
		}
	}
}

// Across many generated households, the sample mean
// should be close to the configured spending anchor.
func TestGenerateMeanNearAnchor(t *testing.T) {
	const n = 5000

	var totals [spending.NumCategories]int64

	for id := int64(0); id < n; id++ {
		p := Generate(Canada2023, 42, id)

		for cat := spending.Category(0); cat < spending.NumCategories; cat++ {
			for month := 0; month < 12; month++ {
				totals[cat] += p.Months[month][cat]
			}
		}
	}

	for cat := spending.Category(0); cat < spending.NumCategories; cat++ {
		mean := float64(totals[cat]) / float64(n)
		target := float64(Canada2023.AnnualAverageCents[cat])

		if math.Abs(mean-target)/target > 0.02 {
			t.Fatalf("category %d: mean %.2f differs from anchor %.2f", cat, mean, target)
		}
	}
}
