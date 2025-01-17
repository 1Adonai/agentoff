package handlers

import (
	"agentoff/internals/server/database"
	"agentoff/internals/server/keys"
	"agentoff/internals/server/ratelimit"
	"agentoff/internals/server/telegram"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	_ "agentoff/internals/server/telegram"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("key"))
var telegramBot *telegram.TelegramBot

func init() {
	var err error
	chat_id, err := strconv.ParseInt(keys.GetEnv("telegram_chat_id"), 10, 64)
	telegramBot, err = telegram.NewTelegramBot(keys.GetEnv("telegram_bot_token"), chat_id)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	tmpl, err := template.ParseFiles(path.Join("templates", templateName), path.Join("templates", "ContactForm.html"))
	if err != nil {
		http.Error(w, "Unable to parse template", http.StatusInternalServerError)
	}
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	data := struct {
		Message string
	}{
		Message: message,
	}
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}
	tmpl.Execute(w, data)

}

func ParseContactForm(r *http.Request) (database.ContactForm, error) {
	err := r.ParseForm()
	if err != nil {
		return database.ContactForm{}, err
	}

	// Get IP address
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	return database.ContactForm{
		Name:         r.FormValue("name"),
		ContactType:  r.FormValue("contactType"),
		ContactInfo:  r.FormValue("contactInfo"),
		SelectOption: r.FormValue("selectOption"),
		Message:      r.FormValue("message"),
		IP:           ip,
	}, nil
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Received a POST request to /contact")

		// Get client's IP address
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		// Check if the request is allowed
		if !ratelimit.IsAllowed(ip) {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			session, _ := store.Get(r, "session-name")
			session.Values["message"] = "Слишком много запросов."
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Парсинг формы
		contactForm, err := ParseContactForm(r)
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
			return
		}
		// Send contact info to Telegram
		err = telegramBot.SendContactInfo(contactForm.Name, contactForm.ContactType, contactForm.ContactInfo, contactForm.Message, contactForm.SelectOption, contactForm.IP)
		if err != nil {
			log.Printf("Failed to send contact info to Telegram: %v", err)
			// Handle error appropriately
		}

		log.Printf("Parsed form: %+v", contactForm)

		if err := database.InsertContact(contactForm); err != nil {
			log.Printf("Error inserting data into database: %v", err)
			http.Error(w, "Ошибка вставки данных в базу", http.StatusInternalServerError)
			return
		}

		// Установка сообщения в сессию
		session, _ := store.Get(r, "session-name")
		session.Values["message"] = "Сообщение успешно отправлено! Мы скоро свяжемся с вами."
		err = session.Save(r, w) // Добавлено сохранение сессии
		if err != nil {
			log.Printf("Error saving session: %v", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Printf("Unsupported request method: %s", r.Method)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index.html")
}
func OsagoHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "OSAGO.html")
}

func KaskoHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "KASKO.html")
}

func HouseHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "HOUSE.html")
}

func DomHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "DOM.html")
}
