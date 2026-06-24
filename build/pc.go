package build

import "math"

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
	{300, "Budget", "AMD Ryzen 3 4100", "GTX 1650", "8GB DDR4 3200MHz", "256GB NVMe SSD", 350},
	{500, "Mid-Range", "Intel i5-12400F", "RTX 3060 12GB", "16GB DDR4 3200MHz", "512GB NVMe SSD", 650},
	{800, "High Performance", "AMD Ryzen 5 7600", "RTX 4060 Ti 8GB", "16GB DDR5 5600MHz", "1TB NVMe SSD", 950},
	{1200, "Enthusiast", "Intel i7-13700KF", "RTX 4070 12GB", "32GB DDR5 6000MHz", "2TB NVMe SSD", 1450},
	{1800, "Extreme", "AMD Ryzen 9 7900X", "RTX 4080 Super 16GB", "32GB DDR5 6000MHz", "2TB NVMe SSD", 1950},
	{2500, "Ultimate", "Intel i9-14900K", "RTX 4090 24GB", "64GB DDR5 6400MHz", "4TB NVMe SSD Gen5", 3200},
	{4000, "God Tier", "AMD Threadripper 7960X", "RTX 5090 32GB", "128GB DDR5 6400MHz", "8TB NVMe SSD Gen5", 5500},
}

func Recommend(budget float64) Build {
	if budget <= 0 {
		return Build{ErrorMessage: "Please enter a valid budget amount."}
	}

	budget = math.Floor(budget)

	var best Tier
	for _, t := range tiers {
		if budget >= t.MinBudget {
			best = t
		}
	}

	if best.Total == 0 {
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
