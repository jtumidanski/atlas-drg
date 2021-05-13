package consumer

import (
	spawn2 "atlas-drg/character/drop/spawn"
	"atlas-drg/drop/gather"
	"atlas-drg/drop/reservation/cancelled"
	"atlas-drg/drop/reservation/reserved"
	"atlas-drg/drop/spawn"
	"atlas-drg/kafka/handler"
	"github.com/sirupsen/logrus"
)

func CreateEventConsumers(l *logrus.Logger) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_SPAWN_DROP_COMMAND", spawn.CommandEventCreator(), spawn.HandleCommand())
	cec( "TOPIC_CANCEL_DROP_RESERVATION_COMMAND", cancelled.CancelDropReservationCommandCreator(), cancelled.HandleCancelDropReservationCommand())
	cec( "TOPIC_RESERVE_DROP_COMMAND", reserved.ReserveDropCommandCreator(), reserved.HandleReserveDropCommand())
	cec( "TOPIC_PICKUP_DROP_COMMAND", gather.GatherDropCommandCreator(), gather.HandleGatherDropCommand())
	cec("TOPIC_SPAWN_CHARACTER_DROP_COMMAND", spawn2.CommandEventCreator(), spawn2.HandleCommand())

}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	go NewConsumer(l, topicToken, "Drop Registry", emptyEventCreator, processor)
}
