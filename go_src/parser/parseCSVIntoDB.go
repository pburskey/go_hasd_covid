package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/pburskey/hasd_covid/dao"
	"github.com/pburskey/hasd_covid/domain"
	"github.com/pburskey/hasd_covid/utility"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Parser struct {
	dao dao.DAO
}

func BuildDbParser(aDAO dao.DAO) *Parser {
	return &Parser{dao: aDAO}
}

func (me *Parser) Parse(fileOrDirectory string) {

	fileMode, err := os.Stat(fileOrDirectory)
	if err != nil {
		log.Fatal("Encountered error opening: %s", fileOrDirectory, err)
	}

	if fileMode.Mode().IsDir() {
		files, err := ioutil.ReadDir(fileOrDirectory)
		if err != nil {
			log.Fatal("Encountered error reading file names in directory: %s", fileOrDirectory, err)
		}

		for _, aFileName := range files {
			log.Println(fmt.Sprintf("Processing file name: %s", aFileName.Name()))

			mungedName := (fileOrDirectory + "/" + aFileName.Name())
			me.Shelve(me.parseCSV(mungedName))
		}
	} else {
		//log.Print(fileOrDirectory)

		me.Shelve(me.parseCSV(fileOrDirectory))

	}
}

func (me *Parser) Shelve(aTime time.Time, data map[string]map[string]*domain.CovidMetric) error {
	if data != nil {

		for categoryName, schoolMap := range data {

			if _, err := me.scanForCategoriesAndSaveIfNecessary(categoryName); err != nil {
				log.Fatal(err)
				return err
			}
			categories, _ := me.dao.GetCategories()
			category := FindCategoryByDescription(categoryName, categories)
			for schoolName, metric := range schoolMap {

				if _, err := me.scanForSchoolsAndSaveIfNecessary(schoolName); err != nil {
					log.Fatal(err)
					return err
				}

				schools, _ := me.dao.GetSchools()
				school := FindSchoolByDescription(schoolName, schools)
				metric.SchoolId = school.Id
				metric.CategoryId = category.Id
				metric.DateTime = aTime
				me.dao.SaveMetric(metric)
			}
		}
	}
	return nil
}

func (me *Parser) scanForSchoolsAndSaveIfNecessary(aDescription string) ([]*domain.Code, error) {
	var schools []*domain.Code
	var err error
	if schools, err = me.dao.GetSchools(); err != nil {
		return nil, err
	}

	if school := FindSchoolByDescription(aDescription, schools); school == nil {
		school = &domain.Code{
			Description: aDescription,
		}
		if _, err = me.dao.SaveSchool(school); err != nil {
			return nil, err
		}
		schools = append(schools, school)

	}

	return schools, err

}

func (me *Parser) scanForCategoriesAndSaveIfNecessary(aDescription string) ([]*domain.Code, error) {
	var categories []*domain.Code
	var err error
	if categories, err = me.dao.GetCategories(); err != nil {
		return nil, err
	}

	if category := FindCategoryByDescription(aDescription, categories); category == nil {
		category = &domain.Code{
			Description: aDescription,
		}
		if _, err = me.dao.SaveCategory(category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, err

}

func FindSchoolByDescription(aDescription string, schools []*domain.Code) *domain.Code {
	var found *domain.Code
	for _, aSchool := range schools {
		if strings.EqualFold(aSchool.Description, aDescription) {
			found = aSchool
			break
		}
	}
	return found
}

func FindCategoryByDescription(aDescription string, objects []*domain.Code) *domain.Code {
	var found *domain.Code
	for _, anObject := range objects {
		if strings.EqualFold(anObject.Description, aDescription) {
			found = anObject
			break
		}
	}
	return found
}

func (me *Parser) parseDateFromFileName(fileName string) time.Time {
	start := strings.Index(fileName, COVID_DATA)

	/*
		strip out the covid_data
	*/
	dateAndTimePortion := fileName[start+len(COVID_DATA) : (len(fileName) - len(CSV))]

	aTime, err := utility.ParseYYYYMMDDHH24MiSS(dateAndTimePortion)
	if err != nil {
		log.Fatal("Unable to parse date and time portion from file name: %s", fileName)
	}
	return aTime
}

func (me *Parser) parseCSV(fileName string) (time.Time, map[string]map[string]*domain.CovidMetric) {

	dataMap := make(map[string]map[string]*domain.CovidMetric)
	aTime := parseDateFromFileName(fileName)

	//fmt.Println("Date and time: %s from file name: %s", aTime, fileName)
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file: %s for processing", fileName, err)
	}
	r := csv.NewReader(csvFile)
	if r == nil {
		return aTime, nil
	}

	var categories utility.Stack
	var schools utility.Stack
	first := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else if first {
			for _, aString := range record {
				if len(aString) > 0 {
					schools.Push(aString)
				}

			}
			first = false
		} else {
			parseCSVRecord(record, dataMap, &categories, &schools)
		}

	}
	//log.Println(dataMap)

	return aTime, dataMap

}

/*
	students
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)
	staff
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)




*/
//
//func parseCSVRecord(record []string, dataMap map[string]map[string]*domain.CovidMetric, categories *utility.Stack, schools *utility.Stack) {
//	if record == nil || len(record) <= 0 {
//		return
//	}
//
//	if utility.CountNonEmptyValuesIn(record) == 2 {
//		var category string = record[1]
//		_, ok := dataMap[category]
//		if !ok {
//			dataMap[category] = make(map[string]*domain.CovidMetric)
//			for _, aSchool := range *schools {
//				//fmt.Println(aSchool)
//				var metric domain.CovidMetric
//				dataMap[category][aSchool] = &metric
//			}
//		}
//		//fmt.Printf("Found category %s\n", category)
//		categories.Push(category)
//	} else {
//		category, _ := categories.Peak()
//		var metricName string = record[1]
//		//fmt.Printf("Category %s Metric: %s ... Record %s\n", category, metricName, record)
//		if metricName == "Active Cases" {
//			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
//				metric.ActiveCases = value
//			}
//			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
//
//		} else if metricName == "Total Positive Cases" {
//			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
//				metric.TotalPositiveCases = value
//			}
//			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
//		} else if metricName == "Probable Cases" {
//			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
//				metric.ProbableCases = value
//			}
//			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
//		} else if metricName == "Resolved (no longer active)" {
//			metricAssignmentFunction := func(metric *domain.CovidMetric, value int) {
//				metric.ResolvedCases = value
//			}
//			updateSchoolMetrics(record, category, dataMap, schools, metricAssignmentFunction)
//		}
//
//	}
//
//}
//
//func updateSchoolMetrics(record []string, category string, dataMap map[string]map[string]*domain.CovidMetric, schools *utility.Stack, metricAssignmentFunction func(*domain.CovidMetric, int)) {
//	i := 0
//	for _, aSchool := range *schools {
//		metric, found := dataMap[category][aSchool]
//		if !found {
//			log.Fatal("poop")
//		}
//		recordValue := record[2+i]
//		if recordValue != "" {
//			value, err := strconv.Atoi(record[2+i])
//			if err != nil {
//				log.Fatal("Encountered a non numeric value in csv data position")
//			}
//			metricAssignmentFunction(metric, value)
//		}
//
//		i++
//	}
//}

/*
Store data in redis

organization {
code, description}

category {
code, description}


organizations: [ hhs, hes, hms ...]
categories: [ staff, students ]





	staff
		school
			active cases
			total positive cases
			probable cases
			resolved ( no longer active)



Category
	School
		Date
			Metrics
*/