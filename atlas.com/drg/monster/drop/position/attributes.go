package position

type DropPositionInputDataContainer struct {
	Data DropPositionData `json:"data"`
}

type DropPositionData struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes DropPositionAttributes `json:"attributes"`
}

type DropPositionAttributes struct {
	InitialX  int16 `json:"initialX"`
	InitialY  int16 `json:"initialY"`
	FallbackX int16 `json:"fallbackX"`
	FallbackY int16 `json:"fallbackY"`
}
