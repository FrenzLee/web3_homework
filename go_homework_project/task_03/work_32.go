package task_03

import (
	"fmt"
	"go_homework_project/task_03/repository"
)

/*
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，
并将结果映射到一个自定义的 Employee 结构体切片中。
*/
func GetEmpByDept(department string) []Employee {
	var employees []Employee
	err := repository.DB_Sqlx.Select(&employees, "select * from employees where department = ?", department)
	if err != nil {
		fmt.Println("查询失败:", err)
	}

	return employees
}

/*
使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中
*/
func GetHighestSalary() Employee {
	var employee Employee

	err := repository.DB_Sqlx.Get(&employee, "select * from employees order by salary desc limit 1")

	if err != nil {
		fmt.Println("查询出错：", err)
	}

	return employee

}

/*
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，
并将结果映射到 Book 结构体切片中，确保类型安全。
*/
func GetBookByPriceOver(price int64) []Book {
	var books []Book
	err := repository.DB_Sqlx.Select(&books, "select * from books where price > ? ", price)

	if err != nil {
		fmt.Println("查询出错：", err)
	}

	return books
}
