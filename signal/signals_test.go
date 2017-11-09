package cube

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/anuvu/cube/service"
	. "github.com/smartystreets/goconvey/convey"
)

type sigH struct {
	s []os.Signal
}

func (s *sigH) handle(sig os.Signal) {
	s.s = append(s.s, sig)
}

func TestSignals(t *testing.T) {
	Convey("Create a signal Router", t, func() {
		s := NewSignalRouter()
		So(s, ShouldNotBeNil)
		Convey("Should be able add handler", func() {
			So(s.IsIgnored(syscall.SIGINT), ShouldBeFalse)
			So(s.IsHandled(syscall.SIGINT), ShouldBeFalse)
			s.Handle(syscall.SIGINT, func(os.Signal) {})
			So(s.IsHandled(syscall.SIGINT), ShouldBeTrue)
		})
		Convey("Should be able to reset handler", func() {
			s.Reset(syscall.SIGINT)
			So(s.IsHandled(syscall.SIGINT), ShouldBeFalse)
		})
		Convey("Should be able to ignore signal", func() {
			s.Ignore(syscall.SIGINT)
			So(s.IsIgnored(syscall.SIGINT), ShouldBeTrue)
			So(s.IsHandled(syscall.SIGINT), ShouldBeFalse)
		})

		Convey("Should be able to start the service", func() {
			sh := &sigH{[]os.Signal{}}
			s.Handle(syscall.SIGINT, sh.handle)
			So(s.IsHandled(syscall.SIGINT), ShouldBeTrue)
			So(len(sh.s), ShouldEqual, 0)
			// Check lifecycle
			svc := s.(service.LifeCycle)
			So(svc, ShouldNotBeNil)
			So(svc.OnConfigure(nil), ShouldBeNil)
			So(svc.IsHealthy(), ShouldBeFalse)
			svc.OnStart()
			So(svc.IsHealthy(), ShouldBeTrue)

			Convey("Should be able to handle a signal", func() {
				// Fire a signal
				p, err := os.FindProcess(os.Getpid())
				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				p.Signal(syscall.SIGINT)

				// Sleep for a second
				time.Sleep(time.Second)
				So(len(sh.s), ShouldEqual, 1)
				So(sh.s[0], ShouldEqual, syscall.SIGINT)
			})

			Convey("Should be able to stop the service", func() {
				So(svc.IsHealthy(), ShouldBeTrue)
				So(svc.OnStop(), ShouldBeNil)
				time.Sleep(time.Second)
				So(svc.IsHealthy(), ShouldBeFalse)
			})
		})
	})
}
