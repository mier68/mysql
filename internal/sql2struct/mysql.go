package sql2struct

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

type TableColumn struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnComment string
}

var DBTypeToStructType = map[string]string{
	"int":"int32",
	"tinyint":"int8",
	"smallint":"int16",
	"mediumint":"int32",
	"bigint":"int64",
	"bit":"int",
	"bool":"bool",
	"enum":"string",
	"set":"string",
	"varchar":"string",
}
func NewDBModel(Info *DBInfo) *DBModel {
	return &DBModel{
		DBInfo: Info,
	}
}

func (m *DBModel) Connect() error {
	var err error
	s := "%s:%s@tcp(%s)/information_schema?" + "charset=%s&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		s,
		m.DBInfo.UserName,
		m.DBInfo.Password,
		m.DBInfo.Host,
		m.DBInfo.Charset,
	)
	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel)GetColumns(dbName,tableName string)([]*TableColumn,error){
	query := "select column_name,data_type,column_key, "+"is_nullable,column_type,column_comment " + "from columns where table_schema = ? adn tabble_name = ?"
	rows ,err := m.DBEngine.Query(query,dbName,tableName)  
	if err != nil{
		return nil,err
	}
	if rows == nil{
		return nil,errors.New("no data")
	}
	defer rows.Close()
	var columns []*TableColumn
	for rows.Next(){
		var column TableColumn
		err := rows.Scan(&column.ColumnName,&column.DataType,&column.ColumnKey,&column.IsNullable,&column.ColumnType,&column.ColumnComment)
		if err != nil{
			return nil,err
		}
		columns = append(columns, &column)
	}
	return columns,nil
} 