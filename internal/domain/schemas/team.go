package schemas

type Team struct {
    Name    string `json:"team_name" db:"team_name"`
    Members []User `json:"members"`
}
  