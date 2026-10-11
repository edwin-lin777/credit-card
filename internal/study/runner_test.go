package study

import (
	"reflect"
	"testing"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/households"
	"github.com/edwin-lin777/credit-card/internal/rewards"
)

func TestRunDeterministic(t *testing.T) {
	cards := []catalog.Card{
		{BaseRate: 0.01},
		{BaseRate: 0.02},
		{BaseRate: 0.015},
	}

	wallets := []rewards.PreparedWallet{
		rewards.PrepareWallet(cards[:2]),
		rewards.PrepareWallet(cards[1:]),
		rewards.PrepareWallet(cards),
	}

	oneWorker := Run(
		wallets,
		households.Canada2023,
		42,
		100,
		1,
	)

	fourWorkers := Run(
		wallets,
		households.Canada2023,
		42,
		100,
		4,
	)

	if !reflect.DeepEqual(oneWorker, fourWorkers) {
		t.Fatalf(
			"worker counts produced different results: %v vs %v",
			oneWorker,
			fourWorkers,
		)
	}
}
