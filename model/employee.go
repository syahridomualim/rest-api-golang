package model

type Employee struct {
	IdEmployee string  `json:"idEmployee,omitempty" db:"id" validate:"required"`
	FirstName  string  `json:"firstName,omitempty" db:"first_name" validate:"required"`
	LastName   string  `json:"lastName,omitempty" db:"last_name" validate:"required"`
	Email      string  `json:"email,omitempty" db:"email" validate:"required"`
	Phone      string  `json:"phone,omitempty" db:"phone" validate:"required"`
	HireDate   string  `json:"hireDate" db:"hire_date"`
	Salary     float64 `json:"salary" db:"salary"`
}
