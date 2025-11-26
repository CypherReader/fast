package domain

type DietType string

const (
	DietVegan      DietType = "vegan"
	DietVegetarian DietType = "vegetarian"
	DietNormal     DietType = "normal"
)

type Recipe struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Diet         DietType `json:"diet"`
	IsSimple     bool     `json:"is_simple"` // Highlight simple recipes
	Calories     int      `json:"calories"`
	Carbs        int      `json:"carbs"` // Net carbs
	Image        string   `json:"image"` // URL or placeholder
}
