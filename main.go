package main

import (
	"fmt"
	"go-load-balancer/loadbalancer"
	"log"
	"net/http"
)

func main() {
	backendServers := []string{
		"http://localhost:9001",
		"http://localhost:9002",
		"http://localhost:9003",
	}
	lb := loadbalancer.NewLoadBalancer(backendServers)
	http.HandleFunc("/", lb.ServeHTTP)
	fmt.Println("Load Balancer is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
