package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client ...
type Client struct {
	database *mongo.Database
}

// CreateOne ...
func (c *Client) CreateOne(ctx context.Context, collection string, document interface{}) (err error) {
	// Hits DB
	_, err = c.database.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		return err
	}

	// Returns
	return nil
}

// ReadOne ...
func (c *Client) ReadOne(ctx context.Context, collection string, query interface{}) (res interface{}, err error) {
	// Hits DB
	err = c.database.Collection(collection).FindOne(ctx, query).Decode(&res)
	if err != nil {
		// Handles no document found
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	// Returns
	return res, nil
}

// UpdateOne ...
func (c *Client) UpdateOne(ctx context.Context, collection string, query interface{}, fields interface{}) (err error) {
	// Hits DB
	_, err = c.database.Collection(collection).UpdateOne(ctx, query, fields)
	if err != nil {
		return err
	}

	// Returns
	return nil
}

// DeleteOne ...
func (c *Client) DeleteOne(ctx context.Context, collection string, query interface{}) (count int64, err error) {
	// Hits DB
	result, err := c.database.Collection(collection).DeleteOne(ctx, query)
	if err != nil {
		return 0, err
	}
	count = result.DeletedCount

	// Returns
	return count, nil
}

// Read ...
func (c *Client) Read(ctx context.Context, collection string, query interface{}) (res []interface{}, err error) {
	// Hits DB
	cursor, err := c.database.Collection(collection).Find(ctx, query)
	if err != nil {
		return nil, err
	}

	// Close connection at the last
	defer cursor.Close(ctx)

	// Binds cursor response
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	// Returns
	return res, nil
}

// ReadWithProjection ...
func (c *Client) ReadWithProjection(ctx context.Context, collection string, query interface{}, projection interface{}) (res []interface{}, err error) {
	// Hits DB
	cursor, err := c.database.Collection(collection).Find(ctx, query, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}

	// Close connection at the last
	defer cursor.Close(ctx)

	// Binds cursor response
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	// Returns
	return res, nil
}
