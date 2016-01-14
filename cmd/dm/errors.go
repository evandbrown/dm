package main

import (
	log "github.com/Sirupsen/logrus"
)

func Check(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("An error occurred")
	}
}
