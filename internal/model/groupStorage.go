package model

type GroupStorage struct {
	groupStorage map[string]string
}

// NewGroupStorage returns pointer to GroupStorage.
func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		groupStorage: map[string]string{
			"БИ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044&",
			"БИ-11.2":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045&",
			"БИ-12.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046&",
			"БМ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081&",
		},
	}
}

// GroupUrl gets study group`s url.
func (g *GroupStorage) GroupUrl(studyGroupId string) (*string, bool) {
	url, exists :=  g.groupStorage[studyGroupId]
	return &url, exists
}