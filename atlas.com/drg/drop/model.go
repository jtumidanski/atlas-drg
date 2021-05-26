package drop

type Drop struct {
	id           uint32
	worldId      byte
	channelId    byte
	mapId        uint32
	itemId       uint32
	equipmentId  uint32
	quantity     uint32
	meso         uint32
	dropType     byte
	x            int16
	y            int16
	ownerId      uint32
	ownerPartyId uint32
	dropTime     uint64
	dropperId    uint32
	dropperX     int16
	dropperY     int16
	playerDrop   bool
	mod          byte
	status       string
}

func (d Drop) Id() uint32 {
	return d.id
}

func (d Drop) ItemId() uint32 {
	return d.itemId
}

func (d Drop) Quantity() uint32 {
	return d.quantity
}

func (d Drop) Meso() uint32 {
	return d.meso
}

func (d Drop) Type() byte {
	return d.dropType
}

func (d Drop) X() int16 {
	return d.x
}

func (d Drop) Y() int16 {
	return d.y
}

func (d Drop) OwnerId() uint32 {
	return d.ownerId
}

func (d Drop) OwnerPartyId() uint32 {
	return d.ownerPartyId
}

func (d Drop) DropTime() uint64 {
	return d.dropTime
}

func (d Drop) DropperId() uint32 {
	return d.dropperId
}

func (d Drop) DropperX() int16 {
	return d.dropperX
}

func (d Drop) DropperY() int16 {
	return d.dropperY
}

func (d Drop) PlayerDrop() bool {
	return d.playerDrop
}

func (d Drop) Mod() byte {
	return d.mod
}

func (d Drop) Status() string {
	return d.status
}

func (d Drop) CancelReservation() {
	d.status = "AVAILABLE"
}

func (d Drop) Reserve() {
	d.status = "RESERVED"
}

func (d Drop) MapId() uint32 {
	return d.mapId
}

func (d Drop) WorldId() byte {
	return d.worldId
}

func (d Drop) ChannelId() byte {
	return d.channelId
}

func (d Drop) CharacterDrop() bool {
	return d.playerDrop
}

func (d Drop) EquipmentId() uint32 {
	return d.equipmentId
}
