package death

import (
	"atlas-drg/kafka/consumer"
	"atlas-drg/monster/drop"
	"log"
)

type DamageEntry struct {
	CharacterId uint32 `json:"characterId"`
	Damage      uint64 `json:"damage"`
}

type MonsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         uint32        `json:"mapId"`
	UniqueId      uint32        `json:"uniqueId"`
	MonsterId     uint32        `json:"monsterId"`
	X             int16         `json:"x"`
	Y             int16         `json:"y"`
	KillerId      uint32        `json:"killerId"`
	DamageEntries []DamageEntry `json:"damageEntries"`
}

func MonsterKilledEventCreator() consumer.EmptyEventCreator {
	return func() interface{} {
		return &MonsterKilledEvent{}
	}
}

func HandleMonsterKilledEvent() consumer.EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*MonsterKilledEvent); ok {
			l.Printf("[INFO] processing death of %d in map %d.", event.MonsterId, event.MapId)
			drop.Processor(l).CreateDrops(event.WorldId, event.ChannelId, event.MapId, event.UniqueId, event.MonsterId, event.X, event.Y, event.KillerId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMonsterKilledEvent]")
		}
	}
}
