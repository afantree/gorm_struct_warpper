package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type mockRow struct {
	ID   uint64   `json:"id" gorm:"column:id"`
	Attr demomap  `json:"attr" gorm:"column:attr;type:obj"`
	Tags demolist `json:"tags" gorm:"column:tags;type:obj"`
}

type demomap struct {
	GObj `json:"-"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (o *demomap) Scan(value interface{}) error {
	return o.GObj.Scan(value, &o)
}

func (o demomap) Value() (driver.Value, error) {
	return o.GObj.Value(o)
}

type demolist struct {
	GObj `json:"-"`
	list []string
}

func (o *demolist) At(index int) string {
	return o.list[index]
}

func (o *demolist) Scan(value interface{}) error {
	return o.GObj.Scan(value, &o.list)
}

func (o demolist) Value() (driver.Value, error) {
	return o.GObj.Value(o.list)
}

func (o *demolist) MarshalJSON() ([]byte, error) {
	a, e := json.Marshal(&o.list)
	return a, e
}

func (o *demolist) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &o.list)
}

func TestJson(t *testing.T) {
	demo := &demolist{
		GObj: GObj{},
		list: []string{"o", "k"},
	}
	b, _ := json.Marshal(demo)
	if !reflect.DeepEqual(`["o","k"]`, string(b)) {
		t.Errorf("Attr got = %v", string(b))
	}
}

func mockDB() *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}
	mock.MatchExpectationsInOrder(false)
	newRows := sqlmock.NewRows([]string{"id", "attr", "tags"}).AddRow(1, []byte(`{"name":"test","age":2}`), []byte(`["o","k"]`))
	mock.ExpectQuery("SELECT .+").WillReturnRows(newRows)

	return gdb
}

func TestObjParser(t *testing.T) {
	db := mockDB()
	t.Run("test", func(t *testing.T) {
		var row mockRow
		if err := db.Model(&mockRow{}).
			Where("id=1").
			Find(&row).Error; err != nil {
			t.Errorf("mysql select error = %v", err)
		}
		if !reflect.DeepEqual(row.Tags, demolist{
			GObj: GObj{},
			list: []string{"o", "k"},
		}) {
			t.Errorf("tags got = %v", row.Tags)
		}
		if !reflect.DeepEqual(row.Attr, demomap{
			GObj: GObj{},
			Name: "test",
			Age:  2,
		}) {
			t.Errorf("Attr got = %v", row.Attr)
		}
	})
}

func TestObjInsert(t *testing.T) {
	db := mockDB()
	r := mockRow{
		ID: 1,
		Attr: demomap{
			Name: "test",
			Age:  2,
		},
		Tags: demolist{
			list: []string{"o", "k"},
		},
	}
	stmt := db.Create(r).Statement
	sql := stmt.SQL.String()
	println(sql)
}
