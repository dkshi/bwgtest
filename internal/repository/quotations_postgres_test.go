package repository

import (
	"testing"
	"time"

	"github.com/dkshi/bwgtest"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestQuotationsPostgres_InsertQuotation(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error creating mock db: %s", err.Error())
	}
	defer db.Close()

	r := NewQuotationsPostgres(db, &pq.Listener{})

	testTable := []struct {
		name           string
		inputQuotation *bwgtest.Quotation
		expectedResult int
		expectedError  error
	}{
		{
			name: "OK",
			inputQuotation: &bwgtest.Quotation{
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       1.2,
				UpdateTime: time.Now(),
				Success:    false,
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mock.ExpectQuery("INSERT INTO quotations").
				WithArgs(testCase.inputQuotation.CodeFrom, testCase.inputQuotation.CodeTo, testCase.inputQuotation.Rate,
					testCase.inputQuotation.UpdateTime, testCase.inputQuotation.Success).
				WillReturnRows(sqlmock.NewRows([]string{"update_id"}).AddRow(1))

			updateID, err := r.InsertQuotation(testCase.inputQuotation)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, updateID)
		})
	}
}

func TestQuotationsPostgres_GetLatestQuotation(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error creating mock db: %s", err.Error())
	}
	defer db.Close()

	r := NewQuotationsPostgres(db, &pq.Listener{})

	testTable := []struct {
		name             string
		codeFrom, codeTo string
		rows             *sqlmock.Rows
		expectedResult   *bwgtest.Quotation
		expectedError    error
	}{
		{
			name:     "OK",
			codeFrom: "USD",
			codeTo:   "EUR",
			rows:     sqlmock.NewRows([]string{"update_id", "code_from", "code_to", "rate", "update_time", "success"}).AddRow(1, "USD", "EUR", 1.2, time.Time{}, true),
			expectedResult: &bwgtest.Quotation{UpdateID: 1,
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       1.2,
				UpdateTime: time.Time{},
				Success:    true,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT \\* FROM quotations WHERE code_from=\\$1 AND code_to=\\$2 AND success=true ORDER BY update_id DESC LIMIT 1").
				WithArgs(testCase.codeFrom, testCase.codeTo).
				WillReturnRows(testCase.rows)

			result, err := r.GetLatestQuotation(testCase.codeFrom, testCase.codeTo)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestQuotationsPostgres_GetQuotation(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error creating mock db: %s", err.Error())
	}
	defer db.Close()

	r := NewQuotationsPostgres(db, &pq.Listener{})

	testTable := []struct {
		name           string
		updateID       int
		rows           *sqlmock.Rows
		expectedResult *bwgtest.Quotation
		expectedError  error
	}{
		{
			name:     "OK",
			updateID: 1,
			rows:     sqlmock.NewRows([]string{"update_id", "code_from", "code_to", "rate", "update_time", "success"}).AddRow(1, "USD", "EUR", 1.2, time.Time{}, true),
			expectedResult: &bwgtest.Quotation{
				UpdateID:   1,
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       1.2,
				UpdateTime: time.Time{},
				Success:    true,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT \\* FROM quotations WHERE update_id=\\$1").
				WithArgs(testCase.updateID).
				WillReturnRows(testCase.rows)

			result, err := r.GetQuotation(testCase.updateID)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestQuotationsPostgres_UpdateQuotation(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("error creating mock db: %s", err.Error())
	}
	defer db.Close()

	r := NewQuotationsPostgres(db, &pq.Listener{})

	testTable := []struct {
		name           string
		inputQuotation *bwgtest.Quotation
		expectedError  error
	}{
		{
			name: "OK",
			inputQuotation: &bwgtest.Quotation{
				UpdateID:   1,
				CodeFrom:   "USD",
				CodeTo:     "EUR",
				Rate:       1.2,
				UpdateTime: time.Time{},
				Success:    false,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mock.ExpectExec("UPDATE quotations").
				WithArgs(testCase.inputQuotation.CodeFrom, testCase.inputQuotation.CodeTo, testCase.inputQuotation.Rate,
					testCase.inputQuotation.UpdateTime, testCase.inputQuotation.Success, testCase.inputQuotation.UpdateID).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := r.UpdateQuotation(testCase.inputQuotation)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
