package http

import (
	"fmt"
	"testing"
	"net/http"
	"io/ioutil"

	"github.com/anuvu/cube/service"
	"github.com/anuvu/cube/config"
	. "github.com/smartystreets/goconvey/convey"
)

type staticFakeConfig struct {
	setConfig bool
}

func (fc staticFakeConfig) Open() error {
	return nil
}

func (fc staticFakeConfig) Close() {
}

func (fc staticFakeConfig) Get(name string, config interface{}) error {
	if name == "staticURL" {
		url := config.(*string)
		*url = "/foo"
	} else if name == "staticRoot" {
		root := config.(*string)
		*root = "/dev"
	} else {
		return fmt.Errorf("unknown key %s", name)
	}

	return nil
}

func TestStaticFileServer(t *testing.T) {
	Convey("http server actually serves stuff", t, func() {
		grp := service.NewGroup("base", nil)
		So(grp.AddService(func() config.Store { return staticFakeConfig{} }), ShouldBeNil)
		So(grp.AddService(NewService), ShouldBeNil)
		So(grp.AddService(NewStaticFileServer), ShouldBeNil)
		So(grp.Configure(), ShouldBeNil)
		So(grp.Start(), ShouldBeNil)

		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/foo/null", port))
		So(err, ShouldBeNil)
		bytes, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(len(bytes), ShouldEqual, 0)
		So(grp.Stop(), ShouldBeNil)
	})
}
