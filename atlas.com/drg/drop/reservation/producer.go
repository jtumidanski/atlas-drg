package reservation

import (
	producer2 "atlas-drg/kafka/producer"
	"context"
	"log"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

var Producer = func(l *log.Logger, ctx context.Context) *producer {
	return &producer{
		l:   l,
		ctx: ctx,
	}
}

type producer struct {
	l   *log.Logger
	ctx context.Context
}

func (m *producer) EmitFailure(dropId uint32, characterId uint32) {
	m.emit(dropId, characterId, "FAILURE")
}

func (m *producer) emit(dropId uint32, characterId uint32, theType string) {
	e := &dropReservationEvent{
		CharacterId: characterId,
		DropId:      dropId,
		Type:        theType,
	}
	producer2.ProduceEvent(m.l, "TOPIC_DROP_RESERVATION_EVENT", producer2.CreateKey(int(dropId)), e)
}

func (m *producer) EmitSuccess(dropId uint32, characterId uint32) {
	m.emit(dropId, characterId, "SUCCESS")
}
