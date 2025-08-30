package task_03

import (
	"errors"
	"fmt"
	"go_homework_project/task_03/repository"
)

func CreateStudent() {
	stus := []Student{
		{Name: "张三", Age: 20, Grade: "三年级"},
		{Name: "李四", Age: 16, Grade: "一年级"},
		{Name: "王五", Age: 21, Grade: "四年级"},
	}

	repository.DB.Create(&stus)
}

func GetAgeOver18() {
	var stus []Student
	res := repository.DB.Where("age > ?", 18).Find(&stus)
	if res.Error != nil {
		panic(res.Error)
	}

	for i := range stus {
		fmt.Println("姓名：", stus[i].Name, "，年龄：", stus[i].Age, "，年级：", stus[i].Grade)
	}
}

func UpdateGrand(name, grade string) {
	var stu Student
	repository.DB.Where("name = ?", name).First(&stu)
	fmt.Println("修改前数据：", stu.Name, stu.Age, stu.Grade)

	repository.DB.Model(&Student{}).Where("name = ?", name).Update("grade", grade)
	repository.DB.Where("name = ?", name).First(&stu)
	fmt.Println("修改后数据：", stu.Name, stu.Age, stu.Grade)
}

func DeleteAgeSmall15() {
	var stus []Student
	repository.DB.Where("age < ?", 15).Find(&stus)
	for i := range stus {
		fmt.Println("删除前的数据：姓名：", stus[i].Name, "，年龄：", stus[i].Age, "，年级：", stus[i].Grade)
	}

	repository.DB.Where("age < ?", 15).Delete(&Student{})
	res := repository.DB.Where("age < ?", 15).Find(&stus)
	if res.RowsAffected == 0 {
		fmt.Println("没有年龄<15的数据")
	}
	for i := range stus {
		fmt.Println("删除后的数据：姓名：", stus[i].Name, "，年龄：", stus[i].Age, "，年级：", stus[i].Grade)
	}

}

/*
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。
如果余额不足，则回滚事务。
*/
func TxTransfer(fromAccountNo string, toAccountNo string, amount int64) error {
	var fromAccount, toAccount Account

	//开始事务
	tx := repository.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//查找转出、转入账户信息是否存在
	if err := tx.Where("account_no = ?", fromAccountNo).First(&fromAccount).Error; err != nil {
		tx.Rollback()
		return errors.New("查询不到转出账户信息")
	}
	if err := tx.Where("account_no = ?", toAccountNo).First(&toAccount).Error; err != nil {
		tx.Rollback()
		return errors.New("查询不到转入账户信息")
	}

	//校验转出账户余额是否充足
	if fromAccount.Balance < amount {
		tx.Rollback()
		return errors.New("转出账户余额不足")
	}

	//更新账户信息
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if err := tx.Save(&fromAccount).Error; err != nil || tx.Save(&toAccount).Error != nil {
		tx.Rollback()
		return errors.New("更新账户余额失败")
	}

	//记录交易信息
	transfer := Transaction{
		From_account_id: fromAccount.ID,
		To_account_id:   toAccount.ID,
		Amount:          amount,
	}
	if err := tx.Save(&transfer).Error; err != nil {
		tx.Rollback()
		return errors.New("记录交易信息失败")
	}

	//提交事务
	return tx.Commit().Error

}
