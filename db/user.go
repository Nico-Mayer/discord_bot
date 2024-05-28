package db

import (
	"database/sql"

	"github.com/disgoorg/snowflake/v2"
)

type DBUser struct {
	ID        snowflake.ID   `json:"id"`
	Name      string         `json:"name"`
	Exp       int            `json:"exp"`
	RiotPUUID sql.NullString `json:"riot_puuid"`
}

func InsertDBUser(discordUserID snowflake.ID, username string) error {
	query := "INSERT INTO users (id, name) VALUES ($1, $2)"
	_, err := DB.Exec(query, discordUserID, username)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(id snowflake.ID) (DBUser, error) {
	var user DBUser
	query := "SELECT * FROM users WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Exp, &user.RiotPUUID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UserInDatabase(id snowflake.ID) bool {
	query := "SELECT id FROM users WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&id)

	return err == nil
}

func (user *DBUser) SetRiotPUUID(puuid string) error {
	query := "UPDATE users SET riot_puuid = $1 WHERE id = $2"

	_, err := DB.Exec(query, puuid, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (user *DBUser) GrantExp(exp int) error {
	query := "UPDATE users SET exp = exp + $1 WHERE id = $2"

	_, err := DB.Exec(query, exp, user.ID)
	if err != nil {
		return err
	}

	return nil
}
