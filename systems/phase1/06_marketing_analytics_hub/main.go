package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/billing"
	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/crm"
	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/email"
	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/report"
)

func main() {
	// 1. Initialize the concrete underlying subsystems
	crmService := &crm.CRMService{}
	emailService := &email.EmailService{}
	billingService := &billing.BillingService{}

	// 2. Inject them into the Hub (The Facade)
	// Notice how the Hub accepts these concrete structs because they 
	// implicitly implement the Hub's defined interfaces!
	hub := report.NewReportHub(crmService, emailService, billingService)

	// 3. Create a mock context and request
	ctx := context.Background()
	req := domain.ReportRequest{
		ReportType: domain.WeeklyEngagement,
	}

	// 4. Generate the report!
	fmt.Println("Client: Requesting Weekly Engagement Report...")
	result, err := hub.GenerateReport(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// 5. Print the payload
	fmt.Printf("\n--- %s ---\n", result.Title)
	for key, val := range result.Payload {
		fmt.Printf("- %s: %s\n", key, val)
	}
}