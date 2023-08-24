package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserBookmarks struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserId    string             `bson:"user_id" json:"userId" bson:"user_id"`
	ArticleId string             `bson:"article_id" json:"articleId" bson:"article_id"`
}