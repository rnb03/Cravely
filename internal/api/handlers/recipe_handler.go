package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"creavely/internal/models"
	"creavely/internal/services"
)

// RecipeHandler handles HTTP requests for recipes
type RecipeHandler struct {
	recipeService *services.RecipeService
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipeService: recipeService}
}

// GetRecipe handles GET requests for a specific recipe
func (h *RecipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	recipe, err := h.recipeService.GetRecipeByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// SearchRecipes handles GET requests for searching recipes
func (h *RecipeHandler) SearchRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse search parameters
	params := models.RecipeSearchParams{
		Query:   query.Get("query"),
		Cuisine: query.Get("cuisine"),
	}

	// Parse ingredients
	if ingredients := query.Get("ingredients"); ingredients != "" {
		params.Ingredients = strings.Split(ingredients, ",")
	}

	// Parse dietary tags
	if dietaryTags := query.Get("dietaryTags"); dietaryTags != "" {
		params.DietaryTags = strings.Split(dietaryTags, ",")
	}

	// Parse max prep time
	if maxPrepTime := query.Get("maxPrepTime"); maxPrepTime != "" {
		if val, err := strconv.Atoi(maxPrepTime); err == nil {
			params.MaxPrepTime = val
		}
	}

	// Parse max cook time
	if maxCookTime := query.Get("maxCookTime"); maxCookTime != "" {
		if val, err := strconv.Atoi(maxCookTime); err == nil {
			params.MaxCookTime = val
		}
	}

	// Parse pagination parameters
	if page := query.Get("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil && val > 0 {
			params.Page = val
		}
	}

	if pageSize := query.Get("pageSize"); pageSize != "" {
		if val, err := strconv.Atoi(pageSize); err == nil && val > 0 {
			params.PageSize = val
		}
	}

	// Search recipes
	response, err := h.recipeService.SearchRecipes(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateRecipe handles POST requests for creating a new recipe
func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdRecipe, err := h.recipeService.CreateRecipe(r.Context(), recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRecipe)
}

// UpdateRecipe handles PUT requests for updating an existing recipe
func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedRecipe, err := h.recipeService.UpdateRecipe(r.Context(), id, recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRecipe)
}

// DeleteRecipe handles DELETE requests for deleting a recipe
func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.recipeService.DeleteRecipe(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
