package config

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type sampleConfig struct {
	Foo string
}

func TestRegistry(t *testing.T) {
	Convey("On a new registry", t, func() {
		r := NewRegistry()
		So(r, ShouldNotBeNil)
		Convey("I should not be able to register a nil type", func() {
			So(r.Register("test", nil), ShouldBeError)
		})
		Convey("I should not be able to register a primitive type", func() {
			So(r.Register("test", 10), ShouldBeError)
		})
		Convey("I should be able to register a struct", func() {
			So(r.Register("test", sampleConfig{}), ShouldBeNil)
			Convey("and retrieve the type", func() {
				So(r.Get("test"), ShouldEqual, reflect.TypeOf(sampleConfig{}))
			})
		})
		Convey("I should be able to register a struct pointer", func() {
			So(r.Register("test", &sampleConfig{}), ShouldBeNil)
			Convey("and retrieve the type", func() {
				So(r.Get("test"), ShouldEqual, reflect.TypeOf(sampleConfig{}))
			})
		})
	})
}
