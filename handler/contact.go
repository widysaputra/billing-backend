package handler

import (
    "encoding/json"
	"database/sql"
	"log"
    "net/http"
    "strconv"
    "strings"

    "github.com/username/billing-backend/model"
)

// ContactResponseData adalah struktur untuk bagian 'data' dalam JSON
type ContactResponseData struct {
	Data  []model.Contact `json:"data"`
	Total int             `json:"total"`
}

// ContactResponse adalah struktur utama untuk respons JSON
type ContactResponse struct {
	Data ContactResponseData `json:"data"`
}

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
     
    // Ambil ID dari path, jika ada. Contoh: /contacts/123 -> id = 123
    // Ini cara yang lebih robust daripada split manual.
    path := strings.TrimPrefix(r.URL.Path, "/contacts") // Hasilnya: "" atau "/123"
    path = strings.Trim(path, "/")                      // Hasilnya: "" atau "123"

    var id int
    if path != "" {
        var err error
        id, err = strconv.Atoi(path)
        if err != nil {
            http.Error(w, "Invalid contact ID in URL", http.StatusBadRequest)
            return
        }
    }

	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			// Kasus: GET /contacts/{id} -> Ambil satu kontak
			contact, err := model.GetContactByID(id)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Contact not found", http.StatusNotFound)
				} else {
					log.Printf("Error getting contact by ID from DB: %v", err)
					http.Error(w, "DB error", http.StatusInternalServerError)
				}
				return
			}
			json.NewEncoder(w).Encode(contact)
		} else {
			// Kasus: GET /contacts -> Ambil semua kontak
			contacts, total, err := model.GetAllContacts()
			if err != nil {
				log.Printf("Error getting contacts from DB: %v", err)
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
			response := ContactResponse{
				Data: ContactResponseData{
					Data:  contacts,
					Total: total,
				},
			}
			json.NewEncoder(w).Encode(response)
		}

	case http.MethodPost:
		// Kasus: POST /contacts -> Buat kontak baru
		if id != 0 {
			http.Error(w, "Cannot specify ID for a new contact", http.StatusBadRequest)
			return
		}
		var c model.Contact
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		newID, err := model.InsertContact(c)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		c.ID = newID
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)

	case http.MethodPut:
		// Kasus: PUT /contacts/{id} -> Update kontak
		if id == 0 {
			http.Error(w, "Contact ID is required for update", http.StatusBadRequest)
			return
		}
		var c model.Contact
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		c.ID = id // Pastikan ID dari URL yang digunakan
		if err := model.UpdateContact(c); err != nil {
			log.Printf("Error updating contact in DB: %v", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(c)

	case http.MethodDelete:
		// Kasus: DELETE /contacts/{id} -> Hapus kontak
		if id == 0 {
			http.Error(w, "Contact ID is required for delete", http.StatusBadRequest)
			return
		}
		if err := model.DeleteContact(id); err != nil {
			log.Printf("Error deleting contact from DB: %v", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Contact deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}