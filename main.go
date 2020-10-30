package main

import (
	"encoding/json"
	"errors"
	"fmt"

	// "http"

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
	RemainingDebt string `json:"remaining_amount,omitempty"`
	NextBillingDate string `json:"next_payment_due_date,omitempty"`
}

func initialize() {
	log.SetFormatter(&log.TextFormatter{})
	trueAccordAPIConnector = trueaccordapiconnector.NewTrueAccordAPIConnector()
}

func findNextPaymentDate(paymentPlan *trueaccordapiconnector.PaymentPlan) (nextPaymentDate time.Time, subtotalOwed float64, err error) {
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

func getPaymentHistory(payments []trueaccordapiconnector.Payment) (totalPayments float64) {
	if(len(payments) == 0) {
		return
	}

	for _, payment := range payments {
		totalPayments += payment.Amount
	}

	return
}

// debtDataEnrichment assumes there is a valid payment plan
func debtDataEnrichment(debt trueaccordapiconnector.Debt, nextPaymentDate time.Time, subtotalOwed, totalPayments float64) (res EnrichedDebt) {
	if(nextPaymentDate.IsZero()) {
		return
	}

	res = EnrichedDebt{
		debt,
		true,
		fmt.Sprintf("%.2f", math.Max(subtotalOwed - totalPayments, 0)),
		nextPaymentDate.Format(time.RFC3339),
	}

	return
}

func logResult(res EnrichedDebt) error {
	out, err  := json.Marshal(res)
	if(err != nil) {
		return err
	}

	fmt.Println(string(out))
	return nil
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
			continue
		}

		if(paymentPlan == nil) {
			res := EnrichedDebt{debt, false, "", ""}

			logError := logResult(res)
			if(logError != nil) {
				err = httphelpers.NewAPIError(logError, fmt.Sprintf("Failed to log result for debtID: %d", debt.ID))
				err.LogError()
			}
			
			continue
		}

		payments, err := trueAccordAPIConnector.GetPayments(paymentPlan.ID)
		if err != nil {
			err.LogError()
		}

		nextPaymentDate, subtotalOwed, findPaymentErr := findNextPaymentDate(paymentPlan)
		if(findPaymentErr != nil) {
			err = httphelpers.NewAPIError(findPaymentErr, fmt.Sprintf("Failed to process payment plan for debtID: %d", debt.ID))
			err.LogError()
		}

		totalPaid := getPaymentHistory(payments)

		res := debtDataEnrichment(debt, nextPaymentDate, subtotalOwed, totalPaid)
		logError := logResult(res)
		if(logError != nil) {
			err = httphelpers.NewAPIError(logError, fmt.Sprintf("Failed to log result for debtID: %d", debt.ID))
			err.LogError()
		}
	}
}
