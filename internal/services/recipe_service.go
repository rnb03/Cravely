package services

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"creavely/internal/models"
	"creavely/pkg/database"
)

// RecipeService handles business logic for recipes
type RecipeService struct {
	db *database.Client
}

// NewRecipeService creates a new recipe service
func NewRecipeService(db *database.Client) *RecipeService {
	return &RecipeService{db: db}
}

// GetRecipeByID retrieves a recipe by its ID
func (s *RecipeService) GetRecipeByID(ctx context.Context, id string) (*models.Recipe, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid recipe ID")
	}

	var recipe models.Recipe
	err = s.db.RecipeCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&recipe)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("recipe not found")
		}
		return nil, err
	}

	return &recipe, nil
}

// SearchRecipes searches for recipes based on the provided parameters
func (s *RecipeService) SearchRecipes(ctx context.Context, params models.RecipeSearchParams) (*models.RecipeResponse, error) {
	filter := bson.M{}

	// Add query filter if provided
	if params.Query != "" {
		filter["$or"] = bson.A{
			bson.M{"title": bson.M{"$regex": params.Query, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": params.Query, "$options": "i"}},
		}
	}

	// Add ingredients filter if provided
	if len(params.Ingredients) > 0 {
		filter["ingredients"] = bson.M{"$all": params.Ingredients}
	}

	// Add cuisine filter if provided
	if params.Cuisine != "" {
		filter["cuisine"] = params.Cuisine
	}

	// Add dietary tags filter if provided
	if len(params.DietaryTags) > 0 {
		filter["dietaryTags"] = bson.M{"$all": params.DietaryTags}
	}

	// Add prep time filter if provided
	if params.MaxPrepTime > 0 {
		filter["prepTime"] = bson.M{"$lte": params.MaxPrepTime}
	}

	// Add cook time filter if provided
	if params.MaxCookTime > 0 {
		filter["cookTime"] = bson.M{"$lte": params.MaxCookTime}
	}

	// Set default pagination values if not provided
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Calculate skip value for pagination
	skip := (params.Page - 1) * params.PageSize

	// Count total matching documents
	totalCount, err := s.db.RecipeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	// Find matching recipes with pagination
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(params.PageSize)).
		SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	cursor, err := s.db.RecipeCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode recipes
	var recipes []models.Recipe
	if err := cursor.All(ctx, &recipes); err != nil {
		return nil, err
	}

	return &models.RecipeResponse{
		Recipes:    recipes,
		TotalCount: int(totalCount),
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// CreateRecipe creates a new recipe
func (s *RecipeService) CreateRecipe(ctx context.Context, recipe models.Recipe) (*models.Recipe, error) {
	// Set creation and update times
	now := time.Now()
	recipe.CreatedAt = now
	recipe.UpdatedAt = now

	// Insert recipe into database
	result, err := s.db.RecipeCollection.InsertOne(ctx, recipe)
	if err != nil {
		return nil, err
	}

	// Set the ID of the recipe
	recipe.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return &recipe, nil
}

// UpdateRecipe updates an existing recipe
func (s *RecipeService) UpdateRecipe(ctx context.Context, id string, recipe models.Recipe) (*models.Recipe, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid recipe ID")
	}

	// Set update time
	recipe.UpdatedAt = time.Now()

	// Update recipe in database
	update := bson.M{
		"$set": recipe,
	}
	result, err := s.db.RecipeCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("recipe not found")
	}

	// Set the ID of the recipe
	recipe.ID = id

	return &recipe, nil
}

// DeleteRecipe deletes a recipe by its ID
func (s *RecipeService) DeleteRecipe(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid recipe ID")
	}

	result, err := s.db.RecipeCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("recipe not found")
	}

	return nil
}
