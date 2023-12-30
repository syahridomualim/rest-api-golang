package repositories

import (
	"context"
	"database/sql"
	"employee-golang/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func Test_repositories_GetEmployee(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name        string
		fields      fields
		wantRs      []*model.Employee
		expectedErr error
		wantErr     bool
	}{
		{
			name: `get all employees`,
			wantRs: []*model.Employee{
				&model.Employee{
					IdEmployee: "1",
					FirstName:  "John",
					LastName:   "Doe",
					Email:      "john.doe@example.com",
					Phone:      "123456789",
					HireDate:   "2023-01-01",
					Salary:     50000.0,
				},
				&model.Employee{
					IdEmployee: "2",
					FirstName:  "Jane",
					LastName:   "Doe",
					Email:      "jane.doe@example.com",
					Phone:      "987654321",
					HireDate:   "2023-01-02",
					Salary:     60000.0,
				},
			},
			wantErr: false,
		},
		{
			name:        "error query",
			wantRs:      nil,
			expectedErr: errors.New("some database error"),
			wantErr:     true,
		},
		{
			name:        "error database",
			wantRs:      nil,
			expectedErr: sql.ErrConnDone,
			wantErr:     true,
		},
	}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal(`failed to create mock`)
	}
	defer db.Close()

	mock.
		ExpectQuery("select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee").
		WillReturnRows(
			sqlmock.NewRows([]string{"employee_id", "first_name", "last_name", "email", "phone", "hire_date", "salary"}).
				AddRow(1, `John`, `Doe`, `john.doe@example.com`, `123456789`, `2023-01-01`, 50000.0).
				AddRow(2, `Jane`, `Doe`, `jane.doe@example.com`, `987654321`, `2023-01-02`, 60000.0))

	mock.
		ExpectQuery("select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee").
		WillReturnError(tests[1].expectedErr)

	mock.
		ExpectQuery("select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee").
		WillReturnError(tests[2].expectedErr)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.DB = db
			c := repositories{
				DB: tt.fields.DB,
			}
			gotRs, err := c.GetEmployee()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmployee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("GetEmployee() gotRs = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}

func Test_repositories_GetEmployeeById(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRs  *model.Employee
		expErr  error
		wantErr bool
	}{
		{
			name: "success get employee",
			args: args{
				"1",
			},
			wantRs: &model.Employee{
				IdEmployee: "1",
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john.doe@example.com",
				Phone:      "123456789",
				HireDate:   "2023-01-01",
				Salary:     50000.0,
			},
			wantErr: false,
		},
		// Add more test cases as needed
		{
			name: "not found employee",
			args: args{
				"3",
			},
			wantRs:  nil,
			expErr:  sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "database error",
			args: args{
				"2",
			},
			wantRs:  nil,
			expErr:  errors.New("some database error"),
			wantErr: true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal("Error creating mock database: ", err)
	}
	defer db.Close()

	for _, tt := range tests {
		if !tt.wantErr {
			rows := sqlmock.NewRows([]string{"employee_id", "first_name", "last_name", "email", "phone", "hire_date", "salary"}).
				AddRow(tt.wantRs.IdEmployee, tt.wantRs.FirstName, tt.wantRs.LastName, tt.wantRs.Email, tt.wantRs.Phone, tt.wantRs.HireDate, tt.wantRs.Salary)
			mock.
				ExpectQuery("select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee where employee_id = ?").
				WithArgs(tt.args.id).
				WillReturnRows(rows)
		} else {
			mock.ExpectQuery("select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee where employee_id = ?").
				WithArgs(tt.args.id).
				WillReturnError(tt.expErr)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.DB = db
			c := repositories{
				DB: tt.fields.DB,
			}
			gotRs, err := c.GetEmployeeById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmployeeById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("GetEmployeeById() gotRs = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}

func TestNewEmployeeRepositories(t *testing.T) {
	tests := []struct {
		name string
		want IEmployeeRepositories
	}{
		// TODO: Add test cases.
		{
			name: "success",
			want: NewEmployeeRepositories(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEmployeeRepositories(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEmployeeRepositories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repositories_InsertEmployee(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		employee *model.Employee
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRs      string
		expectedErr error
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name: "insert new employee",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantRs:  "Successfully inserted a new employee",
			wantErr: false,
		},
		{
			name: "employee already exists",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantRs:      "Employee already exists",
			expectedErr: errors.New("employee already exists"),
			wantErr:     true,
		},
		{
			name: "error while inserting data",
			args: args{
				employee: &model.Employee{
					IdEmployee: "1235",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido3@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantRs:      "",
			expectedErr: errors.New("some database error"),
			wantErr:     true,
		},
		{
			name: "error database",
			args: args{
				employee: &model.Employee{
					IdEmployee: "1234",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido1@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantRs:      "",
			expectedErr: sql.ErrConnDone,
			wantErr:     true,
		},
		{
			name: "error checking employee existence",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantRs:      "",
			expectedErr: errors.New("error checking employee existence"),
			wantErr:     true,
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal("error creating mock")
	}
	defer db.Close()

	for _, tt := range tests {
		if !tt.wantErr {
			mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
				WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
				WillReturnRows(sqlmock.NewRows([]string{" count(*)"}).
					AddRow(0))
			mock.ExpectExec("insert into employee (employee_id, first_name, last_name, email, phone, hire_date, salary) value (?, ?, ?, ?, ?, nullif(?,''), nullif(?, ''))").
				WithArgs(
					tt.args.employee.IdEmployee,
					tt.args.employee.FirstName,
					tt.args.employee.LastName,
					tt.args.employee.Email,
					tt.args.employee.Phone,
					tt.args.employee.HireDate,
					tt.args.employee.Salary).
				WillReturnResult(sqlmock.NewResult(1, 1))
		} else {
			if tt.wantRs == "Employee already exists" {
				mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
					WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
					WillReturnRows(sqlmock.NewRows([]string{" count(*)"}).
						AddRow(1))
			} else {
				if tt.expectedErr.Error() == "error checking employee existence" {
					// Simulate an error during the employee existence check
					mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
						WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
						WillReturnError(errors.New("error checking employee existence"))
				}
				mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
					WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
					WillReturnRows(sqlmock.NewRows([]string{" count(*)"}).
						AddRow(0))
				mock.ExpectExec("insert into employee (employee_id, first_name, last_name, email, phone, hire_date, salary) value (?, ?, ?, ?, ?, nullif(?,''), nullif(?, ''))").
					WithArgs(
						tt.args.employee.IdEmployee,
						tt.args.employee.FirstName,
						tt.args.employee.LastName,
						tt.args.employee.Email,
						tt.args.employee.Phone,
						tt.args.employee.HireDate,
						tt.args.employee.Salary).
					WillReturnError(tt.expectedErr)
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.DB = db
			c := repositories{
				DB: tt.fields.DB,
			}
			gotRs, err := c.InsertEmployee(tt.args.employee)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertEmployee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRs != tt.wantRs {
				t.Errorf("InsertEmployee() gotRs = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}

func Test_repositories_EditEmployee(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		employee *model.Employee
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantRs    string
		expectErr error
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:   "success update",
			wantRs: "Employee was edited",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantErr: false,
		},
		{
			name:   "employee doesnt exist",
			wantRs: "Employee doesn't exists",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantErr:   true,
			expectErr: sql.ErrNoRows,
		},
		{
			name:   "employee doesnt exist",
			wantRs: "",
			args: args{
				employee: &model.Employee{
					IdEmployee: "123",
					FirstName:  "Mualim",
					LastName:   "Syahrido",
					Email:      "syahrido@gmail.com",
					Phone:      "3424235",
					HireDate:   time.Now().String(),
					Salary:     120000.0,
				},
			},
			wantErr:   true,
			expectErr: errors.New("database was error"),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal("error creating mock")
	}
	defer db.Close()

	for _, tt := range tests {
		if !tt.wantErr {
			mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
				WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
				WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
					AddRow(1))
			mock.ExpectExec("update employee set first_name = ?, last_name = ?, email = ?, phone = ?, hire_date = nullif(?, ''), salary = nullif(?, 0) where employee_id = ?").
				WithArgs(
					tt.args.employee.FirstName,
					tt.args.employee.LastName,
					tt.args.employee.Email,
					tt.args.employee.Phone,
					tt.args.employee.HireDate,
					tt.args.employee.Salary,
					tt.args.employee.IdEmployee).
				WillReturnResult(sqlmock.NewResult(1, 1))
		} else {
			mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
				WithArgs(tt.args.employee.IdEmployee, tt.args.employee.Email).
				WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
					AddRow(0))
			mock.ExpectExec("update employee set first_name = ?, last_name = ?, email = ?, phone = ?, hire_date = nullif(?, ''), salary = nullif(?, 0) where employee_id = ?").
				WithArgs(
					tt.args.employee.FirstName,
					tt.args.employee.LastName,
					tt.args.employee.Email,
					tt.args.employee.Phone,
					tt.args.employee.HireDate,
					tt.args.employee.Salary,
					tt.args.employee.IdEmployee).
				WillReturnError(tt.expectErr)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.DB = db
			r := repositories{
				DB: tt.fields.DB,
			}
			gotRs, err := r.UpdateEmployee(tt.args.employee)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditEmployee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRs != tt.wantRs {
				t.Errorf("EditEmployee() gotRs = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}

func Test_repositories_employeeExists(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx        context.Context
		idEmployee string
		email      string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		expectErr error
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:    "test error",
			wantErr: true,
			want:    false,
			args: args{
				ctx:        context.Background(),
				idEmployee: "2345",
				email:      "mualim@data.com",
			},
			expectErr: errors.New("DB error"),
		},
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal("error creating mock")
	}
	defer db.Close()

	mock.ExpectQuery("select count(*) from employee where employee_id = ? or email = nullif(?, '')").
		WithArgs(tests[0].args.idEmployee, tests[0].args.email).
		WillReturnError(errors.New("DB error"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.DB = db
			r := repositories{
				DB: tt.fields.DB,
			}
			got, err := r.employeeExists(tt.args.ctx, &tt.args.idEmployee, &tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("employeeExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("employeeExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
