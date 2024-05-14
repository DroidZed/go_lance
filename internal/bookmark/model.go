package bookmark

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserBookmarks struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserId    string             `bson:"user_id" json:"userId"`
	ArticleId string             `bson:"article_id" json:"articleId"`
}
