package parser

import (
	"encoding/csv"
	"fmt"
	"github.com/pburskey/covid/internal/dao"
	"github.com/pburskey/covid/internal/domain"
	"github.com/pburskey/covid/internal/utility"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type prettyCSVShelf struct {
	dao           dao.DAO
	dataDirectory string
}

func BuildPrettyCSVShelf(aDao dao.DAO, aDirectory string) ShelfI {
	return &prettyCSVShelf{dao: aDao, dataDirectory: aDirectory}
}

func (me *prettyCSVShelf) Shelve(aTime time.Time, metrics []*domain.CovidMetric) error {
	if metrics != nil {

		dataTime := utility.AsYYYYMMDDHH24MiSS(aTime)

		fileName := fmt.Sprintf("covid_data_parsed_%v.csv", dataTime)

		fileMode, err := os.Stat(me.dataDirectory)
		if err != nil {
			log.Fatal("Encountered error opening: %s", me.dataDirectory, err)
		}

		if fileMode.Mode().IsRegular() {
			me.dataDirectory = filepath.Dir(me.dataDirectory)
		}

		aFilePath := filepath.Join(me.dataDirectory, fileName)

		f, err := os.Create(aFilePath)
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

		schools, _ := me.dao.GetSchools()
		categories, _ := me.dao.GetCategories()

		for _, aMetric := range metrics {
			if aMetric != nil {

				school := domain.FindCodeByID(schools, aMetric.SchoolId)

				category := domain.FindCodeByID(categories, aMetric.CategoryId)

				record = []string{dataTime, school.Description, category.Description, strconv.Itoa(aMetric.ActiveCases), strconv.Itoa(aMetric.TotalPositiveCases), strconv.Itoa(aMetric.ProbableCases), strconv.Itoa(aMetric.ResolvedCases)}

				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}

			}

		}
	}
	return nil
}
