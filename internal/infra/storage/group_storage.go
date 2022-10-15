package storage

type GroupStorage struct {
	Groups []*Group
}

type Group struct {
	Id         string
	Name       string
	ExternalId string
}

func NewGroupStorage() *GroupStorage {
	return &GroupStorage{
		Groups: []*Group{
			{
				Id:         "БИ",
				Name:       "21.1",
				ExternalId: "?ii=1&fi=1&c=2&gn=37&",
			},
			{
				Id:         "БИ",
				Name:       "21.2",
				ExternalId: "?ii=1&fi=1&c=2&gn=38&",
			},
		},
	}
}

// GetStudyGroup gets study group object.
func (g *GroupStorage) GetStudyGroup(groupId string) *Group {
	for i := 0; i < len(g.Groups); i++ {
		group := *g.Groups[i]
		if group.Id+"-"+group.Name == groupId {
			return &group
		}
	}
	return nil
}
