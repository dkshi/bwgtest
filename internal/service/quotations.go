package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dkshi/bwgtest"
	"github.com/dkshi/bwgtest/internal/repository"
	"github.com/sirupsen/logrus"
)

const (
	exchangeHost         = "https://api.frankfurter.app/latest"
)

type Conversion struct {
	Base    string             `json:"base"`
	Rates   map[string]float32 `json:"rates"`
	Message string             `json:"message"`
}

type QuotationsService struct {
	repo *repository.Repostiory
}

func NewQuotationsService(repo *repository.Repostiory) *QuotationsService {
	newQuotationsService := &QuotationsService{
		repo: repo,
	}

	go newQuotationsService.ListenQuotationUpdates()

	return newQuotationsService
}

func (s *QuotationsService) InsertQuotation(codeFrom, codeTo string) (int, error) {
	newQuotation := &bwgtest.Quotation{
		CodeFrom: codeFrom,
		CodeTo:   codeTo,
		Success:  false,
	}

	return s.repo.InsertQuotation(newQuotation)
}

func (s *QuotationsService) GetQuotationRate(codeFrom, codeTo string) (float32, error) {
	url := fmt.Sprintf("%s?from=%s&to=%s", exchangeHost, codeFrom, codeTo)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	var conversion Conversion
	err = json.NewDecoder(resp.Body).Decode(&conversion)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("error from api server: %s", conversion.Message)
	}

	rate, ok := conversion.Rates[codeTo]
	if !ok {
		return 0, fmt.Errorf("error getting rate from map")
	}

	return rate, nil
}

func (s *QuotationsService) GetQuotation(updateID int) (*bwgtest.Quotation, error) {
	return s.repo.GetQuotation(updateID)
}

func (s *QuotationsService) UpdateQuotation(q *bwgtest.Quotation) error {
	q.UpdateTime = time.Now()

	rate, err := s.GetQuotationRate(q.CodeFrom, q.CodeTo)
	if err != nil {
		logrus.Printf("error getting quotation rate: %s", err.Error())
		return s.repo.UpdateQuotation(q)
	}

	q.Rate = rate
	q.Success = true

	return s.repo.UpdateQuotation(q)
}

func (s *QuotationsService) GetLatestQuotation(codeFrom, codeTo string) (*bwgtest.Quotation, error) {
	return s.repo.GetLatestQuotation(codeFrom, codeTo)
}

func (s *QuotationsService) ListenQuotationUpdates() {
	quotationChannel := s.repo.GetQuotationChannel()
	for {
		q := <-quotationChannel

		go func(q *bwgtest.Quotation) {
			err := s.UpdateQuotation(q)
			if err != nil {
				logrus.Printf("error updating quotation: %s", err.Error())
			}
		}(&q)
	}
}
