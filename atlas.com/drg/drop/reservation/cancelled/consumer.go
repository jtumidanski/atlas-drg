package cancelled

import (
	"atlas-drg/kafka/consumer"
	"atlas-drg/monster/drop"
	"log"
)

type CancelDropReservationCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func CancelDropReservationCommandCreator() consumer.EmptyEventCreator {
	return func() interface{} {
		return &CancelDropReservationCommand{}
	}
}

func HandleCancelDropReservationCommand() consumer.EventProcessor {
	return func(l *log.Logger, e interface{}) {
		if event, ok := e.(*CancelDropReservationCommand); ok {
			drop.Processor(l).CancelDropReservation(event.DropId, event.CharacterId)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCancelDropReservationCommand]")
		}
	}
}
