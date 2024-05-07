package db

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
)

type Nase struct {
	ID       uuid.UUID `json:"id"`
	UserID   string    `json:"userid"`
	AuthorID string    `json:"authorid"`
	Reason   string    `json:"reason"`
	Created  time.Time `json:"created"`
}

func GetNasenForUser(dbUserID snowflake.ID) ([]Nase, error) {
	var nasen []Nase

	query := "SELECT * FROM nasen WHERE userid = $1 ORDER BY created DESC"

	rows, err := DB.Query(query, dbUserID)
	if err != nil {
		return []Nase{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var nase Nase

		err := rows.Scan(&nase.ID, &nase.UserID, &nase.AuthorID, &nase.Reason, &nase.Created)
		if err != nil {
			return []Nase{}, err
		}

		nasen = append(nasen, nase)
	}

	return nasen, nil
}

func GetNasenCountForUser(dbUserID snowflake.ID) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM nasen WHERE userid = $1"
	err := DB.QueryRow(query, dbUserID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
