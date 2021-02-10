package world

import (
	drop2 "atlas-drg/drop"
	"atlas-drg/json"
	"atlas-drg/monster/drop"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func GetDropsInMap(l *log.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		wid := getWorldId(r)
		cid := getChannelId(r)
		mid := getMapId(r)

		ds, _ := drop.Processor(l).GetDropsForMap(wid, cid, mid)
		rw.WriteHeader(http.StatusOK)
		result := drop2.DropDataListContainer{}
		result.Data = make([]drop2.DropData, 0)

		for _, d := range ds {
			data := drop2.DropData{
				Id:   strconv.Itoa(int(d.Id())),
				Type: "com.atlas.drg.rest.attribute.DropAttributes",
				Attributes: drop2.DropAttributes{
					WorldId:         d.WorldId(),
					ChannelId:       d.ChannelId(),
					MapId:           d.MapId(),
					ItemId:          d.ItemId(),
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
		json.ToJSON(result, rw)
	}
}

func getWorldId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.ParseUint(vars["worldId"], 10, 8)
	if err != nil {
		log.Println("Error parsing worldId as byte")
		return 0
	}
	return byte(value)
}

func getChannelId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.ParseUint(vars["channelId"], 10, 8)
	if err != nil {
		log.Println("Error parsing channelId as byte")
		return 0
	}
	return byte(value)
}

func getMapId(r *http.Request) uint32 {
	vars := mux.Vars(r)
	value, err := strconv.ParseUint(vars["mapId"], 10, 32)
	if err != nil {
		log.Println("Error parsing mapId as uint32")
		return 0
	}
	return uint32(value)
}
