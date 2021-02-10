package reserved

import (
	"atlas-drg/kafka/consumer"
	"atlas-drg/monster/drop"
	"log"
)

type ReserveDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func ReserveDropCommandCreator() consumer.EmptyEventCreator {
	return func() interface{} {
		return &ReserveDropCommand{}
	}
}

func HandleReserveDropCommand() consumer.EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*ReserveDropCommand); ok {
			drop.Processor(l).ReserveDrop(event.DropId, event.CharacterId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCancelDropReservationCommand]")
		}
	}
}
