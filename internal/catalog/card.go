package catalog

import "github.com/edwin-lin777/credit-card/internal/spending"

type Period int

const (
	Monthly Period = iota
	Annual
)

// Cap limits how much spending earns a group's bonus rate before it falls
// back to the card's base rate.
type Cap struct {
	AmountCents int64
	Period      Period
}

// Group is a set of categories that share one cap. A nil Cap means uncapped.
type Group struct {
	Rates map[spending.Category]float64
	Cap   *Cap
}

// Card is one credit card's earn rules.
type Card struct {
	ID             string  `yaml:"id"`
	Name           string  `yaml:"name"`
	AnnualFeeCents int64   `yaml:"annual_fee_cents"`
	BaseRate       float64 `yaml:"base_rate"`
	Groups         []Group `yaml:"groups"`

	SourceURL  string `yaml:"source_url"`
	VerifiedOn string `yaml:"verified_on"`
}
