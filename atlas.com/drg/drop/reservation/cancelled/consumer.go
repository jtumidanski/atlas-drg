package cancelled

import (
	"atlas-drg/kafka/handler"
	"atlas-drg/monster/drop"
	"github.com/sirupsen/logrus"
)

type CancelDropReservationCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func CancelDropReservationCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &CancelDropReservationCommand{}
	}
}

func HandleCancelDropReservationCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*CancelDropReservationCommand); ok {
			drop.Processor(l).CancelDropReservation(event.DropId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler.")
		}
	}
}
