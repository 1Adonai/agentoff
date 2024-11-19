package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ContactForm struct {
	ID           int64
	Name         string
	ContactType  string
	ContactInfo  string
	SelectOption string
	Message      string
	IP           string
	CreatedAt    time.Time
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
        contact_type TEXT,
        contact_info TEXT,
        select_option TEXT,
        message TEXT,
        ip VARCHAR(45),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllContacts() ([]ContactForm, error) {
	query := `SELECT name, contact_type, contact_info, select_option, message, ip, created_at FROM contacts ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactForm
	for rows.Next() {
		var contact ContactForm
		err := rows.Scan(&contact.Name, &contact.ContactType, &contact.ContactInfo, &contact.SelectOption, &contact.Message, &contact.IP, &contact.CreatedAt)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}
func InsertContact(form ContactForm) error {
	query := `
        INSERT INTO contacts (
            name, contact_type, contact_info, select_option, message, ip, created_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := db.Exec(query,
        form.Name,
        form.ContactType,
        form.ContactInfo,
        form.SelectOption,
        form.Message,
        form.IP,
        time.Now(),
    )
    
    return err
}
