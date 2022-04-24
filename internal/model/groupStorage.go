package model

type GroupStorage struct {
	groupStorage map[string]string
}

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

// StudyGroup gets study group`s url
// and check if study group exists.
func (g *GroupStorage) StudyGroup(studyGroupID string) (*string, bool) {
	url, exists :=  g.groupStorage[studyGroupID]
	return &url, exists
}