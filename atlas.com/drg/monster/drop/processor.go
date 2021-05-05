package drop

import (
	drop2 "atlas-drg/drop"
	"atlas-drg/drop/gathered"
	"atlas-drg/drop/reservation"
	"context"
	"log"
	"time"
)

var Processor = func(l *log.Logger) *processor {
	return &processor{l: l}
}

type processor struct {
	l *log.Logger
}

func (d *processor) SpawnDrop(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
	drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
	drop2.Producer(d.l, context.Background()).Emit(worldId, channelId, mapId, drop)
}

func (d *processor) CancelDropReservation(dropId uint32, characterId uint32) {
	drop2.GetRegistry().CancelDropReservation(dropId, characterId)
	reservation.Producer(d.l, context.Background()).EmitFailure(dropId, characterId)
}

func (d *processor) ReserveDrop(dropId uint32, characterId uint32) {
	err := drop2.GetRegistry().ReserveDrop(dropId, characterId)
	if err == nil {
		d.l.Printf("[INFO] reserving %d for %d.", dropId, characterId)
		reservation.Producer(d.l, context.Background()).EmitSuccess(dropId, characterId)
	} else {
		d.l.Printf("[INFO] failed reserving %d for %d.", dropId, characterId)
		reservation.Producer(d.l, context.Background()).EmitFailure(dropId, characterId)
	}
}

func (d *processor) GatherDrop(dropId uint32, characterId uint32) {
	drop, err := drop2.GetRegistry().RemoveDrop(dropId)
	if err == nil {
		d.l.Printf("[INFO] gathering %d for %d.", dropId, characterId)
		gathered.Producer(d.l, context.Background()).Emit(dropId, characterId, drop.MapId())
	}
}

func (d *processor) GetDropById(dropId uint32) (*drop2.Drop, error) {
	return drop2.GetRegistry().GetDrop(dropId)
}

func (d *processor) GetDropsForMap(worldId byte, channelId byte, mapId uint32) ([]drop2.Drop, error) {
	return drop2.GetRegistry().GetDropsForMap(worldId, channelId, mapId)
}
