package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Insurance struct {
    ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    UserID            primitive.ObjectID `bson:"user_id" json:"user_id"`
    PolicyNumber      string             `bson:"policy_number" json:"policy_number"`
    Provider          string             `bson:"provider" json:"provider"`
    Type              string             `bson:"type" json:"type"`
    PremiumAmount     float64            `bson:"premium_amount" json:"premium_amount"`
    CoverageAmount    float64            `bson:"coverage_amount" json:"coverage_amount"`
    CoverageStartDate time.Time          `bson:"coverage_start_date" json:"coverage_start_date"`
    CoverageEndDate   time.Time          `bson:"coverage_end_date" json:"coverage_end_date"`
    TermsAndConditions string            `bson:"terms_and_conditions" json:"terms_and_conditions"`
    Status            string             `bson:"status" json:"status"`
    CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}
