package rest

import (
	drop2 "atlas-drg/monster/drop"
	"atlas-drg/world"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes)
}

func ProduceRoutes(l logrus.FieldLogger) http.Handler {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/drg").Subrouter()
	router.Use(CommonHeader)

	sRouter := router.PathPrefix("/drops/{id}").Subrouter()
	sRouter.HandleFunc("", drop2.HandleGetDropById(l))

	iRouter := router.PathPrefix("/worlds/{worldId}/channels/{channelId}/maps/{mapId}/drops").Subrouter()
	iRouter.HandleFunc("", world.GetDropsInMap(l))

	return router
}
