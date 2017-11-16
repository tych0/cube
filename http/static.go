package http

import (
	"path"
	"net/http"

	"github.com/anuvu/cube/config"
	"github.com/anuvu/cube/service"
)

type StaticServer interface {}

func NewStaticFileServer(ctx service.Context) StaticServer {
	s := staticServer{}
	ctx.AddLifecycle(&service.Lifecycle{
		ConfigHook: s.ConfigHook,
		StartHook: s.StartHook,
	})
	return s
}

type staticServer struct {
	url string
	root string
}

func (s *staticServer) ConfigHook(store config.Store) error {
	if err := store.Get("staticURL", &s.url); err != nil {
		return err
	}

	if err := store.Get("staticRoot", &s.root); err != nil {
		return err
	}

	return nil
}

func (s *staticServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	theFile := path.Join(s.root, path.Clean(req.URL.Path[len(s.url):]))
	http.ServeFile(w, req, theFile)
}

func (s *staticServer) StartHook(serv Service) {
	serv.Register(s.url, s)
}
