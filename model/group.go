package model

// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044& БИ-11.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045& БИ-11.2
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046& БИ-12.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081& БМ-11.1

type GroupStorage struct {
	groups map[string]string
	savedGroup string
}

func CreateGroupStorage() *GroupStorage {
	return &GroupStorage{
		groups: map[string]string{
			"БИ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044&",
			"БИ-11.2":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045&",
			"БИ-12.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046&",
			"БМ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081&",
		},
		savedGroup: "",
	}
}

// Получаем url адрес группы
func (g *GroupStorage) GetGroupUrl(group string) string {
	url, _ := g.groups[group]

	return url
}

// Проверяем, существует ли группа, введенная пользователем
func (g *GroupStorage) Exists(group string) bool {
	if _, exist := g.groups[group]; exist {
		return true
	}

	return false
}

// Проверяем, существует ли группа, введенная пользователем
func (g *GroupStorage) isSaved(group string) bool {
	if g.savedGroup != "" {
		return true
	}

	return false
}
