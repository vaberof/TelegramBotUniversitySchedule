package storage

type StudyGroupId string
type StudyGroupUrl string

type GroupStorage struct {
	groupStorage map[StudyGroupId]StudyGroupUrl
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		groupStorage: map[StudyGroupId]StudyGroupUrl{
			"БИ-21.1": "https://rasp.sgugit.ru/?ii=1&fi=1&c=2&gn=37&",
			"БИ-21.2": "https://rasp.sgugit.ru/?ii=1&fi=1&c=2&gn=38&",
		},
	}
}

// GetStudyGroupUrl gets study group`s url
func (g *GroupStorage) GetStudyGroupUrl(studyGroupId string) *string {
	url := g.groupStorage[StudyGroupId(studyGroupId)]
	stringUrl := string(url)
	return &stringUrl
}
