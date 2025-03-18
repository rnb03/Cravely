package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"creavely/internal/api/handlers"
	"creavely/internal/api/middleware"
	"creavely/internal/services"
	"creavely/pkg/database"
)

// Router sets up the HTTP routes
func NewRouter(db *database.Client) *mux.Router {
	router := mux.NewRouter()

	// Create services
	recipeService := services.NewRecipeService(db)

	// Create handlers
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	// Apply middleware
	router.Use(middleware.Logger)
	router.Use(middleware.CORS)
	router.Use(middleware.RecoverPanic)

	// Set up routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Recipe routes
	apiRouter.HandleFunc("/recipes", recipeHandler.SearchRecipes).Methods("GET")
	apiRouter.HandleFunc("/recipes", recipeHandler.CreateRecipe).Methods("POST")
	apiRouter.HandleFunc("/recipes/{id}", recipeHandler.GetRecipe).Methods("GET")
	apiRouter.HandleFunc("/recipes/{id}", recipeHandler.UpdateRecipe).Methods("PUT")
	apiRouter.HandleFunc("/recipes/{id}", recipeHandler.DeleteRecipe).Methods("DELETE")

	// Health check
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
