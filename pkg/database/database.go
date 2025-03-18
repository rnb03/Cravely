package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client represents the database client
type Client struct {
	*mongo.Client
	RecipeCollection *mongo.Collection
}

// NewClient creates a new database client
func NewClient(uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to MongoDB")

	db := client.Database(dbName)
	recipeCollection := db.Collection("recipes")

	return &Client{
		Client:           client,
		RecipeCollection: recipeCollection,
	}, nil
}

// Close disconnects from the database
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from database: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}
