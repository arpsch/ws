package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Metadata struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Name        string `json:"name,omitempty"`
	Size        int64  `json:"size,omitempty"`
	ContentType string `json:"content_type,omitempty" bson:"content_type,omitempty"`
}
