package main

import (
	"override-gui/backend"
)

type BackendService struct {
	manager *backend.Manager
}

func (g *BackendService) StartServer() backend.ResponseData {
	g.manager = backend.NewServerManager()
	return g.manager.Start()
}
func (g *BackendService) StopServer() backend.ResponseData {
	return g.manager.Stop()
}
func (g *BackendService) ReadConfig() backend.ResponseData {
	return backend.ReadConfig()
}
func (g *BackendService) UpdateConfig(config string) backend.ResponseData {
	return backend.UpdateConfig(config)
}
