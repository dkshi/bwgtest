package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dkshi/bwgtest"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	quotationsTable = "quotations"
)

type QuotationsPostgres struct {
	db               *sqlx.DB
	lis              *pq.Listener
	quotationChannel chan bwgtest.Quotation
}

func NewQuotationsPostgres(db *sqlx.DB, lis *pq.Listener) *QuotationsPostgres {
	newQuotationsPostgres := &QuotationsPostgres{
		db:               db,
		lis:              lis,
		quotationChannel: make(chan bwgtest.Quotation),
	}

	go newQuotationsPostgres.ListenQuotationUpdates()

	return newQuotationsPostgres
}

func (r *QuotationsPostgres) InitSchema(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading migration_up file: %s", err.Error())
	}
	initQuery := string(data)

	_, err = r.db.Exec(initQuery)

	return err
}

func (r *QuotationsPostgres) GetQuotationChannel() chan bwgtest.Quotation {
	return r.quotationChannel
}

func (r *QuotationsPostgres) InsertQuotation(q *bwgtest.Quotation) (int, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (code_from, code_to, rate, update_time, success) VALUES ($1, $2, $3, $4, $5) RETURNING update_id;", quotationsTable)
	res := r.db.QueryRow(insertQuery, q.CodeFrom, q.CodeTo, q.Rate, q.UpdateTime, q.Success)

	var updateID int

	err := res.Scan(&updateID)
	if err != nil {
		return 0, err
	}

	return updateID, nil
}

func (r *QuotationsPostgres) GetLatestQuotation(codeFrom, codeTo string) (*bwgtest.Quotation, error) {
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE code_from=$1 AND code_to=$2 AND success=true ORDER BY update_id DESC LIMIT 1;", quotationsTable)
	res := r.db.QueryRow(selectQuery, codeFrom, codeTo)

	var q bwgtest.Quotation

	err := res.Scan(&q.UpdateID, &q.CodeFrom, &q.CodeTo, &q.Rate, &q.UpdateTime, &q.Success)
	if err != nil {
		return &bwgtest.Quotation{}, err
	}

	return &q, nil
}

func (r *QuotationsPostgres) GetQuotation(updateID int) (*bwgtest.Quotation, error) {
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE update_id=$1", quotationsTable)
	res := r.db.QueryRow(selectQuery, updateID)

	var q bwgtest.Quotation

	err := res.Scan(&q.UpdateID, &q.CodeFrom, &q.CodeTo, &q.Rate, &q.UpdateTime, &q.Success)
	if err != nil {
		return &bwgtest.Quotation{}, err
	}

	return &q, nil
}

func (r *QuotationsPostgres) UpdateQuotation(q *bwgtest.Quotation) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET code_from=$1, code_to=$2, rate=$3, update_time=$4, success=$5 WHERE update_id=$6", quotationsTable)
	_, err := r.db.Exec(updateQuery, q.CodeFrom, q.CodeTo, q.Rate, q.UpdateTime, q.Success, q.UpdateID)

	return err
}

func (r *QuotationsPostgres) ListenQuotationUpdates() {
	for {
		n := <-r.lis.Notify

		var q bwgtest.Quotation

		if err := json.Unmarshal([]byte(n.Extra), &q); err != nil {
			logrus.Printf("error decoding JSON: %s", err.Error())
			continue
		}

		r.quotationChannel <- q
	}
}
