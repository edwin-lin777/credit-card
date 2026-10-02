package catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/edwin-lin777/credit-card/internal/spending"
)

func TestLoadCardWithGroups(t *testing.T) {
	data := `
id: test-bonus
name: Test Bonus Card
annual_fee_cents: 0
base_rate: 0.01

groups:
  - rates:
      groceries: 0.04
      bills: 0.04
    cap:
      amount_cents: 2500000
      period: annual
`

	dir := t.TempDir()
	path := filepath.Join(dir, "card.yaml")

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	card, err := LoadCard(path)
	if err != nil {
		t.Fatal(err)
	}

	if card.BaseRate != 0.01 {
		t.Errorf("BaseRate = %f, want 0.01", card.BaseRate)
	}

	if len(card.Groups) != 1 {
		t.Fatalf("len(Groups) = %d, want 1", len(card.Groups))
	}

	group := card.Groups[0]

	if group.Rates[spending.Groceries] != 0.04 {
		t.Errorf("groceries rate = %f, want 0.04", group.Rates[spending.Groceries])
	}

	if group.Rates[spending.Bills] != 0.04 {
		t.Errorf("bills rate = %f, want 0.04", group.Rates[spending.Bills])
	}

	if group.Cap == nil {
		t.Fatal("Cap is nil")
	}

	if group.Cap.AmountCents != 2500000 {
		t.Errorf("AmountCents = %d, want 2500000", group.Cap.AmountCents)
	}

	if group.Cap.Period != Annual {
		t.Errorf("Period = %d, want Annual", group.Cap.Period)
	}
}
