package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Block struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Index        int                `json:"index,omitempty"`
	PreviousHash string             `json:"previousHash,omitempty"`
	Proof        int                `json:"proof,omitempty"`
	Timestamp    time.Time          `json:"timestamp,omitempty"`
	Miner        string             `json:"miner,omitempty"`
	Signature    string             `json:"signature,omitempty"`
	Transaction  Transaction        `json:"transaction,omitempty"`
}
