// Copyright 2015 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package exporter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/mtail/metrics"
	"github.com/kylelemons/godebug/pretty"
)

var handleJSONTests = []struct {
	name     string
	metrics  []metrics.Metric
	expected string
}{
	{"empty",
		[]metrics.Metric{},
		"null",
	},
	{"single",
		[]metrics.Metric{
			metrics.Metric{
				Name:        "foo",
				Program:     "test",
				Kind:        metrics.Counter,
				LabelValues: []*metrics.LabelValue{&metrics.LabelValue{[]string{}, &metrics.Datum{Value: 1}}},
			},
		},
		`[
  {
    "Name": "foo",
    "Program": "test",
    "Kind": 0,
    "LabelValues": [
      {
        "Value": {
          "Value": 1,
          "Time": "0001-01-01T00:00:00Z"
        }
      }
    ]
  }
]`,
	},
	{"dimensioned",
		[]metrics.Metric{
			metrics.Metric{
				Name:        "foo",
				Program:     "test",
				Kind:        metrics.Counter,
				Keys:        []string{"a", "b"},
				LabelValues: []*metrics.LabelValue{&metrics.LabelValue{[]string{"1", "2"}, &metrics.Datum{Value: 1}}},
			},
		},
		`[
  {
    "Name": "foo",
    "Program": "test",
    "Kind": 0,
    "Keys": [
      "a",
      "b"
    ],
    "LabelValues": [
      {
        "Labels": [
          "1",
          "2"
        ],
        "Value": {
          "Value": 1,
          "Time": "0001-01-01T00:00:00Z"
        }
      }
    ]
  }
]`,
	},
}

func TestHandleJSON(t *testing.T) {
	for _, tc := range handleJSONTests {
		ms := metrics.Store{}
		for _, metric := range tc.metrics {
			ms.Add(&metric)
		}
		e := New(&ms)
		response := httptest.NewRecorder()
		e.HandleJSON(response, &http.Request{})
		if response.Code != 200 {
			t.Errorf("test case %s: response code not 200: %s", tc.name, response.Code)
		}
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf("test case %s: failed to read response: %s", tc.name, err)
		}
		diff := pretty.Compare(string(b), tc.expected)
		if len(diff) > 0 {
			t.Errorf("test case %s: response not expected: %s", tc.name, diff)
		}
	}
}