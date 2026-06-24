package build

import "math"

var maxBudget = 1_000_000.0

type Build struct {
	Budget       float64
	Total        float64
	CPU          string
	GPU          string
	RAM          string
	Storage      string
	Category     string
	ErrorMessage string
}

type Tier struct {
	MinBudget float64
	Category  string
	CPU       string
	GPU       string
	RAM       string
	Storage   string
	Total     float64
}

var tiers = []Tier{
	{MinBudget: 300, Category: "Budget", CPU: "AMD Ryzen 3 4100", GPU: "GTX 1650", RAM: "8GB DDR4 3200MHz", Storage: "256GB NVMe SSD", Total: 350},
	{MinBudget: 500, Category: "Mid-Range", CPU: "Intel i5-12400F", GPU: "RTX 3060 12GB", RAM: "16GB DDR4 3200MHz", Storage: "512GB NVMe SSD", Total: 650},
	{MinBudget: 800, Category: "High Performance", CPU: "AMD Ryzen 5 7600", GPU: "RTX 4060 Ti 8GB", RAM: "16GB DDR5 5600MHz", Storage: "1TB NVMe SSD", Total: 950},
	{MinBudget: 1200, Category: "Enthusiast", CPU: "Intel i7-13700KF", GPU: "RTX 4070 12GB", RAM: "32GB DDR5 6000MHz", Storage: "2TB NVMe SSD", Total: 1450},
	{MinBudget: 1800, Category: "Extreme", CPU: "AMD Ryzen 9 7900X", GPU: "RTX 4080 Super 16GB", RAM: "32GB DDR5 6000MHz", Storage: "2TB NVMe SSD", Total: 1950},
	{MinBudget: 2500, Category: "Ultimate", CPU: "Intel i9-14900K", GPU: "RTX 4090 24GB", RAM: "64GB DDR5 6400MHz", Storage: "4TB NVMe SSD Gen5", Total: 3200},
	{MinBudget: 4000, Category: "God Tier", CPU: "AMD Threadripper 7960X", GPU: "RTX 5090 32GB", RAM: "128GB DDR5 6400MHz", Storage: "8TB NVMe SSD Gen5", Total: 5500},
}

func Recommend(budget float64) Build {
	budget = math.Floor(budget)

	if budget <= 0 {
		return Build{ErrorMessage: "Please enter a valid budget amount greater than $0."}
	}
	if budget > maxBudget {
		return Build{ErrorMessage: "Budget is too large. Maximum supported is $1,000,000."}
	}

	best, ok := findTier(budget)
	if !ok {
		return Build{
			Budget:       budget,
			ErrorMessage: "Budget too low. Minimum ~$300 recommended for a basic PC.",
		}
	}

	return Build{
		Budget:   budget,
		Total:    best.Total,
		CPU:      best.CPU,
		GPU:      best.GPU,
		RAM:      best.RAM,
		Storage:  best.Storage,
		Category: best.Category,
	}
}

func findTier(budget float64) (Tier, bool) {
	for i := len(tiers) - 1; i >= 0; i-- {
		if budget >= tiers[i].MinBudget {
			return tiers[i], true
		}
	}
	return Tier{}, false
}
