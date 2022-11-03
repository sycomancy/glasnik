package main

import "flag"

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "Listen address the server is running")
	flag.Parse()
	svc := NewLoggingService(&adsFetcher{})

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()
}
