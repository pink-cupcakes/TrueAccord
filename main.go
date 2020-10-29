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
		paymentPlans, err := trueAccordAPIConnector.GetPaymentPlan(debt.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"Message": err.ErrorMessage,
				"ClientError": err.ClientErrorMessage,
				"InternalErrorMessage": err.InternalErrorMessage,
			}).Error()

			// Depending on the severity, either recover or break
		}

		if(len(paymentPlans) > 1) {
			log.WithFields(log.Fields{
				"Message": "More than 1 payment plan found per debt",
			}).Info()
		} else if (len(paymentPlans) == 0) {
			// Handle no payment plans
			continue
		}

		paymentPlan := paymentPlans[0]
		payments, err := trueAccordAPIConnector.GetPayments(paymentPlan.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"Message": err.ErrorMessage,
				"ClientError": err.ClientErrorMessage,
				"InternalErrorMessage": err.InternalErrorMessage,
			}).Error()
		}

		log.Println("Payment plan is:")
		log.Println(paymentPlan)
		for _, payment := range payments {
			log.Println(payment)
		}
	}
}
