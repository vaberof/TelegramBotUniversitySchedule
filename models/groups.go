package models

// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044& БИ-11.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045& БИ-11.2
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046& БИ-12.1
// http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081& БМ-11.1

var Groups map[string]string

// Проверяем, существует ли группа
func GroupExists(group string) bool {

	Groups = make(map[string]string)
	Groups["БИ-11.1"] = "http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1044&"
	Groups["БИ-11.2"] = "http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1045&"
	Groups["БИ-12.1"] = "http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1046&"
	Groups["БМ-11.1"] = "http://rasp.sgugit.ru/?ii=1&fi=1&c=1&gn=1081&"


	if _, exist := Groups[group]; exist {
		return true
	} else {
		return false
	}
	return true
}

