package server

import (
	"encoding/json"
	"fmt"
	logging "log"
	"net/http"
	"rolldice-go-api/pkg/mid"

	"rolldice-go-api/internal/dice"
	"rolldice-go-api/internal/healthcheck"
	"rolldice-go-api/internal/user"
	"rolldice-go-api/pkg/log"
	"rolldice-go-api/pkg/model/transform"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan transform.ResultDice)

// Routing setup api routing
func Routing(db *sqlx.DB, logger log.Logger) chi.Router {
	validate = validator.New()

	// setup server routing
	r := chi.NewRouter()

	// homepage welcome page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, "Running...")
	})

	//Web Socket Handler
	r.Get("/play", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logging.Print("upgrade:", err)
			return
		}
		defer c.Close()
		sessionID := fmt.Sprintf("%p", c)
		clients[sessionID] = c
		logging.Println(sessionID)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				logging.Println("read:", err)
				break
			}
			var receivedMessage transform.ResultDice
			err = json.Unmarshal(message, &receivedMessage)
			smessage := transform.ResultDice{
				Dice:      receivedMessage.Dice,
				UserId:    receivedMessage.UserId,
				EventName: receivedMessage.EventName,
			}
			// smessageBytes, err := json.Marshal(smessage)
			// if err != nil {
			// 	logging.Println("read:", err)
			// 	break
			// }
			// logging.Printf("recv: %s", message)
			// err = c.WriteMessage(mt, smessageBytes)
			if err != nil {
				logging.Println("write:", err)
				break
			}
			if smessage.EventName == "rolldice" {
				result_dice, data_result := dice.RollDice()
				MessageResult := transform.ResultDice{
					DiceTotal: int(result_dice.DiceTotal),
				}
				broadcast <- MessageResult
				err = c.WriteMessage(mt, data_result)
			} else {
				broadcast <- smessage
			}

		}
	})

	// register health check route
	healthcheck.RegisterHandlers(r)

	// register v1 api path group
	r.Route("/v1", func(r chi.Router) {
		r.Use(mid.APIVersionCtx("v1"))
		user.RegisterHandlers(r, db, logger, validate)
		dice.RegisterHandlersDice(r, db, logger, validate)
	})

	return r
}

func BroadcastMsg() {
	for {
		message := <-broadcast

		for sessionID, client := range clients {
			messageBytes, err := json.Marshal(message)
			if err != nil {
				logging.Println("Failed to marshal message to JSON:", err)
				continue
			}

			if message.UserId == sessionID {
				err = client.WriteMessage(websocket.TextMessage, messageBytes)
			}
			if err != nil {
				logging.Println("Failed to write message to WebSocket:", err)
				client.Close()
				delete(clients, sessionID)
			}
		}
	}
}
