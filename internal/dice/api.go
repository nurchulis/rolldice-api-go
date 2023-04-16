package dice

import (
	"encoding/json"
	logging "log"
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

func RollDice() (dicetransform.ResultDiceCal, []byte) {
	randomNumber := rand.Intn(6) + 1
	logging.Println("IKII", randomNumber)
	resultDice := dicetransform.ResultDiceCal{
		DiceTotal: int(randomNumber),
	}
	smessageBytess, _ := json.Marshal(resultDice)
	return resultDice, smessageBytess
}
