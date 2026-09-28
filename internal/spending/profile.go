package spending

type Category int

const (
	Groceries Category = iota
	Dining
	Gas
	Transit
	Bills
	Retail
	Travel
	Other

	NumCategories
)

// Profile is one household's spending for a year: twelve months of
// per-category totals, in cents.
// like a 2D array
type Profile struct {
	Months [12][NumCategories]int64
}
