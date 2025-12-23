package sac

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserName  string `json:"username,omitempty"`
}

func SaveSAC(u User, msg string) error {
	return nil
}
