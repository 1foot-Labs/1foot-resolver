package models

type ActiveOrder struct {
	ID              string  `json:"_id"`
	Type            string  `json:"type"`
	MakerAddress    string  `json:"makerAddress"`
	PubKey		 	string  `json:"pubKey"`
	Sha256        	string  `json:"sha256"`
	Status          string  `json:"status"`
	AmountToReceive float64 `json:"amountToReceive"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}
