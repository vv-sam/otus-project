package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vv-sam/otus-project/server/internal/model/history"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NosqlRepository[T uniqueObject] struct {
	rc *redis.Client
	mc *mongo.Client

	redisKey string

	mongoDatabase   string
	mongoCollection string
}

type NosqlRepositoryOptions struct {
	RedisKey        string // Ключ для записи логов истории в Redis
	MongoDatabase   string // Название базы данных в MongoDB
	MongoCollection string // Название коллекции в MongoDB
}

func NewNosqlRepository[T uniqueObject](rc *redis.Client, mc *mongo.Client, opts NosqlRepositoryOptions) (*NosqlRepository[T], error) {
	if opts.MongoDatabase == "" {
		return nil, fmt.Errorf("mongo database is required")
	}

	if opts.MongoCollection == "" {
		return nil, fmt.Errorf("mongo collection is required")
	}

	if opts.RedisKey == "" {
		return nil, fmt.Errorf("redis key is required")
	}

	return &NosqlRepository[T]{
		rc:              rc,
		mc:              mc,
		redisKey:        opts.RedisKey,
		mongoDatabase:   opts.MongoDatabase,
		mongoCollection: opts.MongoCollection,
	}, nil
}

func (r *NosqlRepository[T]) Get(id uuid.UUID) (T, error) {
	db := r.mc.Database(r.mongoDatabase)
	collection := db.Collection(r.mongoCollection)

	res := collection.FindOne(context.Background(), bson.M{"id": id})
	if err := res.Err(); err != nil {
		var zero T
		if err == mongo.ErrNoDocuments {
			return zero, ErrNotFound
		}
		return zero, err
	}

	var item T
	if err := res.Decode(&item); err != nil {
		var zero T
		return zero, err
	}
	return item, nil
}

func (r *NosqlRepository[T]) GetAll() ([]T, error) {
	db := r.mc.Database(r.mongoDatabase)
	collection := db.Collection(r.mongoCollection)

	c, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var items []T
	if err := c.All(context.Background(), &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *NosqlRepository[T]) Add(item T) error {
	db := r.mc.Database(r.mongoDatabase)
	collection := db.Collection(r.mongoCollection)

	_, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		return err
	}

	logData, err := json.Marshal(history.Log[T]{
		Time:   time.Now(),
		Action: "add",
		Id:     item.GetId(),
		Data:   item,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	cmd := r.rc.RPush(context.Background(), r.redisKey, logData)
	if _, err := cmd.Result(); err != nil {
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	return nil
}

func (r *NosqlRepository[T]) Update(id uuid.UUID, item T) error {
	db := r.mc.Database(r.mongoDatabase)
	collection := db.Collection(r.mongoCollection)

	_, err := collection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": item})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}

	logData, err := json.Marshal(history.Log[T]{
		Time:   time.Now(),
		Action: "update",
		Id:     id,
		Data:   item,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	cmd := r.rc.RPush(context.Background(), r.redisKey, logData)
	if _, err := cmd.Result(); err != nil {
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	return nil
}

func (r *NosqlRepository[T]) Delete(id uuid.UUID) error {
	db := r.mc.Database(r.mongoDatabase)
	collection := db.Collection(r.mongoCollection)

	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}

	logData, err := json.Marshal(history.Log[T]{
		Time:   time.Now(),
		Action: "delete",
		Id:     id,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	cmd := r.rc.RPush(context.Background(), r.redisKey, logData)
	if _, err := cmd.Result(); err != nil {
		return fmt.Errorf("failed to push to redis: %w", err)
	}

	return nil
}

func (r *NosqlRepository[T]) GetHistory() ([]history.Log[T], error) {
	rr := r.rc.LRange(context.Background(), r.redisKey, 0, -1)
	res, err := rr.Result()
	if err != nil {
		return nil, err
	}

	logs := make([]history.Log[T], len(res))
	for _, rec := range res {
		var log history.Log[T]
		if err := json.Unmarshal([]byte(rec), &log); err != nil {
			return nil, fmt.Errorf("failed to unmarshal log: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}
