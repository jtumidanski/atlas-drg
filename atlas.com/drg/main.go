package main

import (
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/drop/expired"
	"atlas-drg/kafka/consumers"
	"atlas-drg/logger"
	"atlas-drg/rest"
	tasks "atlas-drg/task"
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	consumers.CreateEventConsumers(l, ctx, wg)

	rest.CreateRestService(l, ctx, wg)

	createTasks(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	drop.ForEachDrop(destroyDrop(l))
	l.Infoln("Service shutdown.")
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
