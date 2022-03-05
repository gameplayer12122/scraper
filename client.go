package main

import (
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 15 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 15 * time.Second,
}

var client = &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second,
}
