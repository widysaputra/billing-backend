package model

import "time"

// MessageHistory merepresentasikan satu baris data di tabel message_history
type MessageHistory struct {
	ID             int       `json:"id"`
	Target         string    `json:"target"`
	Message        string    `json:"message"`
	Status         string    `json:"status"`
	FonnteResponse string    `json:"fonnte_response"`
	SentAt         time.Time `json:"sent_at"`
}

// InsertMessageHistory menyimpan catatan pesan baru ke database
func InsertMessageHistory(m MessageHistory) error {
	_, err := DB.Exec("INSERT INTO message_history (target, message, status, fonnte_response) VALUES (?, ?, ?, ?)",
		m.Target, m.Message, m.Status, m.FonnteResponse)
	return err
}

// GetAllMessages mengambil semua riwayat pesan dari database
func GetAllMessages() ([]MessageHistory, int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM message_history").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := DB.Query("SELECT id, target, message, status, fonnte_response, sent_at FROM message_history ORDER BY sent_at DESC")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []MessageHistory
	for rows.Next() {
		var m MessageHistory
		if err := rows.Scan(&m.ID, &m.Target, &m.Message, &m.Status, &m.FonnteResponse, &m.SentAt); err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}
	return messages, total, rows.Err()
}