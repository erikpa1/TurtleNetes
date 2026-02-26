package dbclient

import (
	"context"
	"errors"
	"fmt"
	"time"
	"turtle/core/lgr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides generic CRUD operations for a specific entity type
type Repository[T any] struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// NewRepository creates a new repository for a specific entity type
func NewRepository[T any](client *Client, collectionName string) *Repository[T] {
	return &Repository[T]{
		collection: client.database.Collection(collectionName),
		timeout:    client.timeout,
	}
}

// InsertOne inserts a single document
func (r *Repository[T]) InsertOne(ctx context.Context, entity *T) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, nil
}

// InsertMany inserts multiple documents
func (r *Repository[T]) InsertMany(ctx context.Context, entities []T) ([]primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	docs := make([]interface{}, len(entities))
	for i := range entities {
		docs[i] = entities[i]
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	ids := make([]primitive.ObjectID, 0, len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			ids = append(ids, oid)
		}
	}

	return ids, nil
}

// FindOne finds a single document matching the filter
func (r *Repository[T]) FindOne(ctx context.Context, filter interface{}) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	return &entity, nil
}

// FindByID finds a document by its ObjectID
func (r *Repository[T]) FindByID(ctx context.Context, id primitive.ObjectID) (*T, error) {
	return r.FindOne(ctx, bson.M{"_id": id})
}

// FindMany finds all documents matching the filter
func (r *Repository[T]) FindMany(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var entities []T
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return entities, nil
}

// FindAll finds all documents in the collection
func (r *Repository[T]) FindAll(ctx context.Context, opts ...*options.FindOptions) ([]T, error) {
	return r.FindMany(ctx, bson.M{}, opts...)
}

// UpdateOne updates a single document matching the filter
func (r *Repository[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update document: %w", err)
	}

	return result.ModifiedCount, nil
}

// UpdateByID updates a document by its ObjectID
func (r *Repository[T]) UpdateByID(ctx context.Context, id primitive.ObjectID, update interface{}) (int64, error) {
	return r.UpdateOne(ctx, bson.M{"_id": id}, update)
}

// UpdateMany updates all documents matching the filter
func (r *Repository[T]) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update documents: %w", err)
	}

	return result.ModifiedCount, nil
}

// ReplaceOne replaces a single document matching the filter
func (r *Repository[T]) ReplaceOne(ctx context.Context, filter interface{}, replacement *T) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return 0, fmt.Errorf("failed to replace document: %w", err)
	}

	return result.ModifiedCount, nil
}

// DeleteOne deletes a single document matching the filter
func (r *Repository[T]) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete document: %w", err)
	}

	return result.DeletedCount, nil
}

// DeleteByID deletes a document by its ObjectID
func (r *Repository[T]) DeleteByID(ctx context.Context, id primitive.ObjectID) (int64, error) {
	return r.DeleteOne(ctx, bson.M{"_id": id})
}

// DeleteMany deletes all documents matching the filter
func (r *Repository[T]) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}

// Count counts documents matching the filter
func (r *Repository[T]) Count(ctx context.Context, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return count, nil
}

// Exists checks if any document matching the filter exists
func (r *Repository[T]) Exists(ctx context.Context, filter interface{}) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Aggregate performs an aggregation pipeline
func (r *Repository[T]) Aggregate(ctx context.Context, pipeline interface{}) ([]T, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return results, nil
}

// Distinct gets distinct values for a field
func (r *Repository[T]) Distinct(ctx context.Context, fieldName string, filter interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	values, err := r.collection.Distinct(ctx, fieldName, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct values: %w", err)
	}

	return values, nil
}

// Close closes the MongoDB client connection
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// GetCollection returns the underlying mongo.Collection for advanced operations
func (r *Repository[T]) GetCollection() *mongo.Collection {
	return r.collection
}

// QueryBuilder provides a fluent interface for building queries
type QueryBuilder[T any] struct {
	repo   *Repository[T]
	filter bson.M
	opts   *options.FindOptions
}

// NewQueryBuilder creates a new query builder
func (r *Repository[T]) NewQueryBuilder() *QueryBuilder[T] {
	return &QueryBuilder[T]{
		repo:   r,
		filter: bson.M{},
		opts:   options.Find(),
	}
}

// Where adds a filter condition
func (qb *QueryBuilder[T]) Where(key string, value interface{}) *QueryBuilder[T] {
	qb.filter[key] = value
	return qb
}

// WhereIn adds an $in filter condition
func (qb *QueryBuilder[T]) WhereIn(key string, values interface{}) *QueryBuilder[T] {
	qb.filter[key] = bson.M{"$in": values}
	return qb
}

// WhereGreaterThan adds a $gt filter condition
func (qb *QueryBuilder[T]) WhereGreaterThan(key string, value interface{}) *QueryBuilder[T] {
	qb.filter[key] = bson.M{"$gt": value}
	return qb
}

// WhereLessThan adds a $lt filter condition
func (qb *QueryBuilder[T]) WhereLessThan(key string, value interface{}) *QueryBuilder[T] {
	qb.filter[key] = bson.M{"$lt": value}
	return qb
}

// Limit sets the limit option
func (qb *QueryBuilder[T]) Limit(limit int64) *QueryBuilder[T] {
	qb.opts.SetLimit(limit)
	return qb
}

// Skip sets the skip option
func (qb *QueryBuilder[T]) Skip(skip int64) *QueryBuilder[T] {
	qb.opts.SetSkip(skip)
	return qb
}

// Sort sets the sort option
func (qb *QueryBuilder[T]) Sort(field string, order int) *QueryBuilder[T] {
	qb.opts.SetSort(bson.D{{Key: field, Value: order}})
	return qb
}

// Execute executes the query and returns results
func (qb *QueryBuilder[T]) Execute(ctx context.Context) ([]T, error) {
	return qb.repo.FindMany(ctx, qb.filter, qb.opts)
}

// First executes the query and returns the first result
func (qb *QueryBuilder[T]) First(ctx context.Context) (*T, error) {
	qb.opts.SetLimit(1)
	results, err := qb.repo.FindMany(ctx, qb.filter, qb.opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[0], nil
}

// Count counts documents matching the query
func (qb *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	return qb.repo.Count(ctx, qb.filter)
}

func Count(ctx context.Context, collection string, filter interface{}) int64 {
	tmp, err := MongoClient.database.Collection(collection).CountDocuments(ctx, filter)

	if err != nil {
		lgr.Error("failed to count documents: %w", err)
		return 0
	} else {
		return tmp
	}

}
