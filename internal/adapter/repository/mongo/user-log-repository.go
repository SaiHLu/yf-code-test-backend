package mongo

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	portrepository "codetest/internal/port/repository"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userLogRepository struct {
	DB *mongo.Client

	database   string
	collection string
}

func NewUserLogRepository(db *mongo.Client, database, collection string) portrepository.UserLogRepository {
	return &userLogRepository{
		DB:         db,
		database:   database,
		collection: collection,
	}
}

func (u *userLogRepository) Create(ctx context.Context, userLog *model.UserLogModel) error {
	coll := u.DB.Database(u.database).Collection(u.collection)
	_, err := coll.InsertOne(ctx, userLog)

	return err
}

func (u *userLogRepository) Find(ctx context.Context, request *dto.QueryUserLogRequest) ([]*model.UserLogModel, int64, error) {
	var userLogs []*model.UserLogModel

	coll := u.DB.Database(u.database).Collection(u.collection)

	total, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	pageSize := int64(request.PageSize)
	offset := pageSize * int64(request.Page-1)
	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetLimit(pageSize).SetSkip(offset).SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		log.Printf("Failed to find documents in collection %s.%s: %v", u.database, u.collection, err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &userLogs); err != nil {
		return nil, 0, err
	}

	return userLogs, total, nil
}
