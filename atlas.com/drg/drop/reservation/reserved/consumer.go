package reserved

import (
	"atlas-drg/kafka/handler"
	"atlas-drg/monster/drop"
	"github.com/opentracing/opentracing-go"
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
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*ReserveDropCommand); ok {
			drop.ReserveDrop(l, span)(event.DropId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler.")
		}
	}
}
