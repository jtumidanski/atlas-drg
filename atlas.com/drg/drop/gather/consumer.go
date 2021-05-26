package gather

import (
	"atlas-drg/kafka/handler"
	"atlas-drg/monster/drop"
	"github.com/sirupsen/logrus"
)

type GatherDropCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func GatherDropCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &GatherDropCommand{}
	}
}

func HandleGatherDropCommand() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*GatherDropCommand); ok {
			drop.GatherDrop(l)(event.DropId, event.CharacterId)
		} else {
			l.Errorf("Unable to cast event provided to handler.")
		}
	}
}
