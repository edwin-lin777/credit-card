package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/households"
	"github.com/edwin-lin777/credit-card/internal/rewards"
	"github.com/edwin-lin777/credit-card/internal/study"
)

// Results for one wallet.
type walletResult struct {
	Cards                  []string `json:"cards"`
	AverageNetRewardsCents int64    `json:"average_net_rewards_cents"`
	TotalNetRewardsCents   int64    `json:"-"`
}

// Complete simulation results.
type studyResult struct {
	Households     int            `json:"households"`
	Workers        int            `json:"workers"`
	Seed           int64          `json:"seed"`
	CardCount      int            `json:"card_count"`
	WalletCount    int            `json:"wallet_count"`
	ElapsedSeconds float64        `json:"elapsed_seconds"`
	Wallets        []walletResult `json:"wallets"`
}

func main() {
	// Read command-line arguments.
	householdCount := flag.Int("households", 1000, "Number of households")
	workers := flag.Int("workers", runtime.GOMAXPROCS(0), "Number of workers")
	seed := flag.Int64("seed", 42, "Random seed")
	output := flag.String("output", "results.json", "Output JSON file")

	flag.Parse()

	if *householdCount <= 0 {
		log.Fatal("households must be greater than zero")
	}

	if *workers <= 0 {
		log.Fatal("workers must be greater than zero")
	}

	// Load all card YAML files.
	paths, err := filepath.Glob("cards/*.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cards []catalog.Card

	for _, path := range paths {
		if filepath.Base(path) == "test-bonus.yaml" {
			continue
		}

		card, err := catalog.LoadCard(path)
		if err != nil {
			log.Fatal(err)
		}

		cards = append(cards, card)
	}

	// Generate all two-card and three-card wallets.
	wallets := buildWallets(cards)

	if len(wallets) == 0 {
		log.Fatal("not enough cards to build wallets")
	}

	// Precompute each wallet's sorted reward offers.
	preparedWallets := make([]rewards.PreparedWallet, len(wallets))

	for i, wallet := range wallets {
		preparedWallets[i] = rewards.PrepareWallet(wallet)
	}

	// Run concurrent Monte Carlo simulation.
	start := time.Now()

	totals := study.Run(
		preparedWallets,
		households.Canada2023,
		*seed,
		*householdCount,
		*workers,
	)

	elapsed := time.Since(start)

	// Build results for all wallets.
	var results []walletResult

	for i, wallet := range wallets {
		var names []string

		for _, card := range wallet {
			names = append(names, card.Name)
		}

		averageCents := int64(math.Round(
			float64(totals[i]) / float64(*householdCount),
		))

		results = append(results, walletResult{
			Cards:                  names,
			AverageNetRewardsCents: averageCents,
			TotalNetRewardsCents:   totals[i],
		})
	}

	// Rank wallets using their exact total rewards.
	sort.Slice(results, func(i, j int) bool {
		if results[i].TotalNetRewardsCents != results[j].TotalNetRewardsCents {
			return results[i].TotalNetRewardsCents > results[j].TotalNetRewardsCents
		}

		return strings.Join(results[i].Cards, ",") <
			strings.Join(results[j].Cards, ",")
	})

	// Print summary.
	fmt.Printf("Cards: %d\n", len(cards))
	fmt.Printf("Wallet combinations: %d\n", len(wallets))
	fmt.Printf("Households: %d\n", *householdCount)
	fmt.Printf("Workers: %d\n", *workers)
	fmt.Printf("Elapsed time: %v\n\n", elapsed)

	fmt.Println("Top 10 wallets:")

	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("%d. $%.2f | %s\n",
			i+1,
			float64(results[i].AverageNetRewardsCents)/100,
			strings.Join(results[i].Cards, " + "),
		)
	}

	// Construct JSON output.
	report := studyResult{
		Households:     *householdCount,
		Workers:        *workers,
		Seed:           *seed,
		CardCount:      len(cards),
		WalletCount:    len(wallets),
		ElapsedSeconds: elapsed.Seconds(),
		Wallets:        results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Save JSON results to disk.
	if err := os.WriteFile(*output, append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResults saved to %s\n", *output)
}
