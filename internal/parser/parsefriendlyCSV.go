package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/pburskey/covid/internal/dao"
	"github.com/pburskey/covid/internal/domain"
	"github.com/pburskey/covid/internal/utility"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type FriendlyCSVParser struct {
	dao        dao.DAO
	shelves    []ShelfI
	categories []*domain.Code
	schools    []*domain.Code
}

func BuildFriendlyParser(adao dao.DAO, ashelves []ShelfI) *FriendlyCSVParser {
	return &FriendlyCSVParser{shelves: ashelves, dao: adao}
}

func (me *FriendlyCSVParser) Parse(fileOrDirectory string) {

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
			aTime, metrics, _ := me.parseCSV(mungedName)

			me.shelve(aTime, metrics)

		}
	} else {
		//log.Print(fileOrDirectory)

		aTime, metrics, _ := me.parseCSV(fileOrDirectory)

		me.shelve(aTime, metrics)

	}
}

func (me *FriendlyCSVParser) shelve(aTime time.Time, metrics []*domain.CovidMetric) error {
	var err error
	if me.shelves != nil {
		for _, aShelf := range me.shelves {
			if err := aShelf.Shelve(aTime, metrics); err != nil {
				return err
			}
		}
	}
	return err
}

func (me *FriendlyCSVParser) parseDateFromFileName(fileName string) time.Time {
	start := strings.Index(fileName, "covid_data_parsed_")

	/*
		strip out the covid_data
	*/
	dateAndTimePortion := fileName[start+len("covid_data_parsed_") : (len(fileName) - len(CSV))]

	aTime, err := utility.ParseYYYYMMDDHH24MiSS(dateAndTimePortion)
	if err != nil {
		log.Fatal("Unable to parse date and time portion from file name: %s", fileName)
	}
	return aTime
}

func (me *FriendlyCSVParser) parseCSV(fileName string) (time.Time, []*domain.CovidMetric, error) {

	var metrics []*domain.CovidMetric
	aTime := me.parseDateFromFileName(fileName)

	//fmt.Println("Date and time: %s from file name: %s", aTime, fileName)
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Could not open file: %s for processing", fileName, err)
	}
	r := csv.NewReader(csvFile)
	if r == nil {
		return aTime, nil, nil
	}

	me.categories, _ = me.dao.GetCategories()
	me.schools, _ = me.dao.GetSchools()
	first := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else if first {
			first = false
		} else {
			aMetric, err := me.parseCSVRecord(aTime, record)
			if err != nil {
				return aTime, nil, err
			}
			metrics = append(metrics, aMetric)
		}

	}
	//log.Println(dataMap)

	return aTime, metrics, nil

}

func (me *FriendlyCSVParser) parseIntFromString(aString string) int {
	if anInt, err := strconv.Atoi(aString); err == nil {
		return anInt
	} else {
		log.Fatal(err)
		panic(err)
	}

}

func (me *FriendlyCSVParser) parseCSVRecord(aTime time.Time, record []string) (*domain.CovidMetric, error) {
	var err error
	if record == nil || len(record) <= 0 {
		return nil, err
	}

	//datetime,school,category,active_cases,total_positive_cases,probable_cases,resolved
	aDateTime, _ := utility.ParseYYYYMMDDHH24MiSS(record[0])

	school := domain.FindCodeByDescription(me.schools, record[1])
	if school == nil {
		school = &domain.Code{
			Id:          "",
			Description: record[1],
			Type:        domain.SCHOOL,
		}
		if _, err := me.dao.SaveSchool(school); err != nil {
			log.Fatal(err)
			return nil, err
		}
		me.schools = append(me.schools, school)
	}
	category := domain.FindCodeByDescription(me.categories, record[2])

	if category == nil {
		category = &domain.Code{
			Id:          "",
			Description: record[2],
			Type:        domain.CATEGORY,
		}
		if _, err := me.dao.SaveCategory(category); err != nil {
			log.Fatal(err)
			panic(err)
		}
		me.categories = append(me.categories, category)
	}

	aMetric := &domain.CovidMetric{
		Id:                 0,
		SchoolId:           school.Id,
		CategoryId:         category.Id,
		DateTime:           aDateTime,
		ActiveCases:        me.parseIntFromString(record[3]),
		TotalPositiveCases: me.parseIntFromString(record[4]),
		ProbableCases:      me.parseIntFromString(record[5]),
		ResolvedCases:      me.parseIntFromString(record[6]),
	}
	return aMetric, err

}
