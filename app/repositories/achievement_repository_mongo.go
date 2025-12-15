package repositories

import (
	"context"

	"UAS-Backend/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type achievementRepositoryMongo struct {
	collection *mongo.Collection
}

func NewAchievementRepositoryMongo(db *mongo.Database) AchievementRepository {
	return &achievementRepositoryMongo{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementRepositoryMongo) Create(a *models.Achievement) error {
	_, err := r.collection.InsertOne(context.Background(), a)
	return err
}

func (r *achievementRepositoryMongo) FindByID(id string) (*models.Achievement, error) {
	var achievement models.Achievement
	err := r.collection.FindOne(
		context.Background(),
		bson.M{"_id": id},
	).Decode(&achievement)

	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *achievementRepositoryMongo) FindByStudentID(studentID string) ([]models.Achievement, error) {
	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"studentId": studentID},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var achievements []models.Achievement
	if err := cursor.All(context.Background(), &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *achievementRepositoryMongo) UpdateStatus(id string, status string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *achievementRepositoryMongo) Delete(id string) error {
	_, err := r.collection.DeleteOne(
		context.Background(),
		bson.M{"_id": id},
	)
	return err
}
