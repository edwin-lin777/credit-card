package rewards

import (
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/households"
)

var benchmarkResult int64

func BenchmarkEvaluateWallet(b *testing.B) {
	paths := []string{
		"../../cards/capital-one-quicksilver-world-elite.yaml",
		"../../cards/scotia-momentum.yaml",
		"../../cards/simplii-cash-back.yaml",
	}

	var cards []catalog.Card

	for _, path := range paths {
		card, err := catalog.LoadCard(path)
		if err != nil {
			b.Fatal(err)
		}
		cards = append(cards, card)
	}

	// Generate one household to benchmark.
	profile := households.Generate(households.Canada2023, 42, 1)

	// Prepare the wallet ONCE.
	prepared := PrepareWallet(cards)

	b.ReportAllocs()
	b.ResetTimer()

	var result int64

	for i := 0; i < b.N; i++ {
		result = EvaluatePreparedWallet(prepared, profile)
	}

	benchmarkResult = result
}
