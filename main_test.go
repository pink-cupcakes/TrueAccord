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

func TestFindNextPaymentDateSuccess(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-05")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, 102.5, subTotalOwed)
}

func TestFindNextPaymentDateSuccessNoAmountOwedBalance(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 0, "WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-09-28")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, float64(0), subTotalOwed)
}

func TestFindNextPaymentDateFailureInvalidInstallAmount(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 0, "2020-09-28"}

	_, _, err := findNextPaymentDate(&testPaymentPlan)

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No installment_amount found"), err)

	testPaymentPlan = trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", -2, "2020-09-28"}

	_, _, err = findNextPaymentDate(&testPaymentPlan)

	assert.NotNil(t, err)
	assert.Equal(t, errors.New("No installment_amount found"), err)
}

func TestFindNextPaymentDateSuccessFutureStartDate(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "WEEKLY", 51.25, "2021-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2021-09-28")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, 51.25, subTotalOwed)
}

func TestFindNextPaymentDateSuccessBiweekly(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 102.5, "BI_WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-12")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, 102.5, subTotalOwed)
}

func TestFindNextPaymentDateSuccessPartialInstallments(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, 110.00, "WEEKLY", 25.00, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-10-26")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, 110.00, subTotalOwed)
}

func TestFindNextPaymentDateSuccessPaymentsInProgress(t *testing.T) {
	now := time.Now()
	startDate := Bod(now.AddDate(0, 0, -9))
	stringStartDate := startDate.Format("2006-01-02")

	successNextPaymentDate := Bod(now.AddDate(0, 0, 5))

	installmentAmount := 51.25
	amountOwed := installmentAmount * 4
	successSubtotalOwed := installmentAmount * 3

	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, amountOwed, "WEEKLY", installmentAmount, stringStartDate}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, successNextPaymentDate, nextPaymentDate)
	assert.Equal(t, successSubtotalOwed, subTotalOwed)
}

func TestFindNextPaymentDateSuccessNegativeAmountOwed(t *testing.T) {
	testPaymentPlan := trueaccordapiconnector.PaymentPlan{0, 0, -102.5, "WEEKLY", 51.25, "2020-09-28"}
	testSuccessNextPaymentDate, err := time.Parse("2006-01-02", "2020-09-28")
	if err != nil {
		t.Errorf(err.Error())
	}

	nextPaymentDate, subTotalOwed, err := findNextPaymentDate(&testPaymentPlan)

	assert.Nil(t, err)
	assert.Equal(t, nextPaymentDate, testSuccessNextPaymentDate)
	assert.Equal(t, 0.00, subTotalOwed)
}

func TestGetPaymentHistorySuccess(t *testing.T) {
	testPaymentAmounts := []float64{51.25, 51.25}
	testPayments := []trueaccordapiconnector.Payment{{testPaymentAmounts[0], "2020-09-29", 0}, {testPaymentAmounts[1], "2020-10-29", 0}}

	testSumResult := testPaymentAmounts[0] + testPaymentAmounts[1]
	totalPayments := getPaymentHistory(testPayments)
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
	totalPayments := getPaymentHistory(testPayments)
	assert.Equal(t, testSumResult, totalPayments)
}

func TestDebtDataEnrichmentSuccess(t *testing.T) {
	now := time.Now()
	nextPaymentDate := Bod(now.AddDate(0, 0, -18))
	stringNextPaymentDate := nextPaymentDate.Format(time.RFC3339)
	subTotalOwed := 102.5
	totalPayments := 51.25
	testDebt := trueaccordapiconnector.Debt{0, 102.5}

	successEnrichedDebt := EnrichedDebt{testDebt, true, "51.25", stringNextPaymentDate}

	enrichedDebt := debtDataEnrichment(testDebt, nextPaymentDate, subTotalOwed, totalPayments)
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
