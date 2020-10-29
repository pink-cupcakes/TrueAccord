package main

import (
	// "http"
	// "fmt"

	"true_accord/shared/trueaccordapiconnector"

	log "github.com/sirupsen/logrus"
)

var trueAccordAPIConnector trueaccordapiconnector.TrueAccordAPIConnector

func Initialize() {
	log.SetFormatter(&log.TextFormatter{})
	trueAccordAPIConnector = trueaccordapiconnector.NewTrueAccordAPIConnector()
}

func main() {
	Initialize()

	debts, err := trueAccordAPIConnector.GetDebts()
	if err != nil {
		log.WithFields(log.Fields{
			"Message": err.ErrorMessage,
			"ClientError": err.ClientErrorMessage,
			"InternalErrorMessage": err.InternalErrorMessage,
		}).Error()
	}

	for _, debt := range debts {
		log.Println(debt)
	}
}
