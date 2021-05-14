package reserved

import (
	"atlas-drg/kafka/handler"
	"atlas-drg/monster/drop"
	"github.com/sirupsen/logrus"
)

type ReserveDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func ReserveDropCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &ReserveDropCommand{}
	}
}

func HandleReserveDropCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*ReserveDropCommand); ok {
			drop.ReserveDrop(l)(event.DropId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler.")
		}
	}
}
