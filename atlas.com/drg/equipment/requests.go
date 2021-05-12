package equipment

import (
	"atlas-drg/rest/requests"
	"fmt"
)

const (
	equipmentServicePrefix string = "/ms/eso/"
	equipmentService              = requests.BaseRequest + equipmentServicePrefix
	equipmentResource             = equipmentService + "equipment"
	equipResource                 = equipmentResource + "/%d"
)

func Create(itemId uint32) (*DataContainer, error) {
	input := &DataContainer{
		Data: DataBody{
			Id:   "0",
			Type: "com.atlas.eso.attribute.EquipmentAttributes",
			Attributes: Attributes{
				ItemId: itemId,
			},
		}}
	resp, err := requests.Post(fmt.Sprintf(equipmentResource), input)
	if err != nil {
		return nil, err
	}

	ro := &DataContainer{}
	err = requests.ProcessResponse(resp, ro)
	if err != nil {
		return nil, err
	}
	return ro, nil
}

func Delete(equipmentId uint32) error {
	_, err := requests.Delete(fmt.Sprintf(equipResource, equipmentId))
	return err
}
