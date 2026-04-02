package domain

type ReportType int
const (
	WeeklyEngagement ReportType = iota
	MonthlyRevenue
	QuarterlyChurn
)
type ReportRequest struct{
	ReportType ReportType
}
type Report struct{
	Title string
	Payload map[string]string
}