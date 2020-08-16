package shared

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectId = primitive.ObjectID

func NewObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

func ObjectIdFromHex(hex string) (ObjectId, error) {
	return primitive.ObjectIDFromHex(hex)
}
