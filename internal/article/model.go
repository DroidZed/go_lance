package article

import "go.mongodb.org/mongo-driver/bson/primitive"

type Article struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	Link     string             `json:"link" bson:"link"`
	Photo    string             `json:"photo" bson:"photo"`
	Tags     []string           `json:"tags" bson:"tags"`
	AuthorId string             `json:"authorId" bson:"author_id"`
}
