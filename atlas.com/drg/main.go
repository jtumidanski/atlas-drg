package main

import (
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/drop/expired"
	"atlas-drg/kafka/consumers"
	"atlas-drg/logger"
	drop2 "atlas-drg/monster/drop"
	"atlas-drg/rest"
	tasks "atlas-drg/task"
	"atlas-drg/tracing"
	"atlas-drg/world"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-drg"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	consumers.CreateEventConsumers(l, ctx, wg)

	rest.CreateService(l, ctx, wg, "/ms/drg", drop2.InitResource, world.InitResource)

	createTasks(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()

	span := opentracing.StartSpan("shutdown")
	drop.ForEachDrop(destroyDrop(l, span))
	span.Finish()
	l.Infoln("Service shutdown.")
}

func destroyDrop(l logrus.FieldLogger, span opentracing.Span) drop.DropOperator {
	return func(d *drop.Drop) {
		drop.GetRegistry().RemoveDrop(d.Id())
		expired.DropExpired(l, span)(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
	}
}

func createTasks(l logrus.FieldLogger) {
	c, err := registries.GetConfiguration()
	if err != nil {
		return
	}

	go tasks.Register(expired.NewDropExpiration(l, c.ItemExpireCheck))
}
