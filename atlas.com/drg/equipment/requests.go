package equipment

import (
	"atlas-drg/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	equipmentServicePrefix  string = "/ms/eso/"
	equipmentService               = requests.BaseRequest + equipmentServicePrefix
	equipmentResource              = equipmentService + "equipment"
	randomEquipmentResource        = equipmentService + "equipment?random=true"
	equipResource                  = equipmentResource + "/%d"
)

func CreateRandom(l logrus.FieldLogger, span opentracing.Span) func(itemId uint32) (*DataContainer, error) {
	return func(itemId uint32) (*DataContainer, error) {
		input := &DataContainer{
			Data: DataBody{
				Id:   "0",
				Type: "com.atlas.eso.attribute.EquipmentAttributes",
				Attributes: Attributes{
					ItemId: itemId,
				},
			}}
		ro := &DataContainer{}
		err := requests.Post(l, span)(fmt.Sprintf(randomEquipmentResource), input, ro, &requests.ErrorListDataContainer{})
		if err != nil {
			return nil, err
		}
		return ro, nil
	}
}

func Delete(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) error {
	return func(equipmentId uint32) error {
		return requests.Delete(l, span)(fmt.Sprintf(equipResource, equipmentId))
	}
}
