package models

// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044& БИ-11.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045& БИ-11.2
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046& БИ-12.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081& БМ-11.1


type Groups struct {
	groups map[string]string
}

func GroupInit() *Groups {
	return &Groups{
		groups: map[string]string{
			"БИ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044&",
			"БИ-11.2":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045&",
			"БИ-12.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046&",
			"БМ-11.1":"http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081&",
		},
	}
}

// Проверяем, существует ли группа, введенная пользователем
func (g *Groups) GroupExists(group string) bool {
	if _, exist := g.groups[group]; exist {
		return true
	}

	return false
}

// Получаем url адрес группы
func (g *Groups) GetGroup(group string) string {
	url, _ := g.groups[group]

	return url
}

