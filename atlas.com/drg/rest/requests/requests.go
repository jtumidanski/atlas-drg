package requests

import (
	json2 "atlas-drg/json"
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	BaseRequest string = "http://atlas-nginx:80"
)

func Get(url string, resp interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}

	err = ProcessResponse(r, resp)
	return err
}

func Post(url string, input interface{}) (*http.Response, error) {
	jsonReq, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	r, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Delete(url string) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	return client.Do(r)
}

func ProcessResponse(r *http.Response, rb interface{}) error {
	err := json2.FromJSON(rb, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func processErrorResponse(r *http.Response, eb interface{}) error {
	if r.ContentLength > 0 {
		err := json2.FromJSON(eb, r.Body)
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}

// objectMap is a simple representation of a json map. key is a string, and value is a nested object
type objectMap map[string]interface{}

// ObjectMapper maps a objectMap to a concrete data type
type ObjectMapper func(objectMap) interface{}

// conditionalMapperProvider returns a string representing the type of object the ObjectMapper handles, as well as the
// ObjectMapper itself.
type conditionalMapperProvider func() (string, ObjectMapper)

// concreteObjectProvider returns a new (empty) concrete object
type concreteObjectProvider func() interface{}

// DataSegment is a slice of concrete objects
type DataSegment []interface{}

// unmarshalRoot will take a raw byte array, and using mapper functions, produce a data DataSegment and includes
// DataSegment representing a jsonapi.org requests response
func UnmarshalRoot(data []byte, root ObjectMapper, options ...conditionalMapperProvider) (DataSegment, DataSegment, error) {
	var dataResult DataSegment
	var includeResult DataSegment

	var single = struct {
		Data     objectMap
		Included []objectMap
	}{}
	err := json.Unmarshal(data, &single)
	if err == nil {
		dataResult = produceSingle(single.Data, root)
		includeResult = produceList(single.Included, includeMapper(options...))
		return dataResult, includeResult, nil
	} else {
		var list = struct {
			Data     []objectMap
			Included []objectMap
		}{}
		err = json.Unmarshal(data, &list)
		if err == nil {
			dataResult = produceList(list.Data, root)
			includeResult = produceList(list.Included, includeMapper(options...))
			return dataResult, includeResult, nil
		}
		return nil, nil, err
	}
}

// includeMapper represents a ObjectMapper for handling the jsonapi.org includes data section
func includeMapper(options ...conditionalMapperProvider) ObjectMapper {
	return func(o objectMap) interface{} {
		return addInclude(o, options...)
	}
}

// produceSingle will produce a DataSegment given a single objectMap from a ObjectMapper
func produceSingle(o objectMap, m ObjectMapper) DataSegment {
	return append(make([]interface{}, 0), m(o))
}

// produceList will produce a DataSegment given a objectMap slice from a ObjectMapper
func produceList(o []objectMap, m ObjectMapper) DataSegment {
	if len(o) > 0 {
		var result = make([]interface{}, 0)
		for _, x := range o {
			result = append(result, m(x))
		}
		return result
	}
	return nil
}

// transformMap will populate the concrete object provided by the concreteObjectProvider from the objectMap and return it.
func transformMap(dpf concreteObjectProvider, x objectMap) interface{} {
	b, err := json.Marshal(x)
	if err == nil {
		r := dpf()
		err = json.Unmarshal(b, r)
		if err == nil {
			return r
		}
	}
	return nil
}

func MapperFunc(dpf concreteObjectProvider) ObjectMapper {
	return func(x objectMap) interface{} {
		return transformMap(dpf, x)
	}
}

func UnmarshalData(tf string, dpf concreteObjectProvider) (string, ObjectMapper) {
	return tf, MapperFunc(dpf)
}

// addInclude processes a map structure representing the jsonapi.org data object to produce a concrete struct.
// given a set of functions which could produce a concrete struct
func addInclude(x objectMap, options ...conditionalMapperProvider) interface{} {
	t := x["type"].(string)
	for _, o := range options {
		tf, mf := o()
		if t == tf {
			return mf(x)
		}
	}
	return nil
}
