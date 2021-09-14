package parser

import (
	"github.com/pburskey/covid/internal/domain"
	"time"
)

type ShelfI interface {
	Shelve(time.Time, []*domain.CovidMetric) error
}
