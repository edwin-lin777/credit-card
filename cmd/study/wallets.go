package main

import "github.com/edwin-lin777/credit-card/internal/catalog"

// buildWallets generates every unique 2-card and 3-card combination.
func buildWallets(cards []catalog.Card) [][]catalog.Card {
	var wallets [][]catalog.Card

	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {

			// Two-card wallet.
			wallets = append(wallets, []catalog.Card{
				cards[i], cards[j],
			})

			// Three-card wallets.
			for k := j + 1; k < len(cards); k++ {
				wallets = append(wallets, []catalog.Card{
					cards[i], cards[j], cards[k],
				})
			}
		}
	}

	return wallets
}
