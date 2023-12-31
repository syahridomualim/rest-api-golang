package config

import "github.com/spf13/viper"

func GetConnection() string {
	v := viper.GetString("datasource.employee.connection")
	if v == "" {
		return ""
	}
	return v
}

func GetEmployees() string {
	v := viper.GetString("app.query.GET_EMPLOYEES")
	if v == "" {
		return "select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee"
	}
	return v
}

func GetEmployeeById() string {
	v := viper.GetString("app.query.GET_EMPLOYEES_BY_ID")
	if v == "" {
		return "select employee_id, first_name, last_name, email, phone, coalesce(hire_date, ''), coalesce(salary, 0.0) from employee where employee_id = ?"
	}
	return v
}

func InsertEmployee() string {
	v := viper.GetString("app.query.INSERT_EMPLOYEE")
	if v == "" {
		return "insert into employee (employee_id, first_name, last_name, email, phone, hire_date, salary) value (?, ?, ?, ?, ?, nullif(?,''), nullif(?, ''))"
	}
	return v
}

func CountEmployee() string {
	v := viper.GetString("app.query.COUNT_EMPLOYEE")
	if v == "" {
		return "select count(*) from employee where employee_id = ? or email = nullif(?, '')"
	}
	return v
}

func EditEmployee() string {
	v := viper.GetString("app.query.EDIT_EMPLOYEE")
	if v == "" {
		return "update employee set first_name = ?, last_name = ?, email = ?, phone = ?, hire_date = nullif(?, ''), salary = nullif(?, 0) where employee_id = ?"
	}
	return v
}

func DeleteEmployee() string {
	v := viper.GetString("app.query.DELETE_EMPLOYEE")
	if v == "" {
		return "delete from employee where employee_id = ?"
	}
	return v
}
