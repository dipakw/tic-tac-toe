package common

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token,omitempty"`
}

func (u *User) WithoutToken() *User {
	return &User{
		ID:   u.ID,
		Name: u.Name,
	}
}
