package transform

// Pagination to represents common parameters for endpoint that implement pagination
type ResultDice struct {
	Dice         int    `json:"dice"`
	UserId       string `json:"userid"`
	EventName    string `json:"event"`
	Status       string `json:"status"`
	DiceTotal    int    `json:"dice_total"`
	RecivedPoint int    `json:"received_point"`
}

type PayloadPlay struct {
	Dice      int    `json:"dice"`
	UserId    string `json:"userid"`
	EventName string `json:"event"`
	Bet       string `json:"bet"`
	BetPoint  int    `json:"bet_point"`
	RoomId    int    `json:"received_point"`
}

type ResultDiceCal struct {
	DiceTotal    int    `json:"dice_total"`
	Status       string `json:"status"`
	RecivedPoint int    `json:"received_point"`
}
