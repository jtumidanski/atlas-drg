package consumers

import (
	spawn2 "atlas-drg/character/drop/spawn"
	"atlas-drg/drop/gather"
	"atlas-drg/drop/reservation/cancelled"
	"atlas-drg/drop/reservation/reserved"
	"atlas-drg/drop/spawn"
	"atlas-drg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	SpawnDropCommand             = "spawn_drop_command"
	CancelDropReservationCommand = "cancel_drop_reservation_command"
	ReserveDropCommand           = "reserve_drop_command"
	PickupDropCommand            = "pickup_drop_command"
	SpawnCharacterDropCommand    = "spawn_character_drop_command"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_SPAWN_DROP_COMMAND", SpawnDropCommand, spawn.CommandEventCreator(), spawn.HandleCommand())
	cec("TOPIC_CANCEL_DROP_RESERVATION_COMMAND", CancelDropReservationCommand, cancelled.CancelDropReservationCommandCreator(), cancelled.HandleCancelDropReservationCommand())
	cec("TOPIC_RESERVE_DROP_COMMAND", ReserveDropCommand, reserved.ReserveDropCommandCreator(), reserved.HandleReserveDropCommand())
	cec("TOPIC_PICKUP_DROP_COMMAND", PickupDropCommand, gather.GatherDropCommandCreator(), gather.HandleGatherDropCommand())
	cec("TOPIC_SPAWN_CHARACTER_DROP_COMMAND", SpawnCharacterDropCommand, spawn2.CommandEventCreator(), spawn2.HandleCommand())

}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "Drop Registry", emptyEventCreator, processor)
}
