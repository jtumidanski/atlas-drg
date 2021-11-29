package drop

import (
	drop2 "atlas-drg/drop"
	"atlas-drg/drop/gathered"
	"atlas-drg/drop/reservation"
	"atlas-drg/equipment"
	"atlas-drg/inventory"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func SpawnDrop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32, dropperX int16, dropperY int16, playerDrop bool, mod byte) {

		dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
		it, _ := inventory.GetInventoryType(itemId)
		var equipmentId uint32
		if it == inventory.TypeValueEquip {
			ro, err := equipment.CreateRandom(l, span)(itemId)
			if err != nil {
				l.WithError(err).Debugf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.", itemId, ownerId)
				return
			}
			eid, err := strconv.Atoi(ro.Data.Id)
			if err != nil {
				l.WithError(err).Debugf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.", itemId, ownerId)
				return
			}
			equipmentId = uint32(eid)
		}

		drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
		drop2.DropEvent(l, span)(worldId, channelId, mapId, drop)
	}
}

func SpawnCharacterDrop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32, dropperX int16, dropperY int16, playerDrop bool, mod byte) {
		dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
		drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
		drop2.DropEvent(l, span)(worldId, channelId, mapId, drop)
	}
}

func CancelDropReservation(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		drop2.GetRegistry().CancelDropReservation(dropId, characterId)
		reservation.DropReservationFailure(l, span)(dropId, characterId)
	}
}

func ReserveDrop(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		err := drop2.GetRegistry().ReserveDrop(dropId, characterId)
		if err == nil {
			l.Debugf("Reserving %d for %d.", dropId, characterId)
			reservation.DropReservationSuccess(l, span)(dropId, characterId)
		} else {
			l.Debugf("Failed reserving %d for %d.", dropId, characterId)
			reservation.DropReservationFailure(l, span)(dropId, characterId)
		}
	}
}

func GatherDrop(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		drop, err := drop2.GetRegistry().RemoveDrop(dropId)
		if err == nil {
			l.Debugf("Gathering %d for %d.", dropId, characterId)
			gathered.DropPickedUp(l, span)(dropId, characterId, drop.MapId())
		}
	}
}

func GetDropById(dropId uint32) (*drop2.Drop, error) {
	return drop2.GetRegistry().GetDrop(dropId)
}

func GetDropsForMap(worldId byte, channelId byte, mapId uint32) ([]*drop2.Drop, error) {
	return drop2.GetRegistry().GetDropsForMap(worldId, channelId, mapId)
}
