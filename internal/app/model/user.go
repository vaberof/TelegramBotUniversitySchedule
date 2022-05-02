package model

type User struct {
	Data map[int64]string // includes user`s chat id and his input group id
}

func NewUser() *User {
	return &User{
		Data: map[int64]string{},
	}
}

// AddData adds user`s unique chat id
// and his input group id to User.
func (u *User) AddData(userChatID int64, groupID string) {
	u.Data[userChatID] = groupID
}
