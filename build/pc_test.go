package build

import (
	"math"
	"testing"
)

func TestRecommend_ValidBudgets(t *testing.T) {
	tests := []struct {
		budget   float64
		expected string
	}{
		{300, "Budget"},
		{400, "Budget"},
		{500, "Mid-Range"},
		{799, "Mid-Range"},
		{800, "High Performance"},
		{1200, "Enthusiast"},
		{2500, "Ultimate"},
		{999999, "God Tier"},
	}

	for _, tc := range tests {
		result := Recommend(tc.budget)
		if result.ErrorMessage != "" {
			t.Errorf("budget %v: unexpected error: %s", tc.budget, result.ErrorMessage)
		}
		if result.Category != tc.expected {
			t.Errorf("budget %v: expected category %q, got %q", tc.budget, tc.expected, result.Category)
		}
	}
}

func TestRecommend_LowBudget(t *testing.T) {
	result := Recommend(50)
	if result.ErrorMessage == "" {
		t.Error("expected error for low budget, got none")
	}
}

func TestRecommend_ZeroBudget(t *testing.T) {
	result := Recommend(0)
	if result.ErrorMessage == "" {
		t.Error("expected error for zero budget, got none")
	}
}

func TestRecommend_NegativeBudget(t *testing.T) {
	result := Recommend(-100)
	if result.ErrorMessage == "" {
		t.Error("expected error for negative budget, got none")
	}
}

func TestRecommend_ExactTierBoundary(t *testing.T) {
	for _, tier := range tiers {
		result := Recommend(tier.MinBudget)
		if result.ErrorMessage != "" {
			t.Errorf("tier %s at min budget %v: unexpected error: %s", tier.Category, tier.MinBudget, result.ErrorMessage)
		}
		if result.Category != tier.Category {
			t.Errorf("tier %s at budget %v: expected %s, got %s", tier.Category, tier.MinBudget, tier.Category, result.Category)
		}
	}
}

func TestRecommend_OverspendNotAllowed(t *testing.T) {
	halfway := (tiers[0].MinBudget + tiers[1].MinBudget) / 2
	result := Recommend(halfway)
	if result.Total != tiers[0].Total {
		t.Errorf("budget %v: expected total %v, got %v", halfway, tiers[0].Total, result.Total)
	}
}

func TestRecommend_FlooredBudget(t *testing.T) {
	result := Recommend(800.99)
	if result.Category != "High Performance" {
		t.Errorf("expected High Performance for 800.99 (floored to 800), got %s", result.Category)
	}
}

func TestRecommend_FieldsNotEmpty(t *testing.T) {
	for _, tier := range tiers {
		result := Recommend(tier.MinBudget)
		if result.CPU == "" || result.GPU == "" || result.RAM == "" || result.Storage == "" {
			t.Errorf("tier %s: empty component field", tier.Category)
		}
	}
}

func TestFindTier(t *testing.T) {
	_, ok := findTier(100)
	if ok {
		t.Error("expected false for budget below all tiers")
	}

	tier, ok := findTier(300)
	if !ok {
		t.Fatal("expected true for budget 300")
	}
	if tier.Category != "Budget" {
		t.Errorf("expected Budget, got %s", tier.Category)
	}
}

func TestMaxBudget(t *testing.T) {
	result := Recommend(maxBudget + 1)
	if result.ErrorMessage == "" {
		t.Error("expected error for exceeding max budget")
	}
}

func TestFloorIsApplied(t *testing.T) {
	result := Recommend(349.99)
	if math.Floor(349.99) != 349 {
		t.Fatal("floor should round down")
	}
	if result.Total != 350 {
		t.Errorf("expected total 350 for floored budget 349, got %v", result.Total)
	}
}
