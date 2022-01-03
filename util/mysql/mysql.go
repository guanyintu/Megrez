package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type History struct {
	ID         int64     `db:"id"`
	Date       time.Time `db:"date"`
	QuestionId int64     `db:"questionid"`
	Uid        string    `db:"uid"`
	Success    int64     `db:"success"`
}
type Question struct {
	ID         int64     `db:"id"`
	Type       int64     `db:"type"`
	Category   string    `db:"category"`
	Question   string    `db:"question"`
	Ans        string    `db:"ans"`
	Option     string    `db:"option"`
	Reason     string    `db:"reason"`
	People     int64     `db:"people"`
	Success    int64     `db:"success"`
	Time       time.Time `db:"time"`
	Author     string    `db:"author"`
	Difficulty int64     `db:"difficulty"`
}
type User struct {
	ID      int64  `db:"id"`
	Uid     string `db:"uid"`
	Icon    int64  `db:"icon"`
	Admin   int64  `db:"admin"`
	Sum     int64  `db:"sum"`
	Success int64  `db:"success"`
}
type Mysql struct {
	db *sqlx.DB
}

func InitMySQL() (Mysql, error) {
	dsn := "mysql:asdfghjkl1+1=2@tcp(127.0.0.1:3306)/questions"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return Mysql{}, err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(5)
	return Mysql{db: db}, nil
}

func (m Mysql) GetCategory() ([]string, error) {
	var category []string
	err := m.db.Select(&category, "SELECT DISTINCT category FROM question")
	return category, err
}

func (m Mysql) GetQuestions(num int, category string) ([]Question, error) {
	var questions []Question
	query := "select q1.* from question q1 inner join (select (min(q2.id) + round(rand()* ( case when (max(q2.id)-?)>min(q2.id) then max(q2.id)-? - min(q2.id) else 0 end ))) as id,min(q2.id) as minId,max(q2.id) as maxId from question q2 where q2.`category` in ?) as t on q1.id >= t.id and q1.id between t.minId and t.maxId limit ?;"
	err := m.db.Select(&questions, query, num-1, num-1, category, num)
	return questions, err
}
func (m Mysql) UpdateData(user User) error {
	sqlStr := "UPDATE user set icon=?,sum=?,success=? where uid=? "
	_, err := m.db.Exec(sqlStr, user.Icon, user.Sum, user.Success, user.Uid)
	return err
}

//func (m Mysql) Profile() []User {
//
//}
