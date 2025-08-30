package task_02

import "fmt"

/*
定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
*/
type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}

func (r *Rectangle) Area() {
	fmt.Println("Rectangle中的Area()")
}

func (r *Rectangle) Perimeter() {
	fmt.Println("Rectangle中的Perimeter()")
}

type Circle struct {
}

func (c *Circle) Area() {
	fmt.Println("Circle中的Area()")
}

func (c *Circle) Perimeter() {
	fmt.Println("Circle中的Perimeter()")
}

/*
使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
*/
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	EmployeeID int
	Person     Person
}

func (e *Employee) PrintInfo() {
	fmt.Println("员工编号：", e.EmployeeID, ",员工姓名：", e.Person.Name, ",员工年龄：", e.Person.Age)
}
