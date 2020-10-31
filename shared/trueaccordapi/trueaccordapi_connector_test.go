package trueaccordapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var trueAccordTestAPIConnector TrueAccordAPIConnector

func TestMain(m *testing.M) {
	trueAccordTestAPIConnector = NewTrueAccordAPIConnector()
	os.Exit(m.Run())
}

func TestGetDebtsSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testDebtsResponse := `[
		{
			"amount": 123.46,
			"id": 0
		},
		{
			"amount": 100,
			"id": 1
		},
		{
			"amount": 4920.34,
			"id": 2
		},
		{
			"amount": 12938,
			"id": 3
		},
		{
			"amount": 9238.02,
			"id": 4
		}
	]`

	var testDebts []Debt

	unmarshalTestResponseErr := json.Unmarshal([]byte(testDebtsResponse), &testDebts)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getDebts),
		httpmock.NewStringResponder(200, testDebtsResponse))

	res, err := trueAccordTestAPIConnector.GetDebts()
	assert.Nil(t, err, "GetDebts success")
	assert.Equal(t, len(testDebts), len(res))
}

func TestGetDebtsEmptyResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testDebtsResponse := `[]`

	var testDebts []Debt

	unmarshalTestResponseErr := json.Unmarshal([]byte(testDebtsResponse), &testDebts)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getDebts),
		httpmock.NewStringResponder(200, testDebtsResponse))

	res, err := trueAccordTestAPIConnector.GetDebts()
	assert.Nil(t, err, "GetDebts empty response")
	assert.Equal(t, testDebts, res)
}

func TestGetDebtsIncorrectResponseFormat(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testDebtsResponse := `[
		{
			"amount": "String_instead_of_float_field",
			"id": 0
		},
		{
			"amount": 100,
			"id": "String_instead_of_int"
		},
		{
			"amount": 4920.34,
			"id": ["Testing nested structures", 1]
		},
		{
			"amount": 12938,
			"id": 3
		},
		{
			"amount": 9238.02,
			"id": 4
		}
	]`

	// Exact URL match
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getDebts),
		httpmock.NewStringResponder(200, testDebtsResponse))

	_, err := trueAccordTestAPIConnector.GetDebts()
	assert.NotNil(t, err, "GetDebts should return error with incorrect response struture")
	assert.Equal(t, "Failed to unmarshal GET debts result", err.InternalErrorMessage)
}

func TestGetDebtsNon200Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getDebts),
		httpmock.NewStringResponder(503, ""))

	_, err := trueAccordTestAPIConnector.GetDebts()
	assert.NotNil(t, err, "GetDebts should return error with non-200 response")
}

func TestGetPaymentPlansSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	debtID := int64(0)

	testPaymentPlansResponse := `[
		{
			"amount_to_pay": 102.5,
			"debt_id": 0,
			"id": 0,
			"installment_amount": 51.25,
			"installment_frequency": "WEEKLY",
			"start_date": "2020-09-28"
		}
	]`

	var testPaymentPlans []*PaymentPlan

	unmarshalTestResponseErr := json.Unmarshal([]byte(testPaymentPlansResponse), &testPaymentPlans)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	expectedQuery := url.Values{
		"debt_id": []string{fmt.Sprintf("%d", debtID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPaymentPlans), expectedQuery,
		httpmock.NewStringResponder(200, testPaymentPlansResponse))

	res, err := trueAccordTestAPIConnector.GetPaymentPlan(debtID)
	assert.Nil(t, err, "GetPaymentPlanSuccess")
	assert.Equal(t, testPaymentPlans[0], res)
}

