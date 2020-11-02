package main

import (
	"encoding/json"
	"errors"
	"fmt"

	// "http"

	"math"
	"time"

	"true_accord/shared/httphelpers"
	trueaccordapiconnector "true_accord/shared/trueaccordapi"

	log "github.com/sirupsen/logrus"
)

var trueAccordAPIConnector trueaccordapiconnector.TrueAccordAPIConnector

const (
	weeklyInterval   = time.Hour * 24 * 7
	biweeklyInterval = time.Hour * 24 * 7 * 2
)

type EnrichedDebt struct {
	trueaccordapiconnector.Debt

	HasPaymentPlan  bool   `json:"is_in_payment_plan"`
	RemainingDebt   string `json:"remaining_amount"`
	NextBillingDate string `json:"next_payment_due_date"`
}

func initialize() {
	log.SetFormatter(&log.TextFormatter{})
	trueAccordAPIConnector = trueaccordapiconnector.NewTrueAccordAPIConnector()
}

// aggregateNextPaymentInfo ... returns the next payment date and amount owed according to payment plan (not debt)
func aggregateNextPaymentInfo(paymentPlan *trueaccordapiconnector.PaymentPlan, totalPaid float64) (nextPaymentDate time.Time, err error) {
	if paymentPlan == nil {
		err = errors.New("No payment plan provided")
		return
	}

	paymentDate, err := time.Parse("2006-01-02", paymentPlan.StartDate)
	if err != nil {
		return
	}

	if paymentPlan.AmountToPay <= float64(0) {
		err = errors.New("No payment plan amount to pay")
		return
	}

	if paymentPlan.InstallmentAmount <= float64(0) && paymentPlan.AmountToPay != float64(0) {
		err = errors.New("No installment_amount found")
		return
	}

	subtotalOwed := paymentPlan.InstallmentAmount

	var paymentInterval time.Duration
	if paymentPlan.InstallmentFrequency == "WEEKLY" {
		paymentInterval = weeklyInterval
	} else if paymentPlan.InstallmentFrequency == "BI_WEEKLY" {
		paymentInterval = biweeklyInterval
	} else {
		err = errors.New("Unhandled payment interval")
		return
	}

	// Retrieve next payment date and amount owed by payment date (independent of actual payments)
	now := time.Now()

	for paymentDate.Before(now) && paymentPlan.AmountToPay > subtotalOwed {
		paymentDate = paymentDate.Add(paymentInterval)
		subtotalOwed += paymentPlan.InstallmentAmount
	}

	return paymentDate, nil
}

// aggregatePayments ... returns the total amount paid
func aggregatePayments(payments []trueaccordapiconnector.Payment) (totalPayments float64) {
	if len(payments) == 0 {
		return
	}

	now := time.Now()

	for _, payment := range payments {
		paymentDate, err := time.Parse("2006-01-02", payment.Date)
		if err != nil {
			return
		}

		if paymentDate.Before(now) {
			totalPayments += payment.Amount
		}
	}

	return
}

// debtDataEnrichment ... returns the debt object with paymentPlan and next payment information
func debtDataEnrichment(debt trueaccordapiconnector.Debt, nextPaymentDate time.Time, paymentPlan *trueaccordapiconnector.PaymentPlan, totalPayments float64) (res EnrichedDebt) {
	if nextPaymentDate.IsZero() {
		return
	}

	remainingAmount := math.Max(paymentPlan.AmountToPay-totalPayments, 0)
	if remainingAmount == 0 {
		res = EnrichedDebt{
			debt,
			true,
			fmt.Sprintf("%.2f", remainingAmount),
			"null",
		}
		return
	}

	res = EnrichedDebt{
		debt,
		true,
		fmt.Sprintf("%.2f", remainingAmount),
		nextPaymentDate.Format(time.RFC3339),
	}

	return
}

func logResult(res EnrichedDebt) error {
	out, err := json.Marshal(res)
	if err != nil {
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

		if paymentPlan == nil {
			res := EnrichedDebt{debt, false, fmt.Sprintf("%.2f", debt.Amount), "null"}

			logError := logResult(res)
			if logError != nil {
				err = httphelpers.NewAPIError(logError, fmt.Sprintf("Failed to log result for debtID: %d", debt.ID))
				err.LogError()
			}

			continue
		}

		payments, err := trueAccordAPIConnector.GetPayments(paymentPlan.ID)
		if err != nil {
			err.LogError()
		}

		totalPaid := aggregatePayments(payments)

		nextPaymentDate, findPaymentErr := aggregateNextPaymentInfo(paymentPlan, totalPaid)
		if findPaymentErr != nil {
			err = httphelpers.NewAPIError(findPaymentErr, fmt.Sprintf("Failed to process payment plan for debtID: %d", debt.ID))
			err.LogError()
		}

		res := debtDataEnrichment(debt, nextPaymentDate, paymentPlan, totalPaid)
		logError := logResult(res)
		if logError != nil {
			err = httphelpers.NewAPIError(logError, fmt.Sprintf("Failed to log result for debtID: %d", debt.ID))
			err.LogError()
		}
	}
}
