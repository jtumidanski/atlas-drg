package drop

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
	getDropById = "get_drop_by_id"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	sRouter := router.PathPrefix("/drops/{id}").Subrouter()
	sRouter.HandleFunc("", registerGetDropById(l))
}

type idHandler func(id uint32) http.HandlerFunc

func parseId(l logrus.FieldLogger, next idHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			l.WithError(err).Errorf("Error parsing id as integer")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(value))(w, r)
	}
}

func registerGetDropById(l logrus.FieldLogger) http.HandlerFunc {
	fl := l.WithFields(logrus.Fields{"originator": "GetDropById", "type": "rest_handler"})
	return rest.RetrieveSpan(getDropById, func(span opentracing.Span) http.HandlerFunc {
		return parseId(fl, func(id uint32) http.HandlerFunc {
			return handleGetDropById(fl)(span)(id)
		})
	})
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func handleGetDropById(fl logrus.FieldLogger) func(span opentracing.Span) func(dropId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(dropId uint32) http.HandlerFunc {
		return func(dropId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				d, err := GetDropById(dropId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					json.ToJSON(&GenericError{Message: err.Error()}, w)
					return
				}

				w.WriteHeader(http.StatusOK)
				result := drop.DropDataContainer{Data: drop.DropData{
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
				}}
				json.ToJSON(result, w)
			}
		}
	}
}
