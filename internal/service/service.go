package service

import (
	"github.com/dkshi/bwgtest"
	"github.com/dkshi/bwgtest/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Quotations interface {
	InsertQuotation(from, to string) (int, error)
	GetQuotation(updateID int) (*bwgtest.Quotation, error)
	GetLatestQuotation(codeFrom, codeTo string) (*bwgtest.Quotation, error)

	ListenQuotationUpdates()
}

type Service struct {
	Quotations
}

func NewService(repo *repository.Repostiory) *Service {
	return &Service{
		Quotations: NewQuotationsService(repo),
	}
}
