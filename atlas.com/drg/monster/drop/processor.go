package drop

import (
	drop2 "atlas-drg/drop"
	"atlas-drg/drop/gathered"
	"atlas-drg/drop/reservation"
	"atlas-drg/map/point"
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var Processor = func(l *log.Logger) *processor {
	return &processor{l: l}
}

type processor struct {
	l *log.Logger
}

func (d *processor) CreateDrops(worldId byte, channelId byte, mapId uint32, monsterUniqueId uint32, monsterId uint32, x int16, y int16, killerId uint32) {
	// TODO determine type of drop
	//    monster is explosive? 3
	//    monster has ffa loot? 2
	//    killer is in party? 1
	dropType := byte(0)

	ns, err := d.GetDropsForMonster(monsterId)
	if err != nil {
		return
	}

	d.l.Printf("[INFO] successfully found %d drops to evaluate.", len(ns))

	ns = d.getSuccessfulDrops(ns, killerId)

	d.l.Printf("[INFO] successfully found %d drops to emit.", len(ns))

	for i, drop := range ns {
		d.createDrop(worldId, channelId, mapId, i+1, monsterUniqueId, x, y, killerId, dropType, drop)
	}
}

func (d *processor) createDrop(worldId byte, channelId byte, mapId uint32, index int, uniqueId uint32, x int16, y int16, killerId uint32, dropType byte, drop MonsterDrop) {
	factor := 0
	if dropType == 3 {
		factor = 40
	} else {
		factor = 25
	}
	newX := x
	if index%2 == 0 {
		newX += int16(factor * ((index + 1) / 2))
	} else {
		newX += int16(-(factor * (index / 2)))
	}
	if drop.ItemId() == 0 {
		d.spawnMeso(worldId, channelId, mapId, uniqueId, x, y, killerId, dropType, drop, newX, y)
	} else {
		d.spawnItem(worldId, channelId, mapId, drop.ItemId(), uniqueId, x, y, killerId, dropType, drop, newX, y)
	}

}

func (d *processor) spawnItem(worldId byte, channelId byte, mapId uint32, itemId uint32, uniqueId uint32, x int16, y int16, killerId uint32, dropType byte, drop MonsterDrop, posX int16, posY int16) {
	quantity := uint32(1)
	if drop.MaximumQuantity() != 1 {
		quantity = uint32(rand.Int31n(int32(drop.MaximumQuantity()-drop.MinimumQuantity()))) + drop.MinimumQuantity()
	}
	d.spawnDrop(worldId, channelId, mapId, itemId, quantity, 0, posX, posY, x, y, uniqueId, killerId, false, dropType)
}

func (d *processor) spawnMeso(worldId byte, channelId byte, mapId uint32, uniqueId uint32, x int16, y int16, killerId uint32, dropType byte, drop MonsterDrop, posX int16, posY int16) {
	mesos := uint32(rand.Int31n(int32(drop.MaximumQuantity()-drop.MinimumQuantity()))) + drop.MinimumQuantity()
	//TODO apply characters meso buff.
	d.spawnDrop(worldId, channelId, mapId, 0, 0, mesos, posX, posY, x, y, uniqueId, killerId, false, dropType)
}

func (d *processor) spawnDrop(worldId byte, channelId byte, mapId uint32, itemId uint32, quantity uint32, mesos uint32, posX int16, posY int16, monsterX int16, monsterY int16, uniqueId uint32, killerId uint32, playerDrop bool, dropType byte) {
	tempX, tempY := d.calculateDropPosition(mapId, posX, posY, monsterX, monsterY)
	tempX, tempY = d.calculateDropPosition(mapId, tempX, tempY, tempX, tempY)
	drop := drop2.GetRegistry().CreateDrop(worldId, channelId, mapId, itemId, quantity, mesos, dropType, tempX, tempY, killerId, 0, uint64(time.Now().UnixNano() / int64(time.Millisecond)), uniqueId, monsterX, monsterY, playerDrop, byte(1))
	drop2.Producer(d.l, context.Background()).Emit(worldId, channelId, mapId, drop)
}

func (d *processor) calculateDropPosition(mapId uint32, initialX int16, initialY int16, fallbackX int16, fallbackY int16) (int16, int16) {
	resp, err := point.MapInformationRequests().CalculateDropPosition(mapId, initialX, initialY, fallbackX, fallbackY)
	if err != nil {
		return fallbackX, fallbackY
	} else {
		return resp.Data().Attributes.X, resp.Data().Attributes.Y
	}
}

func (d *processor) getSuccessfulDrops(ns []MonsterDrop, killerId uint32) []MonsterDrop {
	rs := make([]MonsterDrop, 0)
	for _, drop := range ns {
		if d.evaluateSuccess(killerId, drop) {
			rs = append(rs, drop)
		}
	}
	return rs
}

func (d *processor) evaluateSuccess(killerId uint32, drop MonsterDrop) bool {
	//TODO evaluate card rate for killer, whether it's meso or drop.
	chance := int32(math.Min(float64(drop.Chance()*1), math.MaxUint32))
	return rand.Int31n(999999) < chance
}

func (d *processor) GetDropsForMonster(monsterId uint32) ([]MonsterDrop, error) {
	rest, err := MonsterDropRequests().GetByMonsterId(monsterId)
	if err != nil {
		return nil, err
	}

	ns := make([]MonsterDrop, 0)
	for _, drop := range rest.DataList() {
		id, err := strconv.ParseUint(drop.Id, 10, 32)
		if err != nil {
			break
		}
		n := d.makeDrop(uint32(id), drop.Attributes)
		ns = append(ns, n)
	}
	return ns, nil
}

func (d *processor) makeDrop(id uint32, att MonsterDropAttributes) MonsterDrop {
	return MonsterDrop{
		monsterId:       att.MonsterId,
		itemId:          att.ItemId,
		minimumQuantity: att.MinimumQuantity,
		maximumQuantity: att.MaximumQuantity,
		chance:          att.Chance,
	}
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
