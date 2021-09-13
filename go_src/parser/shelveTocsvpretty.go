package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/pburskey/hasd_covid/dao"
	"github.com/pburskey/hasd_covid/domain"
	"github.com/pburskey/hasd_covid/utility"
	"log"
	"os"
	"strconv"
	"time"
)

type prettyCSVShelf struct {
	dao dao.DAO
}

func BuildPrettyCSVShelf(aDao dao.DAO) ShelfI {
	return &prettyCSVShelf{dao: aDao}
}

func (me *prettyCSVShelf) Shelve(aTime time.Time, metrics []*domain.CovidMetric) error {
	if metrics != nil {

		dataTime := utility.AsYYYYMMDDHH24MiSS(aTime)
		fileName := fmt.Sprintf("covid_data_parsed_%v.csv", dataTime)

		f, err := os.Create(fileName)
		if err != nil {
			log.Fatalln("failed to open file", err)
		}
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		record := []string{"datetime", "school", "category", "active_cases", "total_positive_cases", "probable_cases", "resolved"}

		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}

		for _, aMetric := range metrics {
			if aMetric != nil {

				school, err := me.dao.FindSchoolBy(aMetric.SchoolId)
				if err != nil {

				}
				category, err := me.dao.FindCategoryBy(aMetric.CategoryId)
				if err != nil {

				}
				record = []string{dataTime, school.Description, category.Description, strconv.Itoa(aMetric.ActiveCases), strconv.Itoa(aMetric.TotalPositiveCases), strconv.Itoa(aMetric.ProbableCases), strconv.Itoa(aMetric.ResolvedCases)}

				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}

			}

		}
	}
	return nil
}
