package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Block struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Index        int                `json:"index,omitempty"`
	PreviousHash string             `json:"previousHash,omitempty"`
	Proof        int                `json:"proof,omitempty"`
	Timestamp    time.Time          `json:"timestamp,omitempty"`
	Miner        string             `json:"miner,omitempty"`
	Transaction  Transaction        `json:"transaction,omitempty"`
}
