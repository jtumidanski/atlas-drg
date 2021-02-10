package gathered

import (
	producer2 "atlas-drg/kafka/producer"
	"context"
	"log"
)

type dropPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	MapId       uint32 `json:"mapId"`
}

var Producer = func(l *log.Logger, ctx context.Context) *producer {
	return &producer{
		l:   l,
		ctx: ctx,
	}
}

type producer struct {
	l   *log.Logger
	ctx context.Context
}

func (m *producer) Emit(dropId uint32, characterId uint32, mapId uint32) {
	e := &dropPickedUpEvent{
		CharacterId: characterId,
		DropId:      dropId,
		MapId:       mapId,
	}
	producer2.ProduceEvent(m.l, "TOPIC_PICKUP_DROP_EVENT", producer2.CreateKey(int(dropId)), e)
}
