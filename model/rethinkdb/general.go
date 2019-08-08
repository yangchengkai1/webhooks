package model

import (
	"errors"

	r "github.com/dancannon/gorethink"
)

var (
	errTableCreate = errors.New("Table already exits")
	errDBCreate    = errors.New("DB already exits")
)

// Create -
func Create(DBName, TableName string, Session *r.Session) (*r.Session, error) {
	err := CheckDB(Session, DBName)
	if err == errDBCreate {
		return Session, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = r.DBCreate(DBName).Run(Session)
	if err != nil {
		return nil, err
	}

	err = CheckTable(Session, DBName, TableName)
	if err == errTableCreate {
		return Session, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = r.DB(DBName).TableCreate(TableName).Run(Session)

	return Session, err
}

//CheckDB -
func CheckDB(session *r.Session, dbname string) error {
	var list []interface{}
	var check bool

	cursor, err := r.DBList().Run(session)
	if err != nil {
		return err
	}

	cursor.All(&list)
	cursor.Close()

	for _, db := range list {
		if !check {
			tn := db.(string)
			if tn == dbname {
				return errDBCreate
			}
		}
	}

	return nil
}

//CheckTable -
func CheckTable(session *r.Session, dbname, tablename string) error {
	var list []interface{}
	var check bool

	cursor, err := r.DB(dbname).TableList().Run(session)
	if err != nil {
		return err
	}

	cursor.All(&list)
	cursor.Close()

	for _, table := range list {
		if !check {
			tn := table.(string)
			if tn == tablename {
				return errTableCreate
			}
		}
	}

	return nil
}

// SelectRecord -
func SelectRecord(session *r.Session, DBName, TableName, field, value string) (interface{}, error) {
	var all []interface{}

	acursor, err := r.DB(DBName).Table(TableName).Filter(r.Row.Field(field).Eq(value)).Run(session)
	if err != nil {
		return nil, err
	}

	acursor.All(&all)
	acursor.Close()

	return all, nil
}

// AllRecord -
func AllRecord(session *r.Session, DBName, TableName string) (interface{}, error) {
	var all []interface{}

	acursor, err := r.DB(DBName).Table(TableName).Run(session)
	if err != nil {
		return nil, err
	}

	acursor.All(&all)
	acursor.Close()

	return all, nil
}

// DelateRecord -
func DelateRecord(session *r.Session, DBName, TableName, field, value string) error {
	var delate = map[string]interface{}{
		field: value,
	}
	_, err := r.DB(DBName).Table(TableName).Filter(delate).Delete().Run(session)

	return err
}

// UpdateRecord -
func UpdateRecord(session *r.Session, DBName, TableName, field, value string) (r.WriteResponse, error) {
	var update = map[string]interface{}{
		field: value,
	}

	return r.DB(DBName).Table(TableName).Update(update).RunWrite(session)

}

// Filter -
func Filter(session *r.Session, DBName, TableName string, filter []string) (interface{}, error) {
	var all []interface{}

	acursor, err := r.DB(DBName).Table(TableName).WithFields(filter).Run(session)
	if err != nil {
		return nil, err
	}

	acursor.All(&all)
	acursor.Close()

	return all, nil
}
