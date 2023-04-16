package dice

import (
	"net/http"

	"rolldice-go-api/pkg/log"
	"rolldice-go-api/pkg/mid"

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
