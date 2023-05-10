package transform

// Pagination to represents common parameters for endpoint that implement pagination
type ResultDice struct {
	Dice         int     `json:"dice"`
	UserId       string  `json:"userid"`
	EventName    string  `json:"event"`
	Status       string  `json:"status"`
	DiceTotal    []Dices `json:"dices"`
	RecivedPoint int     `json:"received_point"`
}

type PayloadPlay struct {
	Dice      int    `json:"dice"`
	UserId    string `json:"userid"`
	EventName string `json:"event"`
	Bet       string `json:"bet"`
	BetPoint  int    `json:"bet_point"`
	SessionID int    `json:"session_id"`
}

type Dices struct {
	DiceName  int `json:"dice_name"`
	DiceTotal int `json:"dice_total"`
}

type ResultDiceCal struct {
	DiceTotal    []Dices `json:"dices"`
	Status       string  `json:"status"`
	RecivedPoint int     `json:"received_point"`
}
