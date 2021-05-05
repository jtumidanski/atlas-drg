package spawn

import (
	"atlas-drg/kafka/consumer"
	"atlas-drg/monster/drop"
	"log"
)

type command struct {
	WorldId      byte   `json:"worldId"`
	ChannelId    byte   `json:"channelId"`
	MapId        uint32 `json:"mapId"`
	ItemId       uint32 `json:"itemId"`
	Quantity     uint32 `json:"quantity"`
	Mesos        uint32 `json:"mesos"`
	DropType     byte   `json:"dropType"`
	X            int16  `json:"x"`
	Y            int16  `json:"y"`
	OwnerId      uint32 `json:"ownerId"`
	OwnerPartyId uint32 `json:"ownerPartyId"`
	DropperId    uint32 `json:"dropperId"`
	DropperX     int16  `json:"dropperX"`
	DropperY     int16  `json:"dropperY"`
	PlayerDrop   bool   `json:"playerDrop"`
	Mod          byte   `json:"mod"`
}

func CommandEventCreator() consumer.EmptyEventCreator {
	return func() interface{} {
		return &command{}
	}
}

func HandleCommand() consumer.EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*command); ok {
			drop.Processor(l).SpawnDrop(event.WorldId, event.ChannelId, event.MapId, event.ItemId, event.Quantity,
				event.Mesos, event.DropType, event.X, event.Y, event.OwnerId, event.OwnerPartyId, event.DropperId,
				event.DropperX, event.DropperY, event.PlayerDrop, event.Mod)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [command]")
		}
	}
}