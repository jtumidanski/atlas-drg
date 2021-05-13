package rest

import (
	drop2 "atlas-drg/monster/drop"
	"atlas-drg/world"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	l  *logrus.Logger
	hs *http.Server
}

func NewServer(l *logrus.Logger) *Server {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/drg").Subrouter()
	router.Use(commonHeader)

	sRouter := router.PathPrefix("/drops/{id}").Subrouter()
	sRouter.HandleFunc("", drop2.GetDropById(l))

	iRouter := router.PathPrefix("/worlds/{worldId}/channels/{channelId}/maps/{mapId}/drops").Subrouter()
	iRouter.HandleFunc("", world.GetDropsInMap(l))

	w := l.Writer()
	defer w.Close()

	hs := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     log.New(w, "", 0), // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read requests from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	return &Server{l, &hs}
}

func (s *Server) Run() {
	s.l.Infoln("Starting server on port 8080")
	err := s.hs.ListenAndServe()
	if err != nil {
		s.l.WithError(err).Errorf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
