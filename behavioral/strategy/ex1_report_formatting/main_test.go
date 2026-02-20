package main

import "testing"

func TestReportGenerator_JSON(t *testing.T) {
	json := &JSONStrategy{}
	input := []string{"apple", "banana"}
	generator := NewReportGenerator(json)

	output, err := generator.Generate(input)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if output != `["apple","banana"]` {
		t.Errorf("Expected ['apple','banana'], got %s", output)
	}
}

func TestReportGenerator_CSV(t *testing.T) {
	csv := &CSVStrategy{}
	input := []string{"apple", "banana"}
	generator := NewReportGenerator(csv)

	output, err := generator.Generate(input)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if output != "apple,banana" {
		t.Errorf("Expected 'apple,banana', got %s", output)
	}
}

func TestReportGenerator_NoStrategy(t *testing.T) {
	input := []string{"apple", "banana"}
	generator := &ReportGenerator{}

	_,err := generator.Generate(input)
	if err == nil {
		t.Error("Expected error, instead got nil")
	}
}
