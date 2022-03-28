package world

import (
	"atlas-drg/drop"
	"atlas-drg/json"
	"atlas-drg/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	getDropsInMap = "get_drops_in_map"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	iRouter := router.PathPrefix("/worlds/{worldId}/channels/{channelId}/maps/{mapId}/drops").Subrouter()
	iRouter.HandleFunc("", registerGetDropsInMap(l))
}

type mapHandler func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc

func parseMap(l logrus.FieldLogger, next mapHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		worldId, err := strconv.ParseUint(vars["worldId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		channelId, err := strconv.ParseUint(vars["channelId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mapId, err := strconv.ParseUint(vars["mapId"], 10, 32)
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(byte(worldId), byte(channelId), uint32(mapId))(w, r)
	}
}

func registerGetDropsInMap(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getDropsInMap, func(span opentracing.Span) http.HandlerFunc {
		fl := l.WithFields(logrus.Fields{"originator": "GetDropById", "type": "rest_handler"})
		return parseMap(fl, func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return handleGetDropsInMap(fl)(span)(worldId, channelId, mapId)
		})
	})
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func handleGetDropsInMap(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
		return func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				ds, _ := drop.GetForMap(worldId, channelId, mapId)

				result := drop.DropDataListContainer{}
				result.Data = make([]drop.DropData, 0)

				for _, d := range ds {
					data := drop.DropData{
						Id:   strconv.Itoa(int(d.Id())),
						Type: "com.atlas.drg.rest.attribute.DropAttributes",
						Attributes: drop.DropAttributes{
							WorldId:         d.WorldId(),
							ChannelId:       d.ChannelId(),
							MapId:           d.MapId(),
							ItemId:          d.ItemId(),
							EquipmentId:     d.EquipmentId(),
							Quantity:        d.Quantity(),
							Meso:            d.Meso(),
							DropType:        d.Type(),
							DropX:           d.X(),
							DropY:           d.Y(),
							OwnerId:         d.OwnerId(),
							OwnerPartyId:    d.OwnerPartyId(),
							DropTime:        d.DropTime(),
							DropperUniqueId: d.DropperId(),
							DropperX:        d.DropperX(),
							DropperY:        d.DropperY(),
							CharacterDrop:   d.CharacterDrop(),
							Mod:             d.Mod(),
						},
					}
					result.Data = append(result.Data, data)
				}

				w.WriteHeader(http.StatusOK)
				err := json.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Uanble to encode result.")
				}
			}
		}
	}
}
