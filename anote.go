package main

import (
	"github.com/anonutopia/gowaves"
)

func initAnote() *gowaves.WavesNodeClient {
	anc := &gowaves.WavesNodeClient{
		Host: AnoteNodeURL,
		Port: 443,
	}

	return anc
}
