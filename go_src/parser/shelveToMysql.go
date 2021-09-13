package parser

import (
	"github.com/pburskey/hasd_covid/dao"
	"github.com/pburskey/hasd_covid/domain"
	"time"
)

type mysqlshelf struct {
	dao dao.DAO
}

func BuildMySqlShelf(aDAO dao.DAO) ShelfI {
	return &mysqlshelf{dao: aDAO}
}

func (me *mysqlshelf) Shelve(aTime time.Time, metrics []*domain.CovidMetric) error {
	if metrics != nil {

		for _, aMetric := range metrics {
			if aMetric != nil {
				if _, err := me.dao.SaveMetric(aMetric); err != nil {
					return err
				}
			}

		}
	}
	return nil
}
