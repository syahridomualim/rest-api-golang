package repositories

import (
	"context"
	"database/sql"
	"employee-golang/config"
	"employee-golang/model"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type repositories struct {
	DB *sql.DB
}

var (
	EmployeeDB             *sql.DB
	MutexConfigurationConn = &sync.Mutex{}
)

func NewEmployeeRepositories() IEmployeeRepositories {
	return InitConfiguration()
}
func InitConfiguration() *repositories {
	// connect db
	MutexConfigurationConn.Lock()
	if EmployeeDB == nil {
		configurationDB, err := sql.Open("mysql", config.GetConnection())
		if err != nil {
			logrus.Errorf("failed to create configuration connection %s", err)
		}
		err = configurationDB.Ping()
		if err != nil {
			logrus.Error(err)
		}
		configurationDB.SetConnMaxLifetime(time.Minute * 3)
		configurationDB.SetMaxOpenConns(10)
		configurationDB.SetMaxIdleConns(10)
		EmployeeDB = configurationDB
	}
	MutexConfigurationConn.Unlock()
	return &repositories{
		DB: EmployeeDB,
	}
}

type IEmployeeRepositories interface {
	GetEmployee() (rs []*model.Employee, err error)
	GetEmployeeById(id string) (rs *model.Employee, err error)
	InsertEmployee(employee *model.Employee) (rs string, err error)
	UpdateEmployee(employee *model.Employee) (rs string, err error)
	DeleteEmployee(id string) (rs string, err error)
}

func (r repositories) GetEmployee() (rs []*model.Employee, err error) {
	res := make([]*model.Employee, 0)
	query := config.GetEmployees()
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		data := new(model.Employee)
		err := rows.Scan(
			&data.IdEmployee,
			&data.FirstName,
			&data.LastName,
			&data.Email,
			&data.Phone,
			&data.HireDate,
			&data.Salary,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		logrus.Println(data)
		res = append(res, data)
	}
	return res, nil
}

func (r repositories) GetEmployeeById(id string) (rs *model.Employee, err error) {
	query := config.GetEmployeeById()
	data := &model.Employee{}

	err = r.DB.QueryRowContext(context.Background(), query, id).Scan(
		&data.IdEmployee,
		&data.FirstName,
		&data.LastName,
		&data.Email,
		&data.Phone,
		&data.HireDate,
		&data.Salary,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		logrus.Errorf("Employee %v not found", id)
		return nil, err
	case err != nil:
		logrus.Errorf("Error retrieving employee: %v", err)
		return nil, err
	default:
		logrus.Printf("Employee found: %+v", data)
	}
	return data, nil
}

func (r repositories) InsertEmployee(employee *model.Employee) (rs string, err error) {
	queryInsert := config.InsertEmployee()
	ctx := context.Background()
	// check if the employee with same ID or email already exists
	exists, err := r.employeeExists(ctx, &employee.IdEmployee, &employee.Email)
	if err != nil {
		logrus.Errorf("Error checking employee existence: %v", err)
		return "", err
	}
	if exists {
		return "Employee already exists", errors.New("employee already exists")
	}

	_, err = r.DB.ExecContext(
		ctx, queryInsert,
		employee.IdEmployee, employee.FirstName, employee.LastName,
		employee.Email, employee.Phone, &employee.HireDate, &employee.Salary)
	switch {
	case errors.Is(err, sql.ErrConnDone):
		logrus.Errorf("Error inserting employee: %v", err)
		return "", err
	case err != nil:
		logrus.Errorf("Error inserting employee: %v", err)
		return "", err
	default:
		logrus.Infof("successfully insert new employee %s", rs)
	}
	return "Successfully inserted a new employee", nil
}

func (r repositories) UpdateEmployee(employee *model.Employee) (rs string, err error) {
	ctx := context.Background()
	exists, err := r.employeeExists(ctx, &employee.IdEmployee, &employee.Email)
	if err != nil {
		logrus.Errorf("Error checking employee existence: %v", err)
		return "", err
	}
	if !exists {
		return "Employee doesn't exists", errors.New("employee doesn't exists")
	}
	query := config.EditEmployee()
	_, err = r.DB.ExecContext(ctx, query,
		employee.FirstName, employee.LastName, employee.Email,
		employee.Phone, &employee.HireDate, &employee.Salary, employee.IdEmployee)
	switch {
	case err != nil:
		logrus.Errorf("Error on database %v", err)
		return "", err
	default:
		logrus.Infof("Employee was edited")
	}
	return "Employee was edited", err
}

func (r repositories) DeleteEmployee(employeeId string) (rs string, err error) {
	ctx := context.Background()
	exists, err := r.employeeExists(ctx, &employeeId, nil)
	if err != nil {
		logrus.Errorf("Error checking employee existence: %v", err)
		return "", err
	}
	if !exists {
		return "Employee doesn't exists", errors.New("employee doesn't exists")
	}
	query := config.DeleteEmployee()
	_, err = r.DB.ExecContext(ctx, query, employeeId)
	switch {
	case err != nil:
		logrus.Errorf("Error on database %v", err)
		return "", err
	default:
		logrus.Infof("Employee was deleted")
	}
	return "Employee was deleted", err
}

func (r repositories) employeeExists(ctx context.Context, idEmployee, email *string) (bool, error) {
	var count int
	query := config.CountEmployee()
	err := r.DB.QueryRowContext(ctx, query, idEmployee, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
