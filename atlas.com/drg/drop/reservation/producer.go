package reservation

import (
	producer2 "atlas-drg/kafka/producer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func DropReservationFailure(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	producer := producer2.ProduceEvent(l, span, "TOPIC_DROP_RESERVATION_EVENT")
	return func(dropId uint32, characterId uint32) {
		emitReservation(producer, dropId, characterId, "FAILURE")
	}
}

func emitReservation(producer func(key []byte, event interface{}), characterId uint32, dropId uint32, theType string) {
	e := &dropReservationEvent{
		CharacterId: characterId,
		DropId:      dropId,
		Type:        theType,
	}
	producer(producer2.CreateKey(int(dropId)), e)
}

func DropReservationSuccess(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	producer := producer2.ProduceEvent(l, span, "TOPIC_DROP_RESERVATION_EVENT")
	return func(dropId uint32, characterId uint32) {
		emitReservation(producer, characterId, dropId, "SUCCESS")
	}
}
