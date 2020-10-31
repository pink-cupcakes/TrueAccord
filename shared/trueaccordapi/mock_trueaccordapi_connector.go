package trueaccordapi

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"net/url"
// 	"strconv"
// 	"true_accord/shared/httphelpers"

// 	log "github.com/sirupsen/logrus"
// )

// type testTrueAccordAPIConnector struct{}

// // NewTestTrueAccordAPIConnector ... returns an interface of TrueAccordAPIConnector
// func NewTestTrueAccordAPIConnector() TrueAccordAPIConnector {
// 	return &testTrueAccordAPIConnector{}
// }

// // GetDebts ... returns all the debts from TrueAccord API
// func (tta *testTrueAccordAPIConnector) GetDebts() (debts []Debt, err *httphelpers.APIError) {
// 	debts = []Debt{{0, 123.46}, {1, 100}, {2, 4920.34}, {3, 12938.0}, {4, 9238.02}}

// 	return
// }

// // GetPaymentPlan ... returns a payment plan (of any) for a given debt from TrueAccord API
// func (ta *trueAccordAPIConnector) GetPaymentPlan(debtID int64) (paymentPlan *PaymentPlan, err *httphelpers.APIError) {
// 	params := url.Values{"debt_id": []string{strconv.Itoa(int(debtID))}}
// 	resp, requestErr := ta.makeRequest(getPaymentPlans, "GET", nil, params)
// 	if err != nil {
// 		return
// 	}

// 	defer resp.Body.Close()

// 	b, requestErr := ioutil.ReadAll(resp.Body)
// 	if requestErr != nil {
// 		clientErr := "Failed to GET payment plans"
// 		err = httphelpers.NewAPIError(requestErr, clientErr).SetInternalErrorMessage("Failed to read response body from GET payment plans")
// 		return
// 	}

// 	if resp.StatusCode != 200 {
// 		requestErr = errors.New("Failed to GET payment plans, non-200 response: " + string(b))
// 		clientErr := "Failed to GET payment plans"
// 		err = httphelpers.NewAPIError(requestErr, clientErr)
// 		return
// 	}

// 	var paymentPlans []PaymentPlan
// 	requestErr = json.Unmarshal(b, &paymentPlans)
// 	if requestErr != nil {
// 		clientErr := "Failed to GET payment plans"
// 		err = httphelpers.NewAPIError(requestErr, clientErr)
// 		err.SetInternalErrorMessage("Failed to unmarshal GET payment plans result")
// 		return
// 	}

// 	if len(paymentPlans) > 1 {
// 		/** TrueAccord API should enforce business logic 1:1 debt to paymentPlan.
// 		This just adds monitoring if we come across failures in the business logic.

// 		This logic is not blocking and will default to the first payment plan.
// 		*/
// 		log.WithFields(log.Fields{
// 			"Message": fmt.Sprintf("More than 1 payment plan found for debtID %d", debtID),
// 		}).Info()
// 	} else if len(paymentPlans) == 0 {
// 		return
// 	}

// 	paymentPlan = &paymentPlans[0]
// 	return
// }
