package topic

import (
	"atlas-drg/json"
	"atlas-drg/rest/requests"
	"atlas-drg/retry"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	topicsServicePrefix string = "/ms/tds/"
	topicsService              = requests.BaseRequest + topicsServicePrefix
	topicById                  = topicsService + "topics/%s"
)

var TopicRequests = func(l *log.Logger) *topic {
	return &topic{l: l}
}

type topic struct {
	l *log.Logger
}

func (t *topic) GetTopic(topic string) (*TopicData, error) {
	get := func(attempt int) (bool, interface{}, error) {
		r, err := http.Get(fmt.Sprintf(topicById, topic))
		if err != nil {
			t.l.Printf("[WARN] unable to retrieve topic data for %s, will retry.", topic)
			return true, r, err
		}
		return false, r, nil
	}

	r, err := retry.RetryResponse(get, 10)
	if err != nil {
		t.l.Printf("[ERROR] unable to retrieve topic data for %s", topic)
		return nil, err
	}
	if val, ok := r.(*http.Response); ok {
		return t.decodeResponse(topic, err, val)
	}
	return nil, errors.New("unexpected output from retry function")
}

func (t *topic) decodeResponse(topic string, err error, val *http.Response) (*TopicData, error) {
	td := &TopicDataContainer{}
	err = json.FromJSON(td, val.Body)
	if err != nil {
		t.l.Printf("[ERROR] decoding topic data for %s", topic)
		return nil, err
	}
	return &td.Data, nil
}
