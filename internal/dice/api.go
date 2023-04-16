package dice

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"rolldice-go-api/pkg/log"
	"rolldice-go-api/pkg/mid"
	dicetransform "rolldice-go-api/pkg/model/transform"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

// RegisterHandlers registers handlers for specified path
func RegisterHandlersDice(r chi.Router, db *sqlx.DB, logger log.Logger, validate *validator.Validate) {
	r.Mount("/dice", RegisterHTTPHandlersDice(NewDiceHTTP(db, logger, validate)))
}

// RegisterHTTPHandlers registers http handlers for users endpoint
func RegisterHTTPHandlersDice(http HTTP) http.Handler {
	r := chi.NewRouter()
	r.With(mid.Paginate).Post("/", http.List)
	return r
}

func RollDice(Bet string, Dice int, BetPoint int) (dicetransform.ResultDiceCal, []byte) {
	var (
		status string
		is_win bool
		is_big bool
		point  int
	)
	sum := 0
	for i := 0; i < Dice; i++ {
		randomNumber := rand.Intn(6) + 1
		sum += randomNumber
	}

	if sum >= 6 {
		is_big = true
	} else {
		is_big = false
	}

	if Bet == "small" && is_big == false {
		is_win = true
	} else if Bet == "big" && is_big == true {
		is_win = true
	} else {
		is_win = false
	}

	if is_win == true {
		status = "win"
		point = BetPoint * 2
	} else {
		status = "lose"
		point = -BetPoint
	}

	resultDice := dicetransform.ResultDiceCal{
		DiceTotal:    sum,
		Status:       status,
		RecivedPoint: point,
	}
	smessageBytes, _ := json.Marshal(resultDice)
	return resultDice, smessageBytes
}
