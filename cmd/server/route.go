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

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan Message)

type Message struct {
	Dice   string `json:"dice"`
	UserId string `json:"userid"`
}

// Routing setup api routing
func Routing(db *sqlx.DB, logger log.Logger) chi.Router {
	validate = validator.New()

	// setup server routing
	r := chi.NewRouter()

	// homepage welcome page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, "<html><head><title>Go Starter Kit</title></head><body>Welcome to Go Starter Kit</head></body></html>")
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
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				logging.Println("read:", err)
				break
			}

			var receivedMessage Message
			err = json.Unmarshal(message, &receivedMessage)

			smessage := Message{
				Dice:   receivedMessage.Dice,
				UserId: receivedMessage.UserId,
			}
			smessageBytes, err := json.Marshal(receivedMessage)

			logging.Printf("recv: %s", message)
			err = c.WriteMessage(mt, smessageBytes)
			if err != nil {
				logging.Println("write:", err)
				break
			}

			logging.Println(sessionID)

			broadcast <- smessage
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
			// Skip sending the message back to the sender
			// Convert the message to bytes
			messageBytes, err := json.Marshal(message)
			if err != nil {
				logging.Println("Failed to marshal message to JSON:", err)
				continue
			}

			if message.UserId == sessionID {
				logging.Println("KII", message)
				err = client.WriteMessage(websocket.TextMessage, messageBytes)
			}
			// Write the message to the client
			if err != nil {
				logging.Println("Failed to write message to WebSocket:", err)
				client.Close()
				delete(clients, sessionID)
			}
		}
	}
}
