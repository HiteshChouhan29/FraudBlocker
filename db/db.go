package db

import (
	"database/sql"
	"fmt"

	"github.com/hiteshchouhan29/fraudblocker/models"
)

// AddBlockedNumber adds a blocked number using stored procedure
func AddBlockedNumber(db *sql.DB, number, description string) (int, error) {
	var id int
	err := db.QueryRow(
		"SELECT add_blocked_number($1, $2)",
		number,
		description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// RemoveBlockedNumber removes a blocked number by ID using stored procedure
func RemoveBlockedNumber(db *sql.DB, id int) error {
	var success bool
	err := db.QueryRow(
		"SELECT remove_blocked_number($1)",
		id,
	).Scan(&success)

	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("no phone number found with ID %d", id)
	}

	return nil
}

// IsNumberBlocked checks if a number is blocked using stored procedure
func IsNumberBlocked(db *sql.DB, number string) (bool, error) {
	var isBlocked bool
	err := db.QueryRow(
		"SELECT is_number_blocked($1)",
		number,
	).Scan(&isBlocked)

	if err != nil {
		return false, err
	}

	return isBlocked, nil
}

// GetAllBlockedNumbers retrieves all blocked numbers using stored procedure
func GetAllBlockedNumbers(db *sql.DB) ([]models.PhoneNumber, error) {
	rows, err := db.Query("SELECT * FROM get_all_blocked_numbers()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockedNumbers []models.PhoneNumber
	for rows.Next() {
		var bn models.PhoneNumber
		var createdAt sql.NullTime
		err := rows.Scan(&bn.ID, &bn.Number, &bn.Description, &createdAt)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			bn.CreatedAt = createdAt.Time.Format("Jan 02, 2006 15:04:05")
		}

		blockedNumbers = append(blockedNumbers, bn)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blockedNumbers, nil
}
