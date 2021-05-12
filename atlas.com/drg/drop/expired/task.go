package expired

import (
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/equipment"
	"context"
	"log"
	"time"
)

type DropExpiration struct {
	l        *log.Logger
	interval uint64
}

func NewDropExpiration(l *log.Logger, interval uint64) *DropExpiration {
	return &DropExpiration{l, interval}
}

func (r *DropExpiration) Run() {
	var expire uint64

	c, err := registries.GetConfiguration()
	if err != nil {
		expire = 180000
	} else {
		expire = c.ItemExpireInterval
	}

	ds := drop.GetRegistry().GetAllDrops()
	for _, d := range ds {
		if d.Status() == "AVAILABLE" {
			if d.DropTime()+expire < uint64(time.Now().UnixNano()/int64(time.Millisecond)) {
				_, err := drop.GetRegistry().RemoveDrop(d.Id())
				if err != nil {
					r.l.Printf("Unable to remove drop from registry.")
					continue
				}

				if d.EquipmentId() != 0 {
					err := equipment.Delete(d.EquipmentId())
					if err != nil {
						r.l.Printf("Generating equipment item %d for character %d, they were not awarded this item. Check request in ESO service.")
						return
					}
				}
				Producer(r.l, context.Background()).Emit(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
			}
		}
	}
}

func (r *DropExpiration) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
