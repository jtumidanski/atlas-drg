package point

import (
	"atlas-drg/monster/drop/position"
	"atlas-drg/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapResource                        = mapInformationService + "maps/%d"
	dropPosition                       = mapResource + "/dropPosition"
)

var MapInformationRequests = func() *mapInformationRequests {
	return &mapInformationRequests{}
}

type mapInformationRequests struct {
}

func (a *mapInformationRequests) CalculateDropPosition(mapId uint32, initialX int16, initialY int16, fallbackX int16, fallbackY int16) (*MapPointDataContainer, error) {
	input := &position.DropPositionInputDataContainer{Data: position.DropPositionData{
		Id:   "0",
		Type: "com.atlas.mis.attribute.DropPositionInputAttributes",
		Attributes: position.DropPositionAttributes{
			InitialX:  initialX,
			InitialY:  initialY,
			FallbackX: fallbackX,
			FallbackY: fallbackY,
		},
	}}
	resp, err := requests.Post(fmt.Sprintf(dropPosition, mapId), input)
	if err != nil {
		return nil, err
	}

	ar := &MapPointDataContainer{}
	err = requests.ProcessResponse(resp, ar)
	if err != nil {
		return nil, err
	}

	return ar, nil
}
