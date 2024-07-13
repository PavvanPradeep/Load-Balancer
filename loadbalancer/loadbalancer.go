package loadbalancer

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	Servers   []*Server
	Mutex     sync.Mutex
	lastIndex int
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, serverURL := range servers {
		lb.Servers = append(lb.Servers, &Server{URL: serverURL})
	}
	return lb
}

func (lb *LoadBalancer) GetServerWithLeastConnections() *Server {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	fmt.Println("Checking servers for least connections:")
	var leastConnServer *Server
	minConns := -1
	candidates := []*Server{}

	for _, server := range lb.Servers {
		fmt.Printf("Server %s has %d active connections\n", server.URL, server.ActiveConns)
		if minConns == -1 || server.ActiveConns < minConns {
			minConns = server.ActiveConns
			candidates = []*Server{server}
		} else if server.ActiveConns == minConns {
			candidates = append(candidates, server)
		}
	}

	if len(candidates) == 1 {
		leastConnServer = candidates[0]
	} else {
		// in case least connections algo fails round robin will be used as a failover mechanism
		leastConnServer = candidates[lb.lastIndex%len(candidates)]
		lb.lastIndex++
	}

	fmt.Printf("Selected server %s with %d active connections\n", leastConnServer.URL, leastConnServer.ActiveConns)
	return leastConnServer
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetServerWithLeastConnections()
	if server == nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	server.IncrementConnections()
	defer server.DecrementConnections()

	proxyURL := server.URL + r.URL.Path
	fmt.Printf("Forwarding request to %s\n", proxyURL)
	resp, err := http.Get(proxyURL)
	if err != nil {
		fmt.Printf("Error forwarding request to %s: %v\n", proxyURL, err)
		http.Error(w, "Failed to reach backend server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("Error copying response body: %v\n", err)
	}
}
