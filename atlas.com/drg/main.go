package main

import (
	spawn2 "atlas-drg/character/drop/spawn"
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/drop/expired"
	"atlas-drg/drop/gather"
	"atlas-drg/drop/reservation/cancelled"
	"atlas-drg/drop/reservation/reserved"
	"atlas-drg/drop/spawn"
	"atlas-drg/kafka/consumer"
	"atlas-drg/rest"
	tasks "atlas-drg/task"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := log.New(os.Stdout, "drg ", log.LstdFlags|log.Lmicroseconds)

	createEventConsumers(l)
	createRestService(l)
	createTasks(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Println("[INFO] shutting down via signal:", sig)

	drop.ForEachDrop(destroyDrop(l))
}

func destroyDrop(l *log.Logger) drop.DropOperator {
	return func(d drop.Drop) {
		drop.GetRegistry().RemoveDrop(d.Id())
		expired.Producer(l, context.Background()).Emit(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
	}
}

func createTasks(l *log.Logger) {
	c, err := registries.GetConfiguration()
	if err != nil {
		return
	}

	go tasks.Register(expired.NewDropExpiration(l, c.ItemExpireCheck))
}

func createEventConsumers(l *log.Logger) {
	createEventConsumer(l, "TOPIC_SPAWN_DROP_COMMAND", spawn.CommandEventCreator(), spawn.HandleCommand())
	createEventConsumer(l, "TOPIC_CANCEL_DROP_RESERVATION_COMMAND", cancelled.CancelDropReservationCommandCreator(), cancelled.HandleCancelDropReservationCommand())
	createEventConsumer(l, "TOPIC_RESERVE_DROP_COMMAND", reserved.ReserveDropCommandCreator(), reserved.HandleReserveDropCommand())
	createEventConsumer(l, "TOPIC_PICKUP_DROP_COMMAND", gather.GatherDropCommandCreator(), gather.HandleGatherDropCommand())
	createEventConsumer(l, "TOPIC_SPAWN_CHARACTER_DROP_COMMAND", spawn2.CommandEventCreator(), spawn2.HandleCommand())
}

func createRestService(l *log.Logger) {
	rs := rest.NewServer(l)
	go rs.Run()
}

func createEventConsumer(l *log.Logger, topicToken string, emptyEventCreator consumer.EmptyEventCreator, eventProcessor consumer.EventProcessor) {
	c := consumer.NewConsumer(l, context.Background(), eventProcessor,
		consumer.SetGroupId("Drop Registry"),
		consumer.SetTopicToken(topicToken),
		consumer.SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}
