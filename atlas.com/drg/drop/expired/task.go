package expired

import (
	registries "atlas-drg/configuration"
	"atlas-drg/drop"
	"atlas-drg/equipment"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"time"
)

const DropExpirationTask = "drop_expiration_task"

type DropExpiration struct {
	l        logrus.FieldLogger
	interval uint64
}

func NewDropExpiration(l logrus.FieldLogger, interval uint64) *DropExpiration {
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

	span := opentracing.StartSpan(DropExpirationTask)
	ds := drop.GetRegistry().GetAllDrops()
	for _, d := range ds {
		if d.Status() == "AVAILABLE" {
			if d.DropTime()+expire < uint64(time.Now().UnixNano()/int64(time.Millisecond)) {
				_, err := drop.GetRegistry().RemoveDrop(d.Id())
				if err != nil {
					r.l.WithError(err).Errorf("Unable to remove drop from registry.")
					continue
				}

				if d.EquipmentId() != 0 {
					err := equipment.Delete(r.l, span)(d.EquipmentId())
					if err != nil {
						r.l.WithError(err).Errorf("Deleting equipment item %d.", d.EquipmentId())
						return
					}
				}
				DropExpired(r.l, span)(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
			}
		}
	}
	span.Finish()
}

func (r *DropExpiration) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
