package handler

import (
    "encoding/json"
    "net/http"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Contoh: username dan password hardcode (ganti sesuai kebutuhan)
    if req.Username == "admin" && req.Password == "admin123" {
        // Token bisa diganti dengan JWT atau random string
        resp := LoginResponse{Token: "simple-token-123"}
        json.NewEncoder(w).Encode(resp)
        return
    }

    http.Error(w, "Username atau password salah", http.StatusUnauthorized)
}