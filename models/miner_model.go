package models

type Miner struct {
	Index       int     `json:"index,omitempty"`
	Name        string  `json:"name,omitempty"`
	BlocksMined int     `json:"blocksMined,omitempty"`
	TotalCoins  float64 `json:"totalCoins,omitempty"`
	Work        int     `json:"work,omitempty"`
	WorkPercent float64 `json:"workPercent,omitempty"`
}