func TestGetPaymentPlansMoreThanOneToDebt(t *testing.T) {
	var loggedError bytes.Buffer
	log.SetOutput(&loggedError)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	debtID := int64(0)

	testPaymentPlansResponse := `[
		{
			"amount_to_pay": 102.5,
			"debt_id": 0,
			"id": 0,
			"installment_amount": 51.25,
			"installment_frequency": "WEEKLY",
			"start_date": "2020-09-28"
		},
		{
			"amount_to_pay": 102.5,
			"debt_id": 0,
			"id": 0,
			"installment_amount": 51.25,
			"installment_frequency": "WEEKLY",
			"start_date": "2020-10-28"
		}
	]`

	var testPaymentPlans []*PaymentPlan

	unmarshalTestResponseErr := json.Unmarshal([]byte(testPaymentPlansResponse), &testPaymentPlans)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	expectedQuery := url.Values{
		"debt_id": []string{fmt.Sprintf("%d", debtID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPaymentPlans), expectedQuery,
		httpmock.NewStringResponder(200, testPaymentPlansResponse))

	res, err := trueAccordTestAPIConnector.GetPaymentPlan(debtID)
	assert.Nil(t, err, "GetPaymentPlanSuccess more than one payment plan found")
	assert.Equal(t, testPaymentPlans[0], res)

	/** Logrus is logging additional information in the format:
	"time=\"2020-10-30T14:31:46-07:00\" level=info Message=\"More than 1 payment plan found for debtID 0\"\n"

	TODO: Fix error matching for logrus.
	*/
	assert.NotNil(t, loggedError.String(), "GetPaymentPlanSuccess should log an informational error when more than one payment plan is found")
}

func TestGetPaymentPlansFailureIncorrectPaymentFormat(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	debtID := int64(0)

	testPaymentPlansResponse := `[
		{
			"amount_to_pay": "String_instead_of_float",
			"debt_id": 0,
			"id": 0,
			"installment_amount": 51.25,
			"installment_frequency": 0,
			"start_date": "z"
		}
	]`

	expectedQuery := url.Values{
		"debt_id": []string{fmt.Sprintf("%d", debtID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPaymentPlans), expectedQuery,
		httpmock.NewStringResponder(200, testPaymentPlansResponse))

	_, err := trueAccordTestAPIConnector.GetPaymentPlan(debtID)
	assert.NotNil(t, err, "GetPaymentPlan should return error with incorrect response struture")
}

func TestGetPaymentPlansNon200Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	debtID := int64(0)

	expectedQuery := url.Values{
		"debt_id": []string{fmt.Sprintf("%d", debtID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPaymentPlans), expectedQuery,
		httpmock.NewStringResponder(503, ""))

	_, err := trueAccordTestAPIConnector.GetPaymentPlan(debtID)
	assert.NotNil(t, err, "GetPaymentPlans should return error with non-200 response")
}

func TestGetPaymentsSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	paymentPlanID := int64(0)

	testPaymentsResponse := `[
		{
			"amount": 51.25,
			"date": "2020-09-29",
			"payment_plan_id": 0
		},
		{
			"amount": 51.25,
			"date": "2020-10-29",
			"payment_plan_id": 0
		}
	]`

	var testPayments []Payment

	unmarshalTestResponseErr := json.Unmarshal([]byte(testPaymentsResponse), &testPayments)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	expectedQuery := url.Values{
		"payment_plan_id": []string{fmt.Sprintf("%d", paymentPlanID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPayments), expectedQuery,
		httpmock.NewStringResponder(200, testPaymentsResponse))

	res, err := trueAccordTestAPIConnector.GetPayments(paymentPlanID)
	assert.Nil(t, err, "GetPaymentsSuccess")
	assert.Equal(t, testPayments, res)
}

func TestGetPaymentsFailureIncorrectFormat(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	paymentPlanID := int64(0)

	testPaymentsResponse := `[
		{
			"amount": "String_instead_of_float",
			"date": "2020-09-29",
			"payment_plan_id": 0
		},
		{
			"amount": 51.25,
			"date": "2020-10-29",
			"payment_plan_id": 0
		}
	]`

	expectedQuery := url.Values{
		"payment_plan_id": []string{fmt.Sprintf("%d", paymentPlanID)},
	}

	// Exact URL match
	httpmock.RegisterResponderWithQuery("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPayments), expectedQuery,
		httpmock.NewStringResponder(200, testPaymentsResponse))

	_, err := trueAccordTestAPIConnector.GetPayments(paymentPlanID)
	assert.NotNil(t, err, "GetPayments should return error with incorrect response struture")
}

func TestGetPaymentsNon200Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	paymentPlanID := int64(0)

	// Exact URL match
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/%s", trueAccordAPIURL, getPayments),
		httpmock.NewStringResponder(503, ""))

	_, err := trueAccordTestAPIConnector.GetPaymentPlan(paymentPlanID)
	assert.NotNil(t, err, "GetPaymentPlans should return error with non-200 response")
}
