package expired

import (
	producer2 "atlas-drg/kafka/producer"
	"context"
	"log"
)

type dropExpiredEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
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

func (m *producer) Emit(worldId byte, channelId byte, mapId uint32, id uint32) {
	e := &dropExpiredEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		MapId:     mapId,
		UniqueId:  id,
	}
	producer2.ProduceEvent(m.l, "TOPIC_DROP_EXPIRE_EVENT", producer2.CreateKey(int(mapId)), e)
}
