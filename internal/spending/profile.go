package spending

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

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

func (c *Category) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "groceries":
		*c = Groceries
	case "dining":
		*c = Dining
	case "gas":
		*c = Gas
	case "transit":
		*c = Transit
	case "bills":
		*c = Bills
	case "retail":
		*c = Retail
	case "travel":
		*c = Travel
	case "other":
		*c = Other
	default:
		return fmt.Errorf("unknown spending category %q", value.Value)
	}

	return nil
}

// Profile is one household's spending for a year: twelve months of
// per-category totals, in cents.
type Profile struct {
	Months [12][NumCategories]int64
}
