package rewards

import (
	"sort"

	"github.com/edwin-lin777/credit-card/internal/catalog"
	"github.com/edwin-lin777/credit-card/internal/spending"
)

type walletOffer struct {
	cardIndex  int
	groupIndex int
	category   spending.Category
	rate       float64
}

// Stores the wallet's cards and sorted reward offers.
type PreparedWallet struct {
	cards  []catalog.Card
	offers []walletOffer
}

// Build and sort reward offers ONCE.
func PrepareWallet(cards []catalog.Card) PreparedWallet {
	w := PreparedWallet{cards: cards}

	for i, card := range cards {
		for cat := spending.Category(0); cat < spending.NumCategories; cat++ {
			w.offers = append(w.offers, walletOffer{
				cardIndex:  i,
				groupIndex: -1,
				category:   cat,
				rate:       card.BaseRate,
			})
		}

		for j, group := range card.Groups {
			for cat, rate := range group.Rates {
				if rate <= card.BaseRate {
					continue
				}

				w.offers = append(w.offers, walletOffer{
					cardIndex:  i,
					groupIndex: j,
					category:   cat,
					rate:       rate,
				})
			}
		}
	}

	sort.Slice(w.offers, func(i, j int) bool {
		a, b := w.offers[i], w.offers[j]

		if a.rate != b.rate {
			return a.rate > b.rate
		}
		if a.cardIndex != b.cardIndex {
			return a.cardIndex < b.cardIndex
		}
		if a.groupIndex != b.groupIndex {
			return a.groupIndex < b.groupIndex
		}
		return a.category < b.category
	})

	return w
}

// Evaluate households using the already-sorted offers.
func EvaluatePreparedWallet(w PreparedWallet, p spending.Profile) int64 {
	if len(w.cards) == 0 {
		return 0
	}

	assigned := make([]spending.Profile, len(w.cards))
	annualUsed := make([][]int64, len(w.cards))

	for i, card := range w.cards {
		annualUsed[i] = make([]int64, len(card.Groups))
	}

	for month := 0; month < 12; month++ {
		remaining := p.Months[month]

		monthlyUsed := make([][]int64, len(w.cards))

		for i, card := range w.cards {
			monthlyUsed[i] = make([]int64, len(card.Groups))
		}

		for _, offer := range w.offers {
			amount := remaining[offer.category]

			if amount <= 0 {
				continue
			}

			if offer.groupIndex >= 0 {
				group := w.cards[offer.cardIndex].Groups[offer.groupIndex]

				if group.Cap != nil {
					used := &annualUsed[offer.cardIndex][offer.groupIndex]

					if group.Cap.Period == catalog.Monthly {
						used = &monthlyUsed[offer.cardIndex][offer.groupIndex]
					}

					available := group.Cap.AmountCents - *used

					if available <= 0 {
						continue
					}

					if amount > available {
						amount = available
					}

					*used += amount
				}
			}

			assigned[offer.cardIndex].Months[month][offer.category] += amount
			remaining[offer.category] -= amount
		}
	}

	var total int64

	for i, card := range w.cards {
		total += Evaluate(card, assigned[i])
	}

	return total
}

// Keep the original function for compatibility.
func EvaluateWallet(cards []catalog.Card, p spending.Profile) int64 {
	w := PrepareWallet(cards)
	return EvaluatePreparedWallet(w, p)
}
