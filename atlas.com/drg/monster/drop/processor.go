package drop

import (
	drop2 "atlas-drg/drop"
	"atlas-drg/drop/gathered"
	"atlas-drg/drop/reservation"
	"atlas-drg/equipment"
	"atlas-drg/inventory"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var Processor = func(l logrus.FieldLogger) *processor {
	return &processor{l: l}
}

type processor struct {
	l logrus.FieldLogger
}

func (d *processor) SpawnDrop(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
	it, _ := inventory.GetInventoryType(itemId)
	var equipmentId uint32
	if it == inventory.TypeValueEquip {
		ro, err := equipment.CreateRandom(itemId)
		if err != nil {
			d.l.Debugf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
			return
		}
		eid, err := strconv.Atoi(ro.Data.Id)
		if err != nil {
			d.l.Debugf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.")
			return
		}
		equipmentId = uint32(eid)
	}

	drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
	drop2.DropEvent(d.l)(worldId, channelId, mapId, drop)
}

func (d *processor) SpawnCharacterDrop(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
	drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
	drop2.DropEvent(d.l)(worldId, channelId, mapId, drop)
}

func (d *processor) CancelDropReservation(dropId uint32, characterId uint32) {
	drop2.GetRegistry().CancelDropReservation(dropId, characterId)
	reservation.DropReservationFailure(d.l)(dropId, characterId)
}

func (d *processor) ReserveDrop(dropId uint32, characterId uint32) {
	err := drop2.GetRegistry().ReserveDrop(dropId, characterId)
	if err == nil {
		d.l.Debugf("Reserving %d for %d.", dropId, characterId)
		reservation.DropReservationSuccess(d.l)(dropId, characterId)
	} else {
		d.l.Debugf("Failed reserving %d for %d.", dropId, characterId)
		reservation.DropReservationFailure(d.l)(dropId, characterId)
	}
}

func (d *processor) GatherDrop(dropId uint32, characterId uint32) {
	drop, err := drop2.GetRegistry().RemoveDrop(dropId)
	if err == nil {
		d.l.Debugf("Gathering %d for %d.", dropId, characterId)
		gathered.DropPickedUp(d.l)(dropId, characterId, drop.MapId())
	}
}

func (d *processor) GetDropById(dropId uint32) (*drop2.Drop, error) {
	return drop2.GetRegistry().GetDrop(dropId)
}

func (d *processor) GetDropsForMap(worldId byte, channelId byte, mapId uint32) ([]drop2.Drop, error) {
	return drop2.GetRegistry().GetDropsForMap(worldId, channelId, mapId)
}
