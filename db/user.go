package db

type DBUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	NasenCount int    `json:"nasen_count"`
	Exp        int    `json:"exp"`
	Level      int    `json:"level"`
	RiotPUUID  string `json:"riot_puuid"`
}

func InsertDBUser(discordUserID string, username string) error {
	query := "INSERT INTO users (id, name) VALUES ($1, $2)"
	_, err := DB.Exec(query, discordUserID, username)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(id string) (DBUser, error) {
	var user DBUser
	query := "SELECT * FROM users WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Exp, &user.Level, &user.RiotPUUID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UserInDatabase(id string) bool {
	query := "SELECT id FROM users WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&id)

	return err == nil
}
