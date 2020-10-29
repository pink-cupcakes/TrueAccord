package trueaccordapiconnector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"true_accord/shared/httphelpers"

	// log "github.com/sirupsen/logrus"
)

// TrueAccordAPIConnector ... is an interface of appapi methods called
type TrueAccordAPIConnector interface {
	GetDebts() (debts []*Debt, err *httphelpers.APIError)
	GetPaymentPlan(debtID int64) (paymentPlans []*PaymentPlan, err *httphelpers.APIError)
	GetPayments(paymentPlanID int64) (payments []*Payment, err *httphelpers.APIError)
}

type trueAccordAPIConnector struct{}

var trueAccordAPIURL = os.Getenv("TRUEACCORD_API_URL")

const (
	getDebts        = "debts"
	getPaymentPlans = "payment_plans"
	getPayments        = "payments"
)

// Debt ... is the debt response model returned from TrueAccord API
type Debt struct {
	ID int64 `json:"id"`
	Amount float64 `json:"amount"`
}

// PaymentPlan ... is the payment plan response model returned from TrueAccord API
type PaymentPlan struct {
	ID int64 `json:"id"`
	DebtID int64 `json:"debt_id"`
	AmountToPay float64 `json:"amount_to_pay"`
	Frequency string `json:"installment_frequency"`
	Amount float64 `json:"installment_amount"`
	StartDate string `json:"start_date"`
}

// Payment ... is the customer payment response model returned from TrueAccord API
type Payment struct {
	Amount float64 `json:"amount"`
	Date string `json:"date"`
	PaymentPlanID int64 `json:"payment_plan_id"`
}

// NewTrueAccordAPIConnector ... returns an interface of TrueAccordAPIConnector
func NewTrueAccordAPIConnector() TrueAccordAPIConnector {
	return &trueAccordAPIConnector{}
}

// GetDebts ... returns all the debts from TrueAccord API
func (j *trueAccordAPIConnector) GetDebts() (debts []*Debt, err *httphelpers.APIError) {
	resp, requestErr := j.makeRequest(getDebts, "GET", nil, nil)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to GET debts"

		err = httphelpers.NewAPIError(requestErr, clientErr).SetInternalErrorMessage("Failed to read response body from GET debts")
		return 
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET debts, non-200 response: " + string(b))
		clientErr := "Failed to GET debts"

		err = httphelpers.NewAPIError(requestErr, clientErr)
		return debts, err
	}

	requestErr = json.Unmarshal(b, &debts)
	if requestErr != nil {
		clientErr := "Failed to GET debts"

		err = httphelpers.NewAPIError(requestErr, clientErr)
		err.SetInternalErrorMessage("Failed to unmarshal GET debts result")
		return nil, err
	}

	return
}

// GetPaymentPlan ... returns a payment plan (of any) for a given debt from TrueAccord API
func (j *trueAccordAPIConnector) GetPaymentPlan(debtID int64) (paymentPlans []*PaymentPlan, err *httphelpers.APIError) {
	params := url.Values{"debt_id": []string{strconv.Itoa(int(debtID))}}
	resp, requestErr := j.makeRequest(getPaymentPlans, "GET", nil, params)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to GET payment plans"

		err = httphelpers.NewAPIError(requestErr, clientErr).SetInternalErrorMessage("Failed to read response body from GET payment plans")
		return 
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET payment plans, non-200 response: " + string(b))
		clientErr := "Failed to GET payment plans"

		err = httphelpers.NewAPIError(requestErr, clientErr)
		return paymentPlans, err
	}

	requestErr = json.Unmarshal(b, &paymentPlans)
	if requestErr != nil {
		clientErr := "Failed to GET payment plans"

		err = httphelpers.NewAPIError(requestErr, clientErr)
		err.SetInternalErrorMessage("Failed to unmarshal GET payment plans result")
		return nil, err
	}

	return
}

// GetPayments ... returns the payment activities for a given payment plan from TrueAccord API
func (j *trueAccordAPIConnector) GetPayments(paymentPlanID int64) (payments []*Payment, err *httphelpers.APIError) {
	params := url.Values{"payment_plan_id": []string{strconv.Itoa(int(paymentPlanID))}}
	resp, requestErr := j.makeRequest(getPayments, "GET", nil, params)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to GET payments"

		err = httphelpers.NewAPIError(requestErr, clientErr).SetInternalErrorMessage("Failed to read response body from GET payments")
		return 
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET payments, non-200 response: " + string(b))
		clientErr := "Failed to GET payments"

		err = httphelpers.NewAPIError(requestErr, clientErr)
		return payments, err
	}

	requestErr = json.Unmarshal(b, &payments)
	if requestErr != nil {
		clientErr := "Failed to GET payments"

		err = httphelpers.NewAPIError(requestErr, clientErr).SetInternalErrorMessage("Failed to unmarshal GET payments result")
		return nil, err
	}

	return
}

func (j *trueAccordAPIConnector) makeRequest(endpoint, method string, body []byte, params url.Values) (resp *http.Response, err error) {
	client := &http.Client{}
	URL, err := url.Parse(fmt.Sprintf("%s/%s", trueAccordAPIURL, endpoint))
	if err != nil {
		return nil, err
	}

	if params != nil {
		URL.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, URL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
