package main

import (
	"slices"
	"testing"
)

func TestGetNextServer_RoundRobin(t *testing.T) {
	backend_servers := []string{"server_A", "server_B", "server_C"}
	strategy := &RoundRobinStrategy{}
	lb := NewLoadBalancer(backend_servers, strategy)

	server, err := lb.GetNextServer()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if server != "server_A" {
		t.Errorf("Expected server_A on first call, got %s", server)
	}

	lb.GetNextServer() // server_B
	lb.GetNextServer() // server_C
	server, _ = lb.GetNextServer()

	if server != "server_A" {
		t.Errorf("Expected server A on fourth call again, got %s", server)
	}
}

func TestGetNextServer_Random(t *testing.T) {
	backend_servers := []string{"server_A", "server_B", "server_C"}
	strategy := &RandomStrategy{}
	lb := NewLoadBalancer(backend_servers, strategy)

	server, err := lb.GetNextServer()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !slices.Contains(backend_servers, server) {
		t.Error("server not in the list of backend servers provided")
	}

}

func TestGetNextServer_NoServers(t *testing.T){
	backend_servers := []string{}
	strategy := &RandomStrategy{}
	lb := NewLoadBalancer(backend_servers, strategy)

	_,err := lb.GetNextServer()
	if err == nil{
		t.Error("Expected error, instead got nil")
	}
	if err != nil && err.Error() != "no servers" {
		t.Errorf("Expected error 'no servers' , got %v", err)
	}
}
