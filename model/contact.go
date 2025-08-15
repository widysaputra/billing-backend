package model

import (
    "database/sql"
)

type Contact struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Number string `json:"number"`
}

var DB *sql.DB

func InitDB(db *sql.DB) {
    DB = db
}

func GetAllContacts() ([]Contact, int, error) {
	// Pertama, hitung total kontak
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Kedua, ambil daftar kontak
    rows, err := DB.Query("SELECT id, name, number FROM contacts")
    if err != nil {
		return nil, 0, err
    }
    defer rows.Close()

    var contacts []Contact
    for rows.Next() {
        var c Contact
        if err := rows.Scan(&c.ID, &c.Name, &c.Number); err != nil {
			return nil, 0, err
        }
        contacts = append(contacts, c)
    }
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func InsertContact(c Contact) (int, error) {
    res, err := DB.Exec("INSERT INTO contacts (name, number) VALUES (?, ?)", c.Name, c.Number)
    if err != nil {
        return 0, err
    }
    id, err := res.LastInsertId()
    return int(id), err
}

func GetContactByID(id int) (Contact, error) {
	var c Contact
	row := DB.QueryRow("SELECT id, name, number FROM contacts WHERE id = ?", id)
	err := row.Scan(&c.ID, &c.Name, &c.Number)
	return c, err
}

func UpdateContact(c Contact) error {
    _, err := DB.Exec("UPDATE contacts SET name = ?, number = ? WHERE id = ?", c.Name, c.Number, c.ID)
    return err
}

func DeleteContact(id int) error {
    _, err := DB.Exec("DELETE FROM contacts WHERE id = ?", id)
    return err
}