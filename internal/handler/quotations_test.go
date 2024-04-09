package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dkshi/bwgtest"
	"github.com/dkshi/bwgtest/internal/service"
	mock_service "github.com/dkshi/bwgtest/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandle_getQuotation(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockQuotations, quotation *bwgtest.Quotation)

	testTable := []struct {
		name               string
		inputUpdateID      string
		mockBehaviour      mockBehaviour
		expectedQuotation  *bwgtest.Quotation
		expectedResponse   quotationResponse
		expectedStatusCode int
	}{
		{
			name:          "OK",
			inputUpdateID: "1",
			mockBehaviour: func(s *mock_service.MockQuotations, quotation *bwgtest.Quotation) {
				s.EXPECT().GetQuotation(quotation.UpdateID).Return(quotation, nil)
			},
			expectedQuotation: &bwgtest.Quotation{
				UpdateID:   1,
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       0.9,
				Success:    true,
				UpdateTime: time.Time{},
			},
			expectedResponse: quotationResponse{
				Code:       "USD/EUR",
				Rate:       0.9,
				UpdateTime: time.Time{},
			},
			expectedStatusCode: 200,
		},
		{
			name:          "Blank ID",
			inputUpdateID: "",
			mockBehaviour: func(s *mock_service.MockQuotations, quotation *bwgtest.Quotation) {
			},
			expectedQuotation:  &bwgtest.Quotation{},
			expectedResponse:   quotationResponse{},
			expectedStatusCode: 400,
		},
		{
			name:          "Incorrect ID",
			inputUpdateID: "test",
			mockBehaviour: func(s *mock_service.MockQuotations, quotation *bwgtest.Quotation) {
			},
			expectedQuotation:  &bwgtest.Quotation{},
			expectedResponse:   quotationResponse{},
			expectedStatusCode: 400,
		},
		{
			name:          "Service error",
			inputUpdateID: "1",
			mockBehaviour: func(s *mock_service.MockQuotations, quotation *bwgtest.Quotation) {
				s.EXPECT().GetQuotation(1).Return(quotation, fmt.Errorf("service error"))
			},
			expectedQuotation: &bwgtest.Quotation{},
			expectedResponse: quotationResponse{},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			quotationService := mock_service.NewMockQuotations(c)

			testCase.mockBehaviour(quotationService, testCase.expectedQuotation)

			s := &service.Service{Quotations: quotationService}
			handler := NewHandler(s)

			r := gin.New()
			r.GET("/quotations/get", handler.getQuotation)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/quotations/get?id="+testCase.inputUpdateID, bytes.NewBufferString(""))

			r.ServeHTTP(w, req)

			var q quotationResponse

			if err := json.NewDecoder(w.Body).Decode(&q); err != nil {
				t.Fatalf("error decoding json from response; %s", err.Error())
			}

			assert.Equal(t, testCase.expectedResponse, q)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandle_updateQuotation(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockQuotations, codeFrom, codeTo string, updateID int, err error)

	testTable := []struct {
		name               string
		inputCodeFrom      string
		inputCodeTo        string
		mockBehaviour      mockBehaviour
		expectedUpdateID   int
		expectedStatusCode int
	}{
		{
			name:          "OK",
			inputCodeFrom: "USD",
			inputCodeTo:   "EUR",
			mockBehaviour: func(s *mock_service.MockQuotations, codeFrom, codeTo string, updateID int, err error) {
				s.EXPECT().InsertQuotation(codeFrom, codeTo).Return(updateID, err)
			},
			expectedUpdateID:   1,
			expectedStatusCode: 200,
		},
		{
			name:               "Missing parameters",
			inputCodeFrom:      "",
			inputCodeTo:        "",
			mockBehaviour:      func(s *mock_service.MockQuotations, codeFrom, codeTo string, updateID int, err error) {},
			expectedUpdateID:   0,
			expectedStatusCode: 400,
		},
		{
			name:               "Service error",
			inputCodeFrom:      "USD",
			inputCodeTo:        "EUR",
			mockBehaviour:      func(s *mock_service.MockQuotations, codeFrom, codeTo string, updateID int, err error) {
				s.EXPECT().InsertQuotation(codeFrom, codeTo).Return(0, fmt.Errorf("service error"))
			},
			expectedUpdateID:   0,
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			quotationService := mock_service.NewMockQuotations(c)

			handler := &Handler{service: &service.Service{Quotations: quotationService}}

			testCase.mockBehaviour(quotationService, testCase.inputCodeFrom, testCase.inputCodeTo, testCase.expectedUpdateID, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/update?from=%s&to=%s", testCase.inputCodeFrom, testCase.inputCodeTo), nil)

			router := gin.Default()
			router.GET("/update", handler.updateQuotation)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandle_getLatestQuotation(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockQuotations, codeFrom, codeTo string, quotation *bwgtest.Quotation, err error)

	testTable := []struct {
		name               string
		inputCodeFrom      string
		inputCodeTo        string
		mockBehaviour      mockBehaviour
		expectedQuotation  *bwgtest.Quotation
		expectedStatusCode int
	}{
		{
			name:          "OK",
			inputCodeFrom: "USD",
			inputCodeTo:   "EUR",
			mockBehaviour: func(s *mock_service.MockQuotations, codeFrom, codeTo string, quotation *bwgtest.Quotation, err error) {
				s.EXPECT().GetLatestQuotation(codeFrom, codeTo).Return(quotation, err)
			},
			expectedQuotation: &bwgtest.Quotation{
				UpdateID:   1,
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       0.9,
				Success:    true,
				UpdateTime: time.Time{},
			},
			expectedStatusCode: 200,
		},
		{
			name:          "Missing parameters",
			inputCodeFrom: "",
			inputCodeTo:   "",
			mockBehaviour: func(s *mock_service.MockQuotations, codeFrom, codeTo string, quotation *bwgtest.Quotation, err error) {},
			expectedQuotation: &bwgtest.Quotation{
				UpdateID:   0,
				CodeFrom:   "",
				CodeTo:     "",
				Rate:       0,
				Success:    false,
				UpdateTime: time.Time{},
			},
			expectedStatusCode: 400,
		},
		{
			name:          "Service error",
			inputCodeFrom: "USD",
			inputCodeTo:   "EUR",
			mockBehaviour: func(s *mock_service.MockQuotations, codeFrom, codeTo string, quotation *bwgtest.Quotation, err error) {
				s.EXPECT().GetLatestQuotation(codeFrom, codeTo).Return(&bwgtest.Quotation{}, fmt.Errorf("service error"))
			},
			expectedQuotation: &bwgtest.Quotation{
				UpdateID:   0,
				CodeFrom:   "",
				CodeTo:     "",
				Rate:       0,
				Success:    false,
				UpdateTime: time.Time{},
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			quotationService := mock_service.NewMockQuotations(c)

			handler := &Handler{service: &service.Service{Quotations: quotationService}}

			testCase.mockBehaviour(quotationService, testCase.inputCodeFrom, testCase.inputCodeTo, testCase.expectedQuotation, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/latest?from=%s&to=%s", testCase.inputCodeFrom, testCase.inputCodeTo), nil)

			router := gin.Default()
			router.GET("/latest", handler.getLatestQuotation)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}