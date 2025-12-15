package repositories

import (
	"context"
	"time"

	"github.com/bomboskuy/UAS-Backend/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(achievement *models.Achievement) (string, error)
	FindByID(id string) (*models.Achievement, error)
	Update(id string, achievement *models.Achievement) error
	SoftDelete(id string) error
}

type achievementRepositoryMongo struct {
	collection *mongo.Collection
}

func NewAchievementRepositoryMongo(db *mongo.Database) AchievementRepository {
	return &achievementRepositoryMongo{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementRepositoryMongo) Create(a *models.Achievement) (string, error) {
	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), a)
	if err != nil {
		return "", err
	}

	return a.ID.Hex(), nil
}


func (r *achievementRepositoryMongo) FindByID(id string) (*models.Achievement, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement models.Achievement
	err = r.collection.FindOne(
		context.Background(),
		bson.M{
			"_id":       objID,
			"deletedAt": bson.M{"$exists": false},
		},
	).Decode(&achievement)

	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *achievementRepositoryMongo) Update(id string, achievement *models.Achievement) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	achievement.UpdatedAt = time.Now()

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": achievement},
	)

	return err
}

func (r *achievementRepositoryMongo) SoftDelete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"deletedAt": time.Now()}},
	)

	return err
}
