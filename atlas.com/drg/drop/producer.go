package drop

import (
	producer2 "atlas-drg/kafka/producer"
	"context"
	"log"
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

func (m *producer) Emit(worldId byte, channelId byte, mapId uint32, drop Drop) {
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
	m.l.Printf("[INFO] dropping item %d in map %d.", drop.ItemId(), drop.MapId())
	producer2.ProduceEvent(m.l, "TOPIC_DROP_EVENT", producer2.CreateKey(int(mapId)), e)
}