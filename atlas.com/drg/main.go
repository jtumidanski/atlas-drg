package main

import (
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/drop/expired"
	"atlas-drg/kafka/consumer"
	"atlas-drg/logger"
	"atlas-drg/rest"
	tasks "atlas-drg/task"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

	consumer.CreateEventConsumers(l)

	rest.CreateRestService(l)
	createTasks(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)

	drop.ForEachDrop(destroyDrop(l))
}

func destroyDrop(l logrus.FieldLogger) drop.DropOperator {
	return func(d *drop.Drop) {
		drop.GetRegistry().RemoveDrop(d.Id())
		expired.DropExpired(l)(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
	}
}

func createTasks(l logrus.FieldLogger) {
	c, err := registries.GetConfiguration()
	if err != nil {
		return
	}

	go tasks.Register(expired.NewDropExpiration(l, c.ItemExpireCheck))
}
