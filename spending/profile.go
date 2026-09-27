package rewards

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

type Profile struct {
	Months [12][NumCategories]int64
}