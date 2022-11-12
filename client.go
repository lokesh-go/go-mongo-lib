package mongodb

import "go.mongodb.org/mongo-driver/mongo"

// Client ...
type Client struct {
	database *mongo.Database
}
