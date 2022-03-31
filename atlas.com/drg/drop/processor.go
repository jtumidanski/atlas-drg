package drop

import (
	"atlas-drg/equipment"
	"atlas-drg/inventory"
	"atlas-drg/model"
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
			ro, _, err := equipment.CreateRandom(itemId)(l, span)
			if err != nil {
				l.WithError(err).Debugf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.", itemId, ownerId)
				return
			}
			eid, err := strconv.Atoi(ro.Data().Id)
			if err != nil {
				l.WithError(err).Debugf("Generating equipment item %d for character %d, they were not awarded this item. Invalid ID from ESO service.", itemId, ownerId)
				return
			}
			equipmentId = uint32(eid)
		}

		drop := GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
		emitEvent(l, span)(worldId, channelId, mapId, drop)
	}
}

func SpawnCharacterDrop(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32,
	mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) {
	return func(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32, mesos uint32, dropType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropperId uint32, dropperX int16, dropperY int16, playerDrop bool, mod byte) {
		dropTime := uint64(time.Now().UnixNano() / int64(time.Millisecond))
		drop := GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, equipmentId, quantity, mesos, dropType, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX, dropperY, playerDrop, mod)
		emitEvent(l, span)(worldId, channelId, mapId, drop)
	}
}

func CancelReservation(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		GetRegistry().CancelDropReservation(dropId, characterId)
		emitReservationFailure(l, span)(dropId, characterId)
	}
}

func Reserve(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		err := GetRegistry().ReserveDrop(dropId, characterId)
		if err == nil {
			l.Debugf("Reserving %d for %d.", dropId, characterId)
			emitReservationSuccess(l, span)(dropId, characterId)
		} else {
			l.Debugf("Failed reserving %d for %d.", dropId, characterId)
			emitReservationFailure(l, span)(dropId, characterId)
		}
	}
}

func Gather(l logrus.FieldLogger, span opentracing.Span) func(dropId uint32, characterId uint32) {
	return func(dropId uint32, characterId uint32) {
		drop, err := GetRegistry().RemoveDrop(dropId)
		if err == nil {
			l.Debugf("Gathering %d for %d.", dropId, characterId)
			emitPickedUp(l, span)(dropId, characterId, drop.MapId())
		}
	}
}

func GetById(dropId uint32) (Drop, error) {
	return GetRegistry().GetDrop(dropId)
}

func GetForMap(worldId byte, channelId byte, mapId uint32) ([]Drop, error) {
	return GetRegistry().GetDropsForMap(worldId, channelId, mapId)
}

func AllModelProvider() model.SliceProvider[Drop] {
	return model.FixedSliceProvider[Drop](GetRegistry().GetAllDrops())
}

func Destroy(l logrus.FieldLogger, span opentracing.Span) model.Operator[Drop] {
	return func(d Drop) error {
		rd, err := GetRegistry().RemoveDrop(d.Id())
		if err != nil {
			l.WithError(err).Errorf("Unable to destroy drop %d.", rd.Id())
			return err
		}
		emitExpiredEvent(l, span)(rd.WorldId(), rd.ChannelId(), rd.MapId(), rd.Id())
		return nil
	}
}
