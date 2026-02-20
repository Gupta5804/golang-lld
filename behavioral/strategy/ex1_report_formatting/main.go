package main

import (
	"encoding/json"
	"errors"
	"strings"
)

// the strategy interface (slot)
type FormatStrategy interface {
	Format(data []string) string
}

// concrete strategies (catridges)
type JSONStrategy struct{}

func (j *JSONStrategy) Format(data []string) string {
	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}

type CSVStrategy struct{}

func (c *CSVStrategy) Format(data []string) string {
	return strings.Join(data, ",")
}

// the context (console)
type ReportGenerator struct {
	strategy FormatStrategy
}

func NewReportGenerator(strategy FormatStrategy) *ReportGenerator {
	return &ReportGenerator{
		strategy: strategy,
	}
}

func (r *ReportGenerator) Generate(data []string) (string, error) {
	if r.strategy == nil {
		return "", errors.New("No strategy")
	}
	return r.strategy.Format(data),nil
}
