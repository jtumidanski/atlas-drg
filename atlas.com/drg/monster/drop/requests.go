package drop

import (
	"atlas-drg/rest/requests"
	"fmt"
)

const (
	dropInformationServicePrefix string = "/ms/dis/"
	dropInformationService              = requests.BaseRequest + dropInformationServicePrefix
	monsterDropResource                 = dropInformationService + "monsters/drops?monsterId=%d"
)

var MonsterDropRequests = func() *monsterDropRequests {
	return &monsterDropRequests{}
}

type monsterDropRequests struct {
}

func (a *monsterDropRequests) GetByMonsterId(monsterId uint32) (*MonsterDropDataContainer, error) {
	ar := &MonsterDropDataContainer{}
	err := requests.Get(fmt.Sprintf(monsterDropResource, monsterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
