package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"voting_system/voting_system/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rdb         = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx         = context.Background()
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan models.Message)
	jwtSecret   = []byte("secret")
	voteSession = "vote_session"
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/create_vote", handleCreateVote)

	go handleMessages()

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		switch msg.Type {
		case "vote":
			claims, err := verifyToken(msg.Token)
			if err != nil {
				log.Println("Invalid token:", err)
				continue
			}

			rdb.HIncrBy(ctx, voteSession, msg.Vote, 1)

			result := rdb.HGetAll(ctx, voteSession).Val()
			result["user"] = claims.Username

			for client := range clients {
				err := client.WriteJSON(result)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func handleGetVote(w http.ResponseWriter, r *http.Request) {
	voteOptions := r.URL.Query().Get("options")
	if voteOptions == "" {
		http.Error(w, "Vote options are required", http.StatusBadRequest)
		return
	}

	options := []string{}
	err := json.Unmarshal([]byte(voteOptions), &options)
	if err != nil {
		http.Error(w, "Invalid options format", http.StatusBadRequest)
		return
	}

	sessionID := uuid.New().String()
	voteSession = sessionID

	for _, option := range options {
		rdb.HSet(ctx, sessionID, option, 0)
	}

	w.Write([]byte(fmt.Sprintf("Vote session %s created with options: %v", sessionID, options)))
}

type CreateVoteRequest struct {
	Options []string `json:"options"`
}

func handleCreateVote(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON body
	var reqBody CreateVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate options
	if len(reqBody.Options) == 0 {
		http.Error(w, "Vote options are required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received vote options: %v\n", reqBody.Options)

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Vote session created successfully"))
}

func verifyToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
