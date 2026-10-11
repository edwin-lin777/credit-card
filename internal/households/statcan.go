package households

import "github.com/edwin-lin777/credit-card/internal/spending"

// SpendingAnchors stores average annual Canadian household spending.
// All monetary values are in cents.
type SpendingAnchors struct {
	Year               int
	AnnualAverageCents [spending.NumCategories]int64
}

// Canada2023 contains baseline spending estimates for our simulation.
//
// Source: Statistics Canada, Survey of Household Spending, 2023.
// https://www150.statcan.gc.ca/n1/daily-quotidien/250521/dq250521a-eng.htm
//
// Note: Some categories are approximations using selected StatCan
// subcategories. They do not capture all credit-card-eligible spending.
var Canada2023 = SpendingAnchors{
	Year: 2023,

	AnnualAverageCents: [spending.NumCategories]int64{
		// Food purchased from stores: $8,659
		spending.Groceries: 865900,

		// Food purchased from restaurants: $3,351
		spending.Dining: 335100,

		// Gasoline and other vehicle fuels: $2,567
		spending.Gas: 256700,

		// Public transportation ($1,799) minus air travel ($1,158).
		// Proxy for non-air transportation: $641.
		spending.Transit: 64100,

		// Cell phone ($125.33/month) + internet ($87/month).
		// Partial bills estimate: $2,547.96 annually.
		spending.Bills: 254796,

		// Clothing ($2,739) + household furnishings
		// and equipment ($3,390): $6,129.
		// Partial retail estimate.
		spending.Retail: 612900,

		// Air travel ($1,158) + accommodation ($910).
		// Partial travel estimate: $2,068.
		spending.Travel: 206800,

		// Personal care ($1,860).
		// Partial proxy for other purchases.
		spending.Other: 186000,
	},
}
