package storage

type MessageStorage struct {
	MessageData map[int64]string // includes user`s chat id and his input group id
}

func NewMessageStorage() *MessageStorage {
	return &MessageStorage{
		MessageData: map[int64]string{},
	}
}

// AddMessageData adds user`s unique chat id
// and his input group id to MessageStorage.
func (u *MessageStorage) AddMessageData(userChatID int64, groupID string) {
	u.MessageData[userChatID] = groupID
}
