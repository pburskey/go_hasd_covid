package parser

import (
	"github.com/pburskey/hasd_covid/domain"
	"time"
)

type ShelfI interface {
	Shelve(time.Time, []*domain.CovidMetric) error
}
