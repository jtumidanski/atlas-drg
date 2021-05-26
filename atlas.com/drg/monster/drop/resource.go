package drop

import (
	"atlas-drg/drop"
	"atlas-drg/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func HandleGetDropById(fl logrus.FieldLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		l := fl.WithFields(logrus.Fields{"originator": "GetDropById", "type": "rest_handler"})
		dropId := getDropId(l)(r)
		d, err := GetDropById(dropId)
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			json.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		rw.WriteHeader(http.StatusOK)
		result := drop.DropDataContainer{Data:
		drop.DropData{
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
		json.ToJSON(result, rw)
	}
}

func getDropId(l logrus.FieldLogger) func(r *http.Request) uint32 {
	return func(r *http.Request) uint32 {
		vars := mux.Vars(r)
		value, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			l.WithError(err).Errorf("Error parsing id as integer")
			return 0
		}
		return uint32(value)
	}
}
