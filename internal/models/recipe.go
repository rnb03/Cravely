package models

import "time"

// Recipe represents a cooking recipe
type Recipe struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PrepTime     int       `json:"prepTime"` // in minutes
	CookTime     int       `json:"cookTime"` // in minutes
	Servings     int       `json:"servings"`
	Cuisine      string    `json:"cuisine"`
	DietaryTags  []string  `json:"dietaryTags"` // e.g., "vegetarian", "gluten-free"
	ImageURL     string    `json:"imageUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// RecipeSearchParams defines parameters for recipe search
type RecipeSearchParams struct {
	Query       string   `json:"query"`
	Ingredients []string `json:"ingredients"`
	Cuisine     string   `json:"cuisine"`
	DietaryTags []string `json:"dietaryTags"`
	MaxPrepTime int      `json:"maxPrepTime"`
	MaxCookTime int      `json:"maxCookTime"`
	Page        int      `json:"page"`
	PageSize    int      `json:"pageSize"`
}

// RecipeResponse represents a paginated response of recipes
type RecipeResponse struct {
	Recipes    []Recipe `json:"recipes"`
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	TotalPages int      `json:"totalPages"`
}
