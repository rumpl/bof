package main

import (
	"net/http/httptest"
	"testing"

	"github.com/rumpl/bof/integration-cli/daemon"
)

type DockerCLINetworkSuite struct {
	ds *DockerSuite
}

func (s *DockerCLINetworkSuite) TearDownTest(c *testing.T) {
	s.ds.TearDownTest(c)
}

func (s *DockerCLINetworkSuite) OnTimeout(c *testing.T) {
	s.ds.OnTimeout(c)
}

type DockerNetworkSuite struct {
	server *httptest.Server
	ds     *DockerSuite
	d      *daemon.Daemon
}
