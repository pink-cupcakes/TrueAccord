package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"
	"time"
	trueaccordapiconnector "true_accord/shared/trueaccordapi"

	"github.com/stretchr/testify/assert"
)

// Bod ... returns the beginning of the day for a given timestamp
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.UTC().Location())
}

func TestAggregateNextPaymentInfoSuccess(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-05")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoSuccessNoAmountOwedBalance(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 0, "WEEKLY", 51.25, "2020-09-28"}

	_, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No payment plan amount to pay"), err)
}

func TestAggregateNextPaymentInfoSuccessPartialPayments(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-05")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 51.25)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoSuccessFullyPaidToDate(t *testing.T) {
	now := time.Now()
	nowString := now.Format("2006-01-02")
	testSuccessNextPaymentDate := Bod(now.AddDate(0, 0, 7))

	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, nowString}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 51.25)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoFailureInvalidInstallAmount(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 0, "2020-09-28"}

	_, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No installment_amount found"), err)

	testPaymentPlan = trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", -2, "2020-09-28"}

	_, err = aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No installment_amount found"), err)
}

func TestAggregateNextPaymentInfoSuccessFutureStartDate(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, "2021-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2021-09-28")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoSuccessBiweekly(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "BI_WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-12")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoSuccessPartialInstallments(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 110.00, "WEEKLY", 25.00, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-26")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
}

func TestAggregateNextPaymentInfoSuccessPaymentsInProgress(t *testing.T) {
	now := time.Now()
	startDate := Bod(now.AddDate(0, 0, -9))
	stringStartDate := startDate.Format("2006-01-02")

	successNextPaymentDate := Bod(now.AddDate(0, 0, 5))

	installmentAmount := 51.25
	amountOwed := installmentAmount * 4

	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, amountOwed, "WEEKLY", installmentAmount, stringStartDate}

	nextPaymentDate, err := aggregateNextPaymentInfo(&testPaymentPlan, 0.00)

	assert.Nil(t, err)
	assert.Equal(t, successNextPaymentDate, nextPaymentDate)
}

func TestGetPaymentHistorySuccess(t *testing.T) {
	testPaymentAmounts := []float64{51.25, 51.25}
	testPayments := []trueaccordapiconnector.Payment{{testPaymentAmounts[0], "2020-09-29", 0}, {testPaymentAmounts[1], "2020-10-29", 0}}

	testSumResult := testPaymentAmounts[0] + testPaymentAmounts[1]
	totalPayments := aggregatePayments(testPayments)
	assert.Equal(t, testSumResult, totalPayments)
}

func TestGetPaymentHistorySuccessFutureDates(t *testing.T) {
	now := time.Now()
	firstDate := Bod(now.AddDate(0, 0, -18))
	stringFirstDate := firstDate.Format("2006-01-02")

	secondDate := Bod(now.AddDate(0, 0, 15))
	stringSecondDate := secondDate.Format("2006-01-02")

	testPaymentAmounts := []float64{51.25, 51.25}

	testPayments := []trueaccordapiconnector.Payment{{testPaymentAmounts[0], stringFirstDate, 0}, {testPaymentAmounts[1], stringSecondDate, 0}}

	testSumResult := testPaymentAmounts[0]
	totalPayments := aggregatePayments(testPayments)
	assert.Equal(t, testSumResult, totalPayments)
}

func TestDebtDataEnrichmentSuccess(t *testing.T) {
	now := time.Now()
	nextPaymentDate := Bod(now.AddDate(0, 0, -18))
	stringNextPaymentDate := nextPaymentDate.Format(time.RFC3339)
	paymentPlan := trueaccordapiconnector.PaymentPlan{1, 1, 102.5, "WEEKLY", 10, "2020-10-10"}
	totalPayments := 51.25
	testDebt := trueaccordapiconnector.Debt{0, 102.5}

	successEnrichedDebt := EnrichedDebt{testDebt, true, "51.25", stringNextPaymentDate}

	enrichedDebt := debtDataEnrichment(testDebt, nextPaymentDate, &paymentPlan, totalPayments)
	assert.Equal(t, successEnrichedDebt, enrichedDebt)
}

func TestLogResultSuccess(t *testing.T) {
	var loggedError bytes.Buffer
	log.SetOutput(&loggedError)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	now := time.Now()
	nextPaymentDate := Bod(now.AddDate(0, 0, -18))
	stringNextPaymentDate := nextPaymentDate.Format(time.RFC3339)
	testDebt := trueaccordapiconnector.Debt{0, 102.5}
	testEnrichedDebt := EnrichedDebt{testDebt, true, "51.25", stringNextPaymentDate}

	err := logResult(testEnrichedDebt)
	assert.Nil(t, err, "LogResults should succeed with valid EnrichedDebt")
	assert.NotNil(t, loggedError.String(), "LogResults should log the input to console")
}
