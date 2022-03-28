package character

import (
	"atlas-drg/drop"
	"atlas-drg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameSpawnDrop = "spawn_character_drop_command"
	topicTokenSpawnDrop   = "TOPIC_SPAWN_CHARACTER_DROP_COMMAND"
)

func SpawnDropConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[spawnDropCommand](consumerNameSpawnDrop, topicTokenSpawnDrop, groupId, handleSpawnDrop())
}

type spawnDropCommand struct {
	WorldId      byte   `json:"worldId"`
	ChannelId    byte   `json:"channelId"`
	MapId        uint32 `json:"mapId"`
	ItemId       uint32 `json:"itemId"`
	EquipmentId  uint32 `json:"equipmentId"`
	Quantity     uint32 `json:"quantity"`
	Mesos        uint32 `json:"mesos"`
	DropType     byte   `json:"dropType"`
	X            int16  `json:"x"`
	Y            int16  `json:"y"`
	OwnerId      uint32 `json:"ownerId"`
	OwnerPartyId uint32 `json:"ownerPartyId"`
	DropperId    uint32 `json:"dropperId"`
	DropperX     int16  `json:"dropperX"`
	DropperY     int16  `json:"dropperY"`
	PlayerDrop   bool   `json:"playerDrop"`
	Mod          byte   `json:"mod"`
}

func handleSpawnDrop() kafka.HandlerFunc[spawnDropCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command spawnDropCommand) {
		drop.SpawnCharacterDrop(l, span)(command.WorldId, command.ChannelId, command.MapId, command.ItemId,
			command.EquipmentId, command.Quantity, command.Mesos, command.DropType, command.X, command.Y,
			command.OwnerId, command.OwnerPartyId, command.DropperId, command.DropperX, command.DropperY,
			command.PlayerDrop, command.Mod)
	}
}
