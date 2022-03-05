package main

import (
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 10 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

var client = &http.Client{
	Transport: transport,
	Timeout:   10 * time.Second,
}
