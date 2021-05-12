package drop

type DropDataContainer struct {
	Data DropData `json:"data"`
}

type DropDataListContainer struct {
	Data []DropData `json:"data"`
}

type DropData struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes DropAttributes `json:"attributes"`
}

type DropAttributes struct {
	WorldId         byte   `json:"worldId"`
	ChannelId       byte   `json:"channelId"`
	MapId           uint32 `json:"mapId"`
	ItemId          uint32 `json:"itemId"`
	EquipmentId     uint32 `json:"equipmentId"`
	Quantity        uint32 `json:"quantity"`
	Meso            uint32 `json:"meso"`
	DropType        byte   `json:"dropType"`
	DropX           int16  `json:"dropX"`
	DropY           int16  `json:"dropY"`
	OwnerId         uint32 `json:"ownerId"`
	OwnerPartyId    uint32 `json:"ownerPartyId"`
	DropTime        uint64 `json:"dropTime"`
	DropperUniqueId uint32 `json:"dropperUniqueId"`
	DropperX        int16  `json:"dropperX"`
	DropperY        int16  `json:"dropperY"`
	CharacterDrop   bool   `json:"playerDrop"`
	Mod             byte   `json:"mod"`
}
