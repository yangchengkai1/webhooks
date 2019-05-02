package model

import (
	"errors"

	r "github.com/dancannon/gorethink"
)

var errTable = errors.New("table already exits")

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
				return errTable
			}
		}

	}

	return nil

}

// SelectRecord -
func SelectRecord(session *r.Session, tablename, field, value string) (interface{}, error) {
	var all []interface{}

	acursor, err := r.Table(tablename).Filter(r.Row.Field(field).Eq(value)).Run(session)

	acursor.All(&all)
	acursor.Close()

	return all, err
}

// DelateRecord -
func DelateRecord(session *r.Session, dbname, tablename, field, value string) error {
	var delate = map[string]interface{}{
		field: value,
	}
	_, err := r.DB(dbname).Table(tablename).Filter(delate).Delete().Run(session)

	return err
}

// UpdateRecord -
func UpdateRecord(session *r.Session, tablename, field, value string) error {
	var update = map[string]interface{}{
		field: value,
	}

	_, err := r.Table(tablename).Update(update).Run(session)

	return err
}
