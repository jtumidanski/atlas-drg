package drop

import (
	"atlas-drg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameSpawnCommand      = "spawn_drop_command"
	consumerNameCancelReservation = "cancel_drop_reservation_command"
	consumerNameReserve           = "reserve_drop_command"
	consumerNamePickup            = "pickup_drop_command"
	topicTokenSpawn               = "TOPIC_SPAWN_DROP_COMMAND"
	topicTokenCancelReservation   = "TOPIC_CANCEL_DROP_RESERVATION_COMMAND"
	topicTokenReserve             = "TOPIC_RESERVE_DROP_COMMAND"
	topicTokenPickup              = "TOPIC_PICKUP_DROP_COMMAND"
)

func SpawnConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[spawnCommand](consumerNameSpawnCommand, topicTokenSpawn, groupId, handleSpawn())
}

type spawnCommand struct {
	WorldId      byte   `json:"worldId"`
	ChannelId    byte   `json:"channelId"`
	MapId        uint32 `json:"mapId"`
	ItemId       uint32 `json:"itemId"`
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

func handleSpawn() kafka.HandlerFunc[spawnCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command spawnCommand) {
		SpawnDrop(l, span)(command.WorldId, command.ChannelId, command.MapId, command.ItemId, command.Quantity,
			command.Mesos, command.DropType, command.X, command.Y, command.OwnerId, command.OwnerPartyId, command.DropperId,
			command.DropperX, command.DropperY, command.PlayerDrop, command.Mod)
	}
}

func CancelReservationConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[cancelReservationCommand](consumerNameCancelReservation, topicTokenCancelReservation, groupId, handleCancelDropReservation())
}

type cancelReservationCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func handleCancelDropReservation() kafka.HandlerFunc[cancelReservationCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command cancelReservationCommand) {
		CancelReservation(l, span)(command.DropId, command.CharacterId)
	}
}

func ReserveConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[reserveCommand](consumerNameReserve, topicTokenReserve, groupId, handleReserveCommand())
}

type reserveCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func handleReserveCommand() kafka.HandlerFunc[reserveCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command reserveCommand) {
		Reserve(l, span)(command.DropId, command.CharacterId)
	}
}

func PickupConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[pickupCommand](consumerNamePickup, topicTokenPickup, groupId, handlePickup())
}

type pickupCommand struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func handlePickup() kafka.HandlerFunc[pickupCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, event pickupCommand) {
		Gather(l, span)(event.DropId, event.CharacterId)
	}
}
