package http

import (
	"fmt"
	"testing"
	"io/ioutil"
	"net/http"

	"github.com/anuvu/cube/service"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	port = 8000
	msg = "hello"
)


type fakeConfig struct {}

func (fc fakeConfig) Open() error {
	return nil
}

func (fc fakeConfig) Close() {
}

func (fc fakeConfig) Get(name string, config interface{}) error {
	c := config.(*int)
	*c = port
	return nil
}

type testHandler struct {}

func (th testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}

func TestHTTPServer(t *testing.T) {
	Convey("http server actually serves stuff", t, func() {
		s := NewService(service.NewContext()).(*server)
		s.Register("/foo", testHandler{})
		So(s.ConfigHook(fakeConfig{}), ShouldBeNil)
		So(s.StartHook(), ShouldBeNil)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/foo", port))
		So(err, ShouldBeNil)
		bytes, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldEqual, string(msg))
		So(s.HealthHook(), ShouldBeTrue)
		So(s.StopHook(), ShouldBeNil)
	})
}
