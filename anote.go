package main

import (
	"github.com/anonutopia/gowaves"
)

func initAnote() *gowaves.WavesNodeClient {
	anc := &gowaves.WavesNodeClient{
		Host: "http://localhost",
		Port: 6869,
	}

	return anc
}
