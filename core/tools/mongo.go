package tools

import (
	"crypto/md5"
	"encoding/hex"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObjectID(s string) primitive.ObjectID {
	// First, try to parse if it's already a valid ObjectID hex string
	if len(s) == 24 {
		objectID, err := primitive.ObjectIDFromHex(s)
		if err == nil {
			return objectID
		}
	}

	// If not valid or different length, create hash from string
	hash := md5.Sum([]byte(s))
	// Take first 12 bytes from hash (ObjectID is 12 bytes)
	hexStr := hex.EncodeToString(hash[:])[:24]

	// Convert to ObjectID
	objectID, _ := primitive.ObjectIDFromHex(hexStr)
	return objectID
}
