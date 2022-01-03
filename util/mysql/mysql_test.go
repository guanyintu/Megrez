package mysql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"testing"
)

var mock sqlmock.Sqlmock

func TestMysql(t *testing.T) {
	var err error
	var dbsql *sql.DB
	dbsql, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer dbsql.Close()

	if nil != err {
		t.Errorf("Init sqlmock failed, err %v", err)
	}
	var db *sqlx.DB
	db = sqlx.NewDb(dbsql, "mysql")
	mysql := Mysql{db: db}

	t.Run("mysql", func(t *testing.T) {
		_, err := InitMySQL()
		if err != nil {
			t.Errorf("%v", err)
			return
		}

	})
	t.Run("getGetCategory", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"category"}).
			AddRow("地理").
			AddRow("语文")
		mock.ExpectQuery("SELECT DISTINCT category FROM question").WillReturnRows(rows)
		_, err := mysql.GetCategory()
		if err != nil {
			t.Errorf("%v", err)
			return
		}
	})
	t.Run("GetQuestions", func(t *testing.T) {
		returnum := 10
		category := "[语文]"
		rows := sqlmock.NewRows([]string{"category"}).
			AddRow("地理").
			AddRow("语文")
		mock.ExpectQuery("select q1.* from question q1 inner join (select (min(q2.id) + round(rand()* ( case when (max(q2.id)-?)>min(q2.id) then max(q2.id)-? - min(q2.id) else 0 end ))) as id,min(q2.id) as minId,max(q2.id) as maxId from question q2 where q2.`category` in ?) as t on q1.id >= t.id and q1.id between t.minId and t.maxId limit ?;").
			WithArgs(returnum-1, returnum-1, category, returnum).WillReturnRows(rows)
		_, err := mysql.GetQuestions(returnum, category)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
	})
	t.Run("update", func(t *testing.T) {

	})
}
