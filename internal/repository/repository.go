package repository

import (
	"github.com/dkshi/bwgtest"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Quotations interface {
	InsertQuotation(q *bwgtest.Quotation) (int, error)
	UpdateQuotation(q *bwgtest.Quotation) error
	GetQuotation(updateID int) (*bwgtest.Quotation, error)
	GetLatestQuotation(codeFrom, codeTo string) (*bwgtest.Quotation, error)

	GetQuotationChannel() chan bwgtest.Quotation
	ListenQuotationUpdates()

	InitSchema(path string) error
}

type Repostiory struct {
	Quotations
}

func NewRepository(db *sqlx.DB, lis *pq.Listener) *Repostiory {
	return &Repostiory{
		Quotations: NewQuotationsPostgres(db, lis),
	}
}
