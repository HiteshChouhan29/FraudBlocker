package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"database/sql"

	"github.com/hiteshchouhan29/fraudblocker/db"

	"github.com/hiteshchouhan29/fraudblocker/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB, templates *template.Template, logger *log.Logger) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	blockedNumbers, err := db.GetAllBlockedNumbers(dbConn)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		logger.Printf("Error querying blocked numbers: %v", err)
		return
	}

	for i := range blockedNumbers {
		blockedNumbers[i].Number = utils.FormatPhoneNumber(blockedNumbers[i].Number)
	}

	err = templates.ExecuteTemplate(w, "index.html", blockedNumbers)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		logger.Printf("Template error: %v", err)
	}
}

func AddNumberHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB, logger *log.Logger) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	phoneNumber := strings.TrimSpace(r.FormValue("phone"))
	description := strings.TrimSpace(r.FormValue("description"))

	if !utils.IsValidPhoneNumber(phoneNumber) {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	normalizedPhone := utils.NormalizePhoneNumber(phoneNumber)

	id, err := db.AddBlockedNumber(dbConn, normalizedPhone, description)
	if err != nil {
		http.Error(w, "Error adding phone number", http.StatusInternalServerError)
		logger.Printf("Error adding phone number: %v", err)
		return
	}

	if id == 0 {
		http.Error(w, "Phone number already blocked", http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RemoveNumberHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB, logger *log.Logger) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.RemoveBlockedNumber(dbConn, id)
	if err != nil {
		if strings.Contains(err.Error(), "no phone number found") {
			http.Error(w, "Phone number not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error removing phone number", http.StatusInternalServerError)
		logger.Printf("Error removing phone number: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CheckNumberHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB, logger *log.Logger) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	phoneNumber := r.URL.Query().Get("phone")
	if phoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	normalizedPhone := utils.NormalizePhoneNumber(phoneNumber)

	isBlocked, err := db.IsNumberBlocked(dbConn, normalizedPhone)
	if err != nil {
		http.Error(w, "Error checking phone number", http.StatusInternalServerError)
		logger.Printf("Error checking phone number: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"blocked": isBlocked})
}

func RunMigrations(db *sql.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	return nil
}
