package db

import (
	"sort"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	NasenCount int    `json:"nasen_count"`
	Exp        int    `json:"exp"`
	Level      int    `json:"level"`
}

func InsertUser(id string, username string) error {
	query := "INSERT INTO users (id, name) VALUES ($1, $2)"
	_, err := DB.Exec(query, id, username)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) InDatabase() bool {
	var id string
	query := "SELECT id FROM users WHERE id = $1"
	err := DB.QueryRow(query, user.ID).Scan(&id)

	return err == nil
}

func (user *User) GetNasen() []Nase {
	var nasen []Nase
	query := "SELECT * FROM nasen WHERE userid = $1 ORDER BY created DESC"
	rows, err := DB.Query(query, user.ID)
	if err != nil {
		return nasen
	}

	for rows.Next() {
		var nase Nase
		err := rows.Scan(&nase.ID, &nase.UserID, &nase.AuthorID, &nase.Reason, &nase.Created)
		if err != nil {
			return nasen
		}
		nasen = append(nasen, nase)
	}
	return nasen
}

func GetUser(id string) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Exp, &user.Level)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetLeaderboard() ([]User, error) {
	var leaderboard []User
	query := "SELECT * FROM users"
	rows, err := DB.Query(query)
	if err != nil {
		return leaderboard, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Exp, &user.Level)
		if err != nil {
			return leaderboard, err
		}
		leaderboard = append(leaderboard, user)
	}

	for i := range leaderboard {
		leaderboard[i].NasenCount = len(leaderboard[i].GetNasen())
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].NasenCount > leaderboard[j].NasenCount
	})

	if len(leaderboard) > 10 {
		leaderboard = leaderboard[:10]
	}

	return leaderboard, nil
}

func (user *User) GiveExp(exp int) error {
	query := "UPDATE users SET exp = $1 WHERE id = $2"
	_, err := DB.Exec(query, user.Exp+exp, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CalcLevel(expNeededPerLevel int) (newLevel int, oldLevel int, levelUp bool, err error) {
	oldLevel = user.Level

	if user.Exp >= oldLevel*expNeededPerLevel {
		newLevel = user.Exp / expNeededPerLevel
		levelUp = true
	} else {
		newLevel = oldLevel
		levelUp = false
	}

	query := "UPDATE users SET level = $1 WHERE id = $2"
	_, err = DB.Exec(query, newLevel, user.ID)
	if err != nil {
		return newLevel, oldLevel, levelUp, err
	}

	return newLevel, oldLevel, levelUp, nil
}
