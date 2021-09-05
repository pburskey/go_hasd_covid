package domain

import "time"

type CovidMetric struct {
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

type CodeType struct {
	Id          uint
	Description string
}

type Code struct {
	codeType    CodeType
	id          string
	description string
}

var CODETYPE_SCHOOL = &CodeType{
	Id:          1,
	Description: "School",
}
var CODETYPE_CATEGORY = &CodeType{
	Id:          2,
	Description: "Category",
}

//type HASDCovidMetric struct {
//	id          uint
//	category    string
//	school      string
//	dateAndTime time.Time
//	metric      *domain.CovidMetric
//}
