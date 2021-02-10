package gather

import (
	"atlas-drg/kafka/consumer"
	"atlas-drg/monster/drop"
	"log"
)

type GatherDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func GatherDropCommandCreator() consumer.EmptyEventCreator {
	return func() interface{} {
		return &GatherDropCommand{}
	}
}

func HandleGatherDropCommand() consumer.EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*GatherDropCommand); ok {
			drop.Processor(l).GatherDrop(event.DropId, event.CharacterId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleGatherDropCommand]")
		}
	}
}
