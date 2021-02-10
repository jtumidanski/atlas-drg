package drop

import (
	"atlas-drg/rest/requests"
)

type MonsterDropDataContainer struct {
	data     requests.DataSegment
	included requests.DataSegment
}

type MonsterDropData struct {
	Id         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes MonsterDropAttributes `json:"attributes"`
}

type MonsterDropAttributes struct {
	MonsterId       uint32 `json:"monsterId"`
	ItemId          uint32 `json:"itemId"`
	MaximumQuantity uint32 `json:"maximumQuantity"`
	MinimumQuantity uint32 `json:"minimumQuantity"`
	Chance          uint32 `json:"chance"`
}

func (a *MonsterDropDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := requests.UnmarshalRoot(data, requests.MapperFunc(EmptyMonsterDropData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}


func (a *MonsterDropDataContainer) Data() *MonsterDropData {
	if len(a.data) >= 1 {
		return a.data[0].(*MonsterDropData)
	}
	return nil
}

func (a *MonsterDropDataContainer) DataList() []MonsterDropData {
	var r = make([]MonsterDropData, 0)
	for _, x := range a.data {
		r = append(r, *x.(*MonsterDropData))
	}
	return r
}

func EmptyMonsterDropData() interface{} {
	return &MonsterDropData{}
}