package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type httpConfig struct {
	Port int `json:"port"`
}

func TestJSONStore(t *testing.T) {
	Convey("On a json store", t, func() {
		r := NewRegistry()
		So(r.Register("http", httpConfig{}), ShouldBeNil)
		Convey("should error for bad file", func() {
			s := NewJSONStore("some_random_file", r)
			So(s.Load(), ShouldBeError)
		})

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
			s := NewJSONStore("bad_cfg.json", r)
			So(s.Load(), ShouldBeError)
		})
	})
}
