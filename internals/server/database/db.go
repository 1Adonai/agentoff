package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ContactForm struct {
	Name         string
	ContactType  string
	ContactInfo  string
	SelectOption string
	Message      string
}

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./internals/server/database/contacts.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS contacts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        contactType TEXT, -- Тип контакта 
        contactInfo TEXT, -- Сам контакт 
        selectOption TEXT,
        message TEXT
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllContacts() ([]ContactForm, error) {
	query := `SELECT name, contactType, contactInfo, selectOption, message FROM contacts ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactForm
	for rows.Next() {
		var contact ContactForm
		err := rows.Scan(&contact.Name, &contact.ContactType, &contact.ContactInfo, &contact.SelectOption, &contact.Message)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}
func InsertContact(form ContactForm) error {
	insertQuery := `INSERT INTO contacts (name, contactType, contactInfo, selectOption, message) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(insertQuery, form.Name, form.ContactType, form.ContactInfo, form.SelectOption, form.Message)
	return err
}
