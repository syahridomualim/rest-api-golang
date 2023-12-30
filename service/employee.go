package service

import (
	"employee-golang/model"
	"employee-golang/repositories"
	"github.com/sirupsen/logrus"
)

type service struct {
	repository repositories.IEmployeeRepositories
}

func NewEmployeeService() IEmployeeService {
	return &service{
		repository: repositories.NewEmployeeRepositories(),
	}
}

type IEmployeeService interface {
	GetEmployees() (rs []*model.Employee, err error)
	GetEmployeeById(id string) (rs *model.Employee, err error)
	InsertEmployee(employee *model.Employee) (rs string, err error)
	UpdateEmployee(employee *model.Employee) (rs string, err error)
}

func (s service) InsertEmployee(employee *model.Employee) (rs string, err error) {
	rs, err = s.repository.InsertEmployee(employee)
	if err != nil {
		logrus.Error("Error is been occurred")
		return "", err
	}
	return rs, nil
}

func (s service) UpdateEmployee(employee *model.Employee) (rs string, err error) {
	rs, err = s.repository.UpdateEmployee(employee)
	if err != nil {
		logrus.Error("Error is been occurred")
		return "", err
	}
	return rs, nil
}

func (s service) GetEmployees() (rs []*model.Employee, err error) {
	rs, err = s.repository.GetEmployee()
	if err != nil {
		logrus.Error("Error is been occurred")
		return nil, err
	}
	return rs, nil
}

func (s service) GetEmployeeById(id string) (rs *model.Employee, err error) {
	rs, err = s.repository.GetEmployeeById(id)
	if err != nil {
		logrus.Error("Error is been occurred")
		return nil, err
	}
	return rs, nil
}
