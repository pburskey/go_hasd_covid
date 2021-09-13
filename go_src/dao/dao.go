package dao

import "github.com/pburskey/hasd_covid/domain"

type DAO interface {
	GetSchools() ([]*domain.Code, error)
	GetCategories() ([]*domain.Code, error)
	GetMetric(uint) (*domain.CovidMetric, error)
	GetMetricsBySchool(*domain.Code) ([]*domain.CovidMetric, error)
	SaveMetric(*domain.CovidMetric) (*domain.CovidMetric, error)
	SaveSchool(*domain.Code) (*domain.Code, error)
	SaveCategory(*domain.Code) (*domain.Code, error)
}
