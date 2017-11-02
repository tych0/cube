package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type httpConfig struct {
	Port int `json:"port"`
}

func (h *httpConfig) UnmarshalJSON(b []byte) error {
	v := map[string]int{}
	if e := json.Unmarshal(b, &v); e != nil {
		return e
	}
	if p, ok := v["port"]; ok {
		h.Port = p
	} else {
		return fmt.Errorf("port must be present")
	}
	return nil
}

type loggerConfig struct {
	File string `json:"file"`
}

func TestJSONStore(t *testing.T) {
	Convey("On a json store", t, func() {
		r := NewRegistry()
		So(r.Register("http", httpConfig{}), ShouldBeNil)

		goodJSON := strings.NewReader(`{"http": {"port": 8080}}
			{"logger": {"file": "/var/log/test.log"}}`)
		s := NewJSONStore(goodJSON, r)
		So(s, ShouldNotBeNil)
		So(s.Registry(), ShouldEqual, r)
		defer s.Close()
		Convey("Should be able to load the file", func() {
			So(s.Open(), ShouldBeNil)
			Convey("should be able load http config", func() {
				v, err := s.Get("http")
				So(err, ShouldBeNil)
				So(v, ShouldNotBeNil)
				cfg := v.(*httpConfig)
				So(cfg, ShouldNotBeNil)
				So(cfg.Port, ShouldEqual, 8080)
			})
			Convey("should not find randon config", func() {
				v, err := s.Get("some_random_key")
				So(err, ShouldBeError)
				So(v, ShouldBeNil)
			})
			Convey("should not be able to find unregistered logger type", func() {
				v, err := s.Get("logger")
				So(err, ShouldBeError)
				So(v, ShouldBeNil)
			})
			Convey("should not be able to find wrong logger type", func() {
				So(r.Register("logger", &httpConfig{}), ShouldBeNil)
				v, err := s.Get("logger")
				So(err, ShouldBeError)
				So(v, ShouldBeNil)
			})
			Convey("should be able to find logger after registering the type", func() {
				So(r.Register("logger", &loggerConfig{}), ShouldBeNil)
				v, err := s.Get("logger")
				So(err, ShouldBeNil)
				So(v, ShouldNotBeNil)
				cfg := v.(*loggerConfig)
				So(cfg, ShouldNotBeNil)
				So(cfg.File, ShouldEqual, "/var/log/test.log")
			})
		})
	})
}

func TestBadJSON(t *testing.T) {
	Convey("Load bad json data", t, func() {
		r := NewRegistry()
		So(r.Register("http", httpConfig{}), ShouldBeNil)
		Convey("should be a json parse error", func() {
			badJSON := strings.NewReader(`{"http": {"portx": "8080"}`)
			s := NewJSONStore(badJSON, r)
			So(s.Open(), ShouldBeError)
		})

		Convey("should error out on bad http config", func() {
			badKeyJSON := strings.NewReader(`{"http": {"portx": "8080"}}`)
			s := NewJSONStore(badKeyJSON, r)
			So(s.Open(), ShouldNotBeNil)
		})
	})
}
