package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService(t *testing.T) {
	Convey("On a service object", t, func() {
		s := NewBaseService("test")
		So(s, ShouldNotBeNil)
		Convey("Name must be test", func() {
			So(s.Name(), ShouldEqual, "test")
		})
		Convey("And ID must not be nil", func() {
			So(s.ID, ShouldNotBeNil)
			So(len(s.ID()), ShouldEqual, 16)
		})
	})
}
