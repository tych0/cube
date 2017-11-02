package config

import (
	"encoding/json"
	"fmt"
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

func TestJSONStore(t *testing.T) {
	Convey("On a json store", t, func() {
		r := NewRegistry()
		So(r.Register("http", httpConfig{}), ShouldBeNil)
		Convey("should error for bad file", func() {
			s := NewJSONStore("some_random_file", r)
			So(s.Load(), ShouldBeError)
		})

		// Contents of cfg_test: {"http": {"port": 8080}}
		s := NewJSONStore("cfg_test.json", r)
		So(s, ShouldNotBeNil)
		So(s.Registry(), ShouldEqual, r)
		Convey("Should be able to load the file", func() {
			So(s.Load(), ShouldBeNil)
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
		})
	})
}

func TestBadJSON(t *testing.T) {
	Convey("Load bad json data", t, func() {
		r := NewRegistry()
		So(r.Register("http", httpConfig{}), ShouldBeNil)
		Convey("should be a json parse error", func() {
			// Contents of bad_cfg.json: {"http": {"portx": "8080"}
			s := NewJSONStore("bad_cfg.json", r)
			So(s.Load(), ShouldBeError)
			Convey("should error out on bad http config", func() {
				e := s.(*jsonStore).processData([]byte(`{"http": {"portx": "8080"}}`))
				So(e, ShouldNotBeNil)
				fmt.Println(s.Get("http"))
			})
		})
	})
}
