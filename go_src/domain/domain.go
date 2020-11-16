package domain

import "time"

type CovidMetric struct {
	ActiveCases        int
	TotalPositiveCases int
	ProbableCases      int
	ResolvedCases      int
}

type DataPoint struct {
	Id       uint
	Metric   CovidMetric
	Category string
	School   string
	DateTime time.Time
}

//type HASDCovidMetric struct {
//	id          uint
//	category    string
//	school      string
//	dateAndTime time.Time
//	metric      *domain.CovidMetric
//}
