package drop

import (
	registries "atlas-drg/configuration"
	"atlas-drg/equipment"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"time"
)

const ExpirationTaskName = "drop_expiration_task"

type ExpirationTask struct {
	l        logrus.FieldLogger
	interval uint64
}

func NewExpirationTask(l logrus.FieldLogger, interval uint64) *ExpirationTask {
	return &ExpirationTask{l, interval}
}

func (r *ExpirationTask) Run() {
	var expire uint64

	c, err := registries.GetConfiguration()
	if err != nil {
		expire = 180000
	} else {
		expire = c.ItemExpireInterval
	}

	span := opentracing.StartSpan(ExpirationTaskName)
	ds := GetRegistry().GetAllDrops()
	for _, d := range ds {
		if d.Status() == "AVAILABLE" {
			if d.DropTime()+expire < uint64(time.Now().UnixNano()/int64(time.Millisecond)) {
				_, err := GetRegistry().RemoveDrop(d.Id())
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
				emitExpiredEvent(r.l, span)(d.WorldId(), d.ChannelId(), d.MapId(), d.Id())
			}
		}
	}
	span.Finish()
}

func (r *ExpirationTask) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
