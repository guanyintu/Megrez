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

func InitMySQL(dsn string) (Mysql, error) {
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

func (m Mysql) GetQuestions(num int, category []string) ([]Question, error) {
	var questions []Question
	if len(category) == 0 {
		category, _ = m.GetCategory()
	}
	query := "select q1.* from question q1 inner join (select (min(q2.id) + round(rand()* ( case when (max(q2.id)-?)>min(q2.id) then max(q2.id)-? - min(q2.id) else 0 end ))) as id,min(q2.id) as minId,max(q2.id) as maxId from question q2 where q2.`category` in (?)) as t on q1.id >= t.id and q1.id between t.minId and t.maxId limit ?;"
	sql, args, err := sqlx.In(query, num-1, num-1, category, num)
	if err != nil {
		return nil, err
	}
	err = m.db.Select(&questions, sql, args...)
	return questions, err
}
func (m Mysql) UpdateData(user User) error {
	sqlStr := "UPDATE user set icon=?,sum=?,success=? where uid=? "
	_, err := m.db.Exec(sqlStr, user.Icon, user.Sum, user.Success, user.Uid)
	return err
}

func (m Mysql) Profile(users []string) ([]User, error) {
	var res []User
	res = make([]User, len(users))
	for _, item := range users {
		var tmp []User
		err := m.db.Select(&tmp, "select * from user where id = ?", item)
		if err != nil {
			return nil, err
		}
		if len(tmp) == 0 {
			_, err = m.db.Exec("INSERT INTO user(uid,icon,admin,sum,success) value (?,0,0,0,0)", item)
			tmp = append(tmp, User{Uid: item, Icon: 0, Admin: 0, Sum: 0, Success: 0})
		}
		if err != nil {
			return nil, err
		}
		res = append(res, tmp[0])
	}

	return res, nil
}
func (m Mysql) GetProfile(uid string) (User, error) {
	var res []User
	err := m.db.Select(&res, "select * from user where id = ?", uid)
	if err != nil {
		return User{}, err
	}
	if len(res) == 0 {
		res = append(res, User{Uid: uid, Icon: 0, Admin: 0, Sum: 0, Success: 0})
	}
	return res[0], nil
}
