package drop

import (
	"errors"
	"sync"
)

type dropRegistry struct {
	mutex            sync.Mutex
	dropMap          map[uint32]Drop
	dropLocks        map[uint32]*sync.Mutex
	mapLocks         map[mapKey]*sync.Mutex
	dropsInMap       map[mapKey][]uint32
	dropReservations map[uint32]uint32
}

var registry *dropRegistry
var once sync.Once

var uniqueId = uint32(1000000001)

func GetRegistry() *dropRegistry {
	once.Do(func() {
		registry = &dropRegistry{
			mutex:            sync.Mutex{},
			dropMap:          make(map[uint32]Drop),
			dropLocks:        make(map[uint32]*sync.Mutex),
			mapLocks:         make(map[mapKey]*sync.Mutex),
			dropsInMap:       make(map[mapKey][]uint32),
			dropReservations: make(map[uint32]uint32),
		}
	})
	return registry
}

type mapKey struct {
	worldId   byte
	channelId byte
	mapId     uint32
}

func (d *dropRegistry) CreateDrop(worldId byte, channelId byte, mapId uint32, itemId uint32, equipmentId uint32, quantity uint32,
	mesos uint32, theType byte, x int16, y int16, ownerId uint32, ownerPartyId uint32, dropTime uint64, dropperId uint32,
	dropperX int16, dropperY int16, playerDrop bool, mod byte) Drop {

	mk := mapKey{
		worldId:   worldId,
		channelId: channelId,
		mapId:     mapId,
	}

	d.mutex.Lock()
	ids := existingIds(d.dropMap)
	currentUniqueId := uniqueId
	for contains(ids, currentUniqueId) {
		currentUniqueId = currentUniqueId + 1
		if currentUniqueId > 2000000000 {
			currentUniqueId = 1000000001
		}
		uniqueId = currentUniqueId
	}
	dropMutex := sync.Mutex{}
	d.dropLocks[currentUniqueId] = &dropMutex

	if _, ok := d.mapLocks[mk]; !ok {
		mapMutex := sync.Mutex{}
		d.mapLocks[mk] = &mapMutex
	}

	d.mutex.Unlock()

	dropMutex.Lock()
	drop := Drop{
		id:           currentUniqueId,
		worldId:      worldId,
		channelId:    channelId,
		mapId:        mapId,
		itemId:       itemId,
		equipmentId:  equipmentId,
		quantity:     quantity,
		meso:         mesos,
		dropType:     theType,
		x:            x,
		y:            y,
		ownerId:      ownerId,
		ownerPartyId: ownerPartyId,
		dropTime:     dropTime,
		dropperId:    dropperId,
		dropperX:     dropperX,
		dropperY:     dropperY,
		playerDrop:   playerDrop,
		mod:          mod,
		status:       "AVAILABLE",
	}

	d.dropMap[drop.Id()] = drop
	if lock, ok := d.mapLocks[mk]; ok {
		lock.Lock()
		d.dropsInMap[mk] = append(d.dropsInMap[mk], drop.Id())
		lock.Unlock()
	}
	dropMutex.Unlock()
	return drop
}

func (d *dropRegistry) CancelDropReservation(dropId uint32, characterId uint32) {
	if lock, ok := d.dropLocks[dropId]; ok {
		lock.Lock()

		if _, ok := d.dropMap[dropId]; !ok {
			lock.Unlock()
			return
		}
		if val, ok := d.dropReservations[dropId]; ok {
			if val != characterId {
				lock.Unlock()
				return
			}
		} else {
			lock.Unlock()
			return
		}

		drop := d.dropMap[dropId]
		if drop.Status() != "RESERVED" {
			lock.Unlock()
			return
		}

		drop = drop.CancelReservation()
		d.dropMap[dropId] = drop
		delete(d.dropReservations, dropId)
		lock.Unlock()
	}
}

func (d *dropRegistry) ReserveDrop(dropId uint32, characterId uint32) error {
	if lock, ok := d.dropLocks[dropId]; ok {
		lock.Lock()
		if _, ok := d.dropMap[dropId]; !ok {
			lock.Unlock()
			return errors.New("unable to locate drop")
		}

		drop := d.dropMap[dropId]
		if drop.Status() == "AVAILABLE" {
			drop = drop.Reserve()
			d.dropMap[dropId] = drop
			d.dropReservations[dropId] = characterId
			lock.Unlock()
			return nil
		} else {
			if locker, ok := d.dropReservations[dropId]; ok && locker == characterId {
				lock.Unlock()
				return nil
			} else {
				lock.Unlock()
				return errors.New("reserved by another party")
			}
		}
	} else {
		return errors.New("unable to lock drop")
	}
}

func (d *dropRegistry) RemoveDrop(dropId uint32) (*Drop, error) {
	var drop Drop
	if lock, ok := d.dropLocks[dropId]; ok {
		lock.Lock()
		if drop, ok = d.dropMap[dropId]; ok {
			delete(d.dropMap, dropId)
			mk := mapKey{
				worldId:   drop.WorldId(),
				channelId: drop.ChannelId(),
				mapId:     drop.MapId(),
			}
			if mapLock, ok := d.mapLocks[mk]; ok {
				mapLock.Lock()
				if _, ok := d.dropsInMap[mk]; ok {
					index := indexOf(dropId, d.dropsInMap[mk])
					if index >= 0 && index < len(d.dropsInMap[mk]) {
						d.dropsInMap[mk] = remove(d.dropsInMap[mk], index)
					}
				}
				mapLock.Unlock()
			}
		}
		delete(d.dropReservations, dropId)
		lock.Unlock()
	}
	return &drop, nil
}

func (d *dropRegistry) GetDrop(dropId uint32) (*Drop, error) {
	if lock, ok := d.dropLocks[dropId]; ok {
		lock.Lock()
		if drop, ok := d.dropMap[dropId]; ok {
			lock.Unlock()
			return &drop, nil
		} else {
			lock.Unlock()
			return nil, errors.New("drop not found")
		}
	}
	return nil, errors.New("drop lock not found")
}

func (d *dropRegistry) GetDropsForMap(worldId byte, channelId byte, mapId uint32) ([]Drop, error) {
	mk := mapKey{
		worldId:   worldId,
		channelId: channelId,
		mapId:     mapId,
	}

	if mapLock, ok := d.mapLocks[mk]; ok {
		var drops []Drop
		mapLock.Lock()
		for _, dropId := range d.dropsInMap[mk] {
			if drop, ok := d.dropMap[dropId]; ok {
				drops = append(drops, drop)
			}
		}
		mapLock.Unlock()
		return drops, nil
	}
	return nil, errors.New("map lock not found")
}

func (d *dropRegistry) GetAllDrops() []Drop {
	var drops []Drop
	for _, drop := range d.dropMap {
		drops = append(drops, drop)
	}
	return drops
}

func existingIds(drops map[uint32]Drop) []uint32 {
	var ids []uint32
	for i := range drops {
		ids = append(ids, i)
	}
	return ids
}

func contains(ids []uint32, id uint32) bool {
	for _, element := range ids {
		if element == id {
			return true
		}
	}
	return false
}

func indexOf(uniqueId uint32, data []uint32) int {
	for k, v := range data {
		if uniqueId == v {
			return k
		}
	}
	return -1 //not found.
}

func remove(s []uint32, i int) []uint32 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
