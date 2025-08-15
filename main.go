package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/username/billing-backend/handler"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"github.com/username/billing-backend/model"
)

var db *sql.DB

// kirim pesan ke Fonnte
func sendFonnte(targets []string, message string, schedule string) error {
    token := os.Getenv("FONNTE_TOKEN")
    if token == "" {
        return fmt.Errorf("FONNTE_TOKEN env empty")
    }

    // build multipart/form-data body (sama seperti contoh PHP curl)
    var buf bytes.Buffer
    w := multipart.NewWriter(&buf)

    if err := w.WriteField("target", strings.Join(targets, ",")); err != nil {
        return err
    }
    if err := w.WriteField("message", message); err != nil {
        return err
    }
    if schedule != "" {
        if err := w.WriteField("schedule", schedule); err != nil {
            return err
        }
    }
    _ = w.Close()

    req, err := http.NewRequest("POST", "https://api.fonnte.com/send", &buf)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", token)
    req.Header.Set("Content-Type", w.FormDataContentType())

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading Fonnte response body: %v", err)
        // Lanjutkan proses meskipun body tidak bisa dibaca
    }
    fonnteRespString := string(bodyBytes)
    log.Println("Fonnte Response:", fonnteRespString)

    status := "success"
    if resp.StatusCode >= 300 {
        status = "failed"
    }

    // Simpan ke database
    history := model.MessageHistory{
        Target:         strings.Join(targets, ","),
        Message:        message,
        Status:         status,
        FonnteResponse: fonnteRespString,
    }
    if err := model.InsertMessageHistory(history); err != nil {
        log.Printf("Failed to save message history to DB: %v", err)
        // Jangan return error ke client, cukup log saja
    }

    if status == "failed" {
        return fmt.Errorf("Fonnte status %d: %s", resp.StatusCode, fonnteRespString)
    }
    return nil
}

func main() {
    // load .env jika ada (abaikan error kalau tidak ada)
    _ = godotenv.Load()

	// Koneksi ke database
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    var err error
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("DB connection error:", err)
    }
    if err = db.Ping(); err != nil {
        log.Fatal("DB ping error:", err)
    }
	// End koneksi ke database

	model.InitDB(db) // inisialisasi DB di model

    // debug: pastikan token terbaca
    log.Println("FONNTE_TOKEN (first 4 chars):", previewToken(os.Getenv("FONNTE_TOKEN")))

    http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        target := r.URL.Query().Get("target")
        msg := r.URL.Query().Get("message")
        schedule := r.URL.Query().Get("schedule") // ambil schedule dari query

        if target == "" || msg == "" {
            http.Error(w, "target & message required", http.StatusBadRequest)
            return
        }
        targets := strings.Split(target, ",")

        if err := sendFonnte(targets, msg, schedule); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write([]byte("Message sent!"))
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
	http.HandleFunc("/login", handler.LoginHandler)
	// Daftarkan handler untuk /contacts (tanpa slash) untuk GET all dan POST
	http.HandleFunc("/contacts", handler.ContactsHandler)
	// Daftarkan handler untuk /contacts/ (dengan slash) untuk GET by ID, PUT, dan DELETE
	http.HandleFunc("/contacts/", handler.ContactsHandler)
	http.HandleFunc("/messages", handler.MessagesHandler)
    log.Println("Server running on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func previewToken(t string) string {
	if len(t) <= 4 {
		return t
	}
	return t[:4] + "***"
}
