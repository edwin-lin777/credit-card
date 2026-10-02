package catalog

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/edwin-lin777/credit-card/internal/spending"
)

type Period int

const (
	Monthly Period = iota
	Annual
)

func (p *Period) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "monthly":
		*p = Monthly
	case "annual":
		*p = Annual
	default:
		return fmt.Errorf("unknown cap period %q", value.Value)
	}

	return nil
}

// Cap limits how much spending earns a group's bonus rate before it falls
// back to the card's base rate.
type Cap struct {
	AmountCents int64  `yaml:"amount_cents"`
	Period      Period `yaml:"period"`
}

type Group struct {
	Rates map[spending.Category]float64 `yaml:"rates"`
	Cap   *Cap                          `yaml:"cap"`
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
