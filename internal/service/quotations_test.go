package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_GetQuotationRate(t *testing.T) {
	quotationsService := &QuotationsService{}

	testCases := []struct {
		name          string
		codeFrom      string
		codeTo        string
		expectedError error
	}{
		{"OK", "USD", "EUR", nil},
		{"InvalidCodeTo", "USD", "SOS", fmt.Errorf("error from api server: not found")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			url := fmt.Sprintf("%s?from=%s&to=%s", exchangeHost, testCase.codeFrom, testCase.codeTo)
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("error getting response: %s", err.Error())
			}

			var conversion Conversion
			err = json.NewDecoder(resp.Body).Decode(&conversion)
			if err != nil {
				t.Fatalf("error decoding json: %s", err.Error())
			}

			correctRate := conversion.Rates[testCase.codeTo]

			rate, err := quotationsService.GetQuotationRate(testCase.codeFrom, testCase.codeTo)

			assert.Equal(t, correctRate, rate)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
