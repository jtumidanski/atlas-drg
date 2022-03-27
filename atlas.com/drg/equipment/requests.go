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

func CreateRandom(itemId uint32) requests.PostRequest[Attributes] {
	input := &DataContainer{
		Data: DataBody{
			Id:   "0",
			Type: "com.atlas.eso.attribute.EquipmentAttributes",
			Attributes: Attributes{
				ItemId: itemId,
			},
		}}
	return requests.MakePostRequest[Attributes](fmt.Sprintf(randomEquipmentResource), input)
}

func Delete(l logrus.FieldLogger, span opentracing.Span) func(equipmentId uint32) error {
	return func(equipmentId uint32) error {
		return requests.Delete(l, span)(fmt.Sprintf(equipResource, equipmentId))
	}
}
