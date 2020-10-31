package trueaccordapi

import "true_accord/shared/httphelpers"

type testTrueAccordAPIConnector struct{}

// NewTestTrueAccordAPIConnector ... returns an interface of TrueAccordAPIConnector
func NewTestTrueAccordAPIConnector() TrueAccordAPIConnector {
	return &testTrueAccordAPIConnector{}
}

// GetDebts ... returns all the debts from TrueAccord API
func (tta *testTrueAccordAPIConnector) GetDebts() (debts []Debt, err *httphelpers.APIError) {
	debts = []Debt{{0, 123.46}, {1, 100}, {2, 4920.34}, {3, 12938.0}, {4, 9238.02}}

	return
}
