package model

type User struct {
	Data map[int64]string // includes user`s chat id and his input group id
}

func NewUser() *User {
	return &User{
		Data: map[int64]string{},
	}
}

// AddUser adds new user to User with unique chat id
// and sets user input group id
func (u *User) AddUser(userChatID int64, groupID string) {
	u.Data[userChatID] = groupID
}
