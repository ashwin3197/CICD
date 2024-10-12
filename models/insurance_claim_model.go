package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsuranceClaim struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PolicyID    string             `json:"policy_id" bson:"policy_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	ClaimAmount float64            `json:"claim_amount" bson:"claim_amount"`
	ClaimDate   time.Time          `json:"claim_date" bson:"claim_date"`
	Status      string             `json:"status" bson:"status"`
}
