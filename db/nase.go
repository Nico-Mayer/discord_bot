package db

import (
	"time"

	"github.com/google/uuid"
)

type Nase struct {
	ID       uuid.UUID `json:"id"`
	UserID   string    `json:"userid"`
	AuthorID string    `json:"authorid"`
	Reason   string    `json:"reason"`
	Created  time.Time `json:"created"`
}

func GetNasenCount(userId string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM nasen WHERE userid = $1"
	err := DB.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GiveNase(users []User, authorId string, reason string) error {
	query := "INSERT INTO nasen (id, userid, authorid, reason, created) VALUES ($1, $2, $3, $4, $5)"
	for _, user := range users {

		nasenId := uuid.New()
		created := time.Now()

		if !user.InDatabase() {
			err := InsertUser(user.ID, user.Name)
			if err != nil {
				return err
			}
		}

		_, err := DB.Exec(query, nasenId, user.ID, authorId, reason, created)
		if err != nil {
			return err
		}
	}
	return nil
}
