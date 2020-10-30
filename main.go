package main

import (
	// "http"
	// "fmt"

	"true_accord/shared/trueaccordapiconnector"

	log "github.com/sirupsen/logrus"
)

var trueAccordAPIConnector trueaccordapiconnector.TrueAccordAPIConnector

func initialize() {
	log.SetFormatter(&log.TextFormatter{})
	trueAccordAPIConnector = trueaccordapiconnector.NewTrueAccordAPIConnector()
}

func main() {
	initialize()

	debts, err := trueAccordAPIConnector.GetDebts()
	if err != nil {
		err.LogError()
	}

	for _, debt := range debts {
		paymentPlans, err := trueAccordAPIConnector.GetPaymentPlan(debt.ID)
		if err != nil {
			err.LogError()

			// TODO: Depending on the severity, either recover or break
		}

		if(len(paymentPlans) > 1) {
			/** Warning: the TrueAccord API should enforce business logic 1:1 debt to paymentPlan.
				This just adds monitoring if we come across failures in the business logic.

				This logic is not blocking and will default to the first payment plan.
			*/
			log.WithFields(log.Fields{
				"Message": fmt.Sprintf("More than 1 payment plan found for debtID %d", debt.ID),
			}).Info()
		} else if (len(paymentPlans) == 0) {
			// TODO: Handle no payment plans
			continue
		}

		paymentPlan := paymentPlans[0]
		payments, err := trueAccordAPIConnector.GetPayments(paymentPlan.ID)
		if err != nil {
			err.LogError()
		}

		log.Println("Payment plan is:")
		log.Println(paymentPlan)
		for _, payment := range payments {
			log.Println(payment)
		}
	}
}
