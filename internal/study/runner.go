package study

import (
	"sync"

	"github.com/edwin-lin777/credit-card/internal/households"
	"github.com/edwin-lin777/credit-card/internal/rewards"
)

// Run evaluates all wallets across a number of households
// using a pool of concurrent workers.
func Run(
	wallets []rewards.PreparedWallet,
	anchors households.SpendingAnchors,
	seed int64,
	simulations int,
	workerCount int,
) []int64 {
	totals := make([]int64, len(wallets))

	if simulations <= 0 || len(wallets) == 0 {
		return totals
	}

	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > simulations {
		workerCount = simulations
	}

	jobs := make(chan int64, workerCount*2)
	workerTotals := make([][]int64, workerCount)

	var wg sync.WaitGroup

	// Start workers.
	for i := 0; i < workerCount; i++ {
		workerTotals[i] = make([]int64, len(wallets))

		wg.Add(1)

		go func(localTotals []int64) {
			defer wg.Done()

			for id := range jobs {
				profile := households.Generate(anchors, seed, id)

				for j, wallet := range wallets {
					localTotals[j] += rewards.EvaluatePreparedWallet(wallet, profile)
				}
			}
		}(workerTotals[i])
	}

	// Send household IDs to workers.
	for id := 0; id < simulations; id++ {
		jobs <- int64(id)
	}

	close(jobs)

	// Wait for all workers to finish.
	wg.Wait()

	// Combine each worker's totals.
	for _, localTotals := range workerTotals {
		for i, amount := range localTotals {
			totals[i] += amount
		}
	}

	return totals
}
