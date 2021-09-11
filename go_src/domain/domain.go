package domain

import "time"

var (
	SCHOOL = &CodeType{
		Id:          "1",
		Description: "School",
	}

	CATEGORY = &CodeType{
		Id:          "2",
		Description: "Category",
	}
)

type CovidMetric struct {
	Id                 string
	SchoolId           string
	CategoryId         string
	DateTime           time.Time
	ActiveCases        int
	TotalPositiveCases int
	ProbableCases      int
	ResolvedCases      int
}

type RawDataPoint struct {
	Id       uint
	Metric   CovidMetric
	Category string
	School   string
	DateTime time.Time
}

//
//type School struct {
//	Id          string
//	Description string
//}
//
//type Category struct {
//	Id          string
//	Description string
//}

type Code struct {
	Id          string
	Description string
	Type        *CodeType
}

type CodeType struct {
	Id          string
	Description string
}

//type HASDCovidMetric struct {
//	id          uint
//	category    string
//	school      string
//	dateAndTime time.Time
//	metric      *domain.CovidMetric
//}
