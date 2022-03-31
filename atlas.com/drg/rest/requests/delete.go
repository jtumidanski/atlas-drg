package requests

import (
	"bytes"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Delete(l logrus.FieldLogger, span opentracing.Span) func(url string, input interface{}) error {
	return func(url string, input interface{}) error {
		var req *http.Request
		var err error
		if input != nil {
			var jsonReq []byte
			jsonReq, err = json.Marshal(input)
			if err != nil {
				return err
			}
			req, err = http.NewRequest(http.MethodDelete, url, bytes.NewReader(jsonReq))
		} else {
			req, err = http.NewRequest(http.MethodDelete, url, nil)
		}

		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		err = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			l.WithError(err).Errorf("Unable to decorate request headers with OpenTracing information.")
		}
		r, err := http.DefaultClient.Do(req)

		if input != nil {
			l.WithFields(logrus.Fields{"method": http.MethodDelete, "status": r.Status, "path": url, "input": input, "response": ""}).Debugf("Printing request.")
		} else {
			l.WithFields(logrus.Fields{"method": http.MethodDelete, "status": r.Status, "path": url, "response": ""}).Debugf("Printing request.")
		}

		return err
	}
}
