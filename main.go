package main

import (
	"errors"
	// "http"
	// "fmt"
	"math"
	"time"

	"true_accord/shared/httphelpers"
	"true_accord/shared/trueaccordapiconnector"

	log "github.com/sirupsen/logrus"
)

var trueAccordAPIConnector trueaccordapiconnector.TrueAccordAPIConnector

const (
	weeklyInterval = time.Hour * 24 * 7
	biweeklyInterval = time.Hour * 24 * 7 * 2
)

type EnrichedDebt struct {
	trueaccordapiconnector.Debt

	HasPaymentPlan bool `json:"is_in_payment_plan"`
	RemainingDebt float64 `json:"remaining_amount,omitempty"`
	NextBillingDate string `json:"next_payment_due_date,omitempty"`
}

func initialize() {
	log.SetFormatter(&log.TextFormatter{})
	trueAccordAPIConnector = trueaccordapiconnector.NewTrueAccordAPIConnector()
}

func findNextPayment(paymentPlan *trueaccordapiconnector.PaymentPlan) (nextPaymentDate time.Time, subtotalOwed float64, err error) {
	if(paymentPlan == nil) {
		err = errors.New("No payment plan provided")
		return
	}

	paymentDate, err := time.Parse("2006-01-02", paymentPlan.StartDate)
	if err != nil {
		return
	}

	var paymentInterval time.Duration
	if(paymentPlan.InstallmentFrequency == "WEEKLY") {
		paymentInterval = weeklyInterval
	} else if(paymentPlan.InstallmentFrequency == "BI_WEEKLY") {
		paymentInterval = biweeklyInterval
	} else {
		err = errors.New("Unhandled payment interval")
		return
	}

	now := time.Now()
	for(paymentDate.Before(now) && paymentPlan.AmountToPay > subtotalOwed) {
		paymentDate = paymentDate.Add(paymentInterval)
		subtotalOwed += paymentPlan.InstallmentAmount
	}

	return paymentDate, math.Min(subtotalOwed, paymentPlan.AmountToPay), nil
}

func debtDataEnrichment(debt trueaccordapiconnector.Debt, paymentPlan *trueaccordapiconnector.PaymentPlan, payments []*trueaccordapiconnector.Payment) (res EnrichedDebt, err *httphelpers.APIError) {
	res.Debt = debt

	if(paymentPlan == nil) {
		return
	}

	res.HasPaymentPlan = true
	return
}

func main() {
	initialize()

	debts, err := trueAccordAPIConnector.GetDebts()
	if err != nil {
		err.LogError()
	}

	for _, debt := range debts {
		paymentPlan, err := trueAccordAPIConnector.GetPaymentPlan(debt.ID)
		if err != nil {
			err.LogError()

			// TODO: Depending on the severity, either recover or break
		}

		if(paymentPlan == nil) {
			log.Println("NO PAYMENT PLAN FOUND")
			continue
		}

		paymentDate, total, paymentErr := findNextPayment(paymentPlan)
		if(paymentErr != nil) {
			log.Println(paymentErr)
		}
		println(paymentDate.Format(time.RFC3339))
		println(total)

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
