package drop

import (
	"atlas-drg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type dropEvent struct {
	WorldId         byte   `json:"worldId"`
	ChannelId       byte   `json:"channelId"`
	MapId           uint32 `json:"mapId"`
	UniqueId        uint32 `json:"uniqueId"`
	ItemId          uint32 `json:"itemId"`
	Quantity        uint32 `json:"quantity"`
	Meso            uint32 `json:"meso"`
	DropType        byte   `json:"dropType"`
	DropX           int16  `json:"dropX"`
	DropY           int16  `json:"dropY"`
	OwnerId         uint32 `json:"ownerId"`
	OwnerPartyId    uint32 `json:"ownerPartyId"`
	DropTime        uint64 `json:"dropTime"`
	DropperUniqueId uint32 `json:"dropperUniqueId"`
	DropperX        int16  `json:"dropperX"`
	DropperY        int16  `json:"dropperY"`
	PlayerDrop      bool   `json:"playerDrop"`
	Mod             byte   `json:"mod"`
}

func emitEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, drop *Drop) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_DROP_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, drop *Drop) {
		e := &dropEvent{
			WorldId:         worldId,
			ChannelId:       channelId,
			MapId:           mapId,
			UniqueId:        drop.Id(),
			ItemId:          drop.ItemId(),
			Quantity:        drop.Quantity(),
			Meso:            drop.Meso(),
			DropType:        drop.Type(),
			DropX:           drop.X(),
			DropY:           drop.Y(),
			OwnerId:         drop.OwnerId(),
			OwnerPartyId:    drop.OwnerPartyId(),
			DropTime:        drop.DropTime(),
			DropperUniqueId: drop.DropperId(),
			DropperX:        drop.DropperX(),
			DropperY:        drop.DropperY(),
			PlayerDrop:      drop.PlayerDrop(),
			Mod:             drop.Mod(),
		}
		l.Debugf("Dropping item %d in map %d.", drop.ItemId(), drop.MapId())
		producer(kafka.CreateKey(int(mapId)), e)
	}
}

type dropExpiredEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
}

func emitExpiredEvent(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, id uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_DROP_EXPIRE_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, id uint32) {
		e := &dropExpiredEvent{
			WorldId:   worldId,
			ChannelId: channelId,
			MapId:     mapId,
			UniqueId:  id,
		}
		producer(kafka.CreateKey(int(mapId)), e)
	}
}

type dropPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	MapId       uint32 `json:"mapId"`
}

func emitPickedUp(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32, mapId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_PICKUP_DROP_EVENT")
	return func(dropId uint32, characterId uint32, mapId uint32) {
		e := &dropPickedUpEvent{
			CharacterId: characterId,
			DropId:      dropId,
			MapId:       mapId,
		}
		producer(kafka.CreateKey(int(dropId)), e)
	}
}

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func emitReservationFailure(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_DROP_RESERVATION_EVENT")
	return func(dropId uint32, characterId uint32) {
		emitReservation(producer, dropId, characterId, "FAILURE")
	}
}

func emitReservation(producer func(key []byte, event interface{}), characterId uint32, dropId uint32, theType string) {
	e := &dropReservationEvent{
		CharacterId: characterId,
		DropId:      dropId,
		Type:        theType,
	}
	producer(kafka.CreateKey(int(dropId)), e)
}

func emitReservationSuccess(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_DROP_RESERVATION_EVENT")
	return func(dropId uint32, characterId uint32) {
		emitReservation(producer, characterId, dropId, "SUCCESS")
	}
}
