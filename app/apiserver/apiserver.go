package apiserver

import (
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	//logger *logrus.
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		//logger: logrus.New(),
	}
}

func (a *APIServer) Start() error {
	return nil
}