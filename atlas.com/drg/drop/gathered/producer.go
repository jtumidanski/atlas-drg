package gathered

import (
	producer2 "atlas-drg/kafka/producer"
	"github.com/sirupsen/logrus"
)

type dropPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	MapId       uint32 `json:"mapId"`
}

func DropPickedUp(l logrus.FieldLogger) func(dropId uint32, characterId uint32, mapId uint32) {
	producer := producer2.ProduceEvent(l, "TOPIC_PICKUP_DROP_EVENT")
	return func(dropId uint32, characterId uint32, mapId uint32) {
		e := &dropPickedUpEvent{
			CharacterId: characterId,
			DropId:      dropId,
			MapId:       mapId,
		}
		producer(producer2.CreateKey(int(dropId)), e)
	}
}
