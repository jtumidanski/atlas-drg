package expired

import (
	producer2 "atlas-drg/kafka/producer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type dropExpiredEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
}

func DropExpired(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, id uint32) {
	producer := producer2.ProduceEvent(l, span, "TOPIC_DROP_EXPIRE_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, id uint32) {
		e := &dropExpiredEvent{
			WorldId:   worldId,
			ChannelId: channelId,
			MapId:     mapId,
			UniqueId:  id,
		}
		producer(producer2.CreateKey(int(mapId)), e)
	}
}
