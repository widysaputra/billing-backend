package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/username/billing-backend/model"
)

// MessageResponseData adalah struktur untuk bagian 'data' dalam JSON
type MessageResponseData struct {
	Data  []model.MessageHistory `json:"data"`
	Total int                    `json:"total"`
}

// MessageResponse adalah struktur utama untuk respons JSON
type MessageResponse struct {
	Data MessageResponseData `json:"data"`
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		messages, total, err := model.GetAllMessages()
		if err != nil {
			log.Printf("Error getting message history from DB: %v", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		response := MessageResponse{
			Data: MessageResponseData{
				Data:  messages,
				Total: total,
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
