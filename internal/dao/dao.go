package dao

import "github.com/pburskey/covid/internal/domain"

type DAO interface {
	GetSchools() ([]*domain.Code, error)
	GetCategories() ([]*domain.Code, error)
	GetMetric(int) (*domain.CovidMetric, error)
	GetMetricsBySchool(*domain.Code) ([]*domain.CovidMetric, error)
	SaveMetric(*domain.CovidMetric) (*domain.CovidMetric, error)
	SaveSchool(*domain.Code) (*domain.Code, error)
	SaveCategory(*domain.Code) (*domain.Code, error)
	GetMetricsByCategory(*domain.Code) ([]*domain.CovidMetric, error)
	GetMetricsBySchoolAndCategory(school *domain.Code, category *domain.Code) ([]*domain.CovidMetric, error)
}
