package statements

import (
	"fmt"
	"errors"
	"context"
	"database/sql"
)

// ERRORS
var (
	ErrValueIsNotDB = errors.New("db value is not *sql.DB")
)

/* 
NewGetAllStmt return get stmt with 1 field that find all record

ctx should be withValue "db" *sql.DB
*/ 
func NewGetAllStmt(ctx context.Context, table, field string) (*sql.Stmt, error) {
	return NewGetStmt(ctx, "*", table, field)
}


/* 
NewGetStmt return get stmt with 1 field that find one field of record

ctx should be withValue "db" *sql.DB
*/ 
func NewGetStmt(ctx context.Context, get, table, field string) (*sql.Stmt, error) {
	db, ok := ctx.Value("db").(*sql.DB)
	if !ok {
		return nil, ErrValueIsNotDB
	}
	query := fmt.Sprintf("select %s from %s where %s = ?", get ,table, field)
	return db.Prepare(query)
}
/*
NewDeleteStmt return delete stms with 1 field that delete in table record with this field
	params: ctx - with key "db" - *sql.DB, table and field -
	return: *sql.Stms and  error
	errors: "db value is not *sql.DB"
*/
func NewDeleteStmt(ctx context.Context, table, field string) (*sql.Stmt, error) {
	db, ok := ctx.Value("db").(*sql.DB)
	if !ok {
		return  nil, ErrValueIsNotDB
	}

	query := fmt.Sprintf("delete from %s where %s = ?", table, field)
	return db.Prepare(query)
}

/*
NewInsertStmt return delete stms with fields that insert into table values
	params: ctx - with key "db" - *sql.DB, table and fields - string
	return: *sql.Stms and  error
	errors: "db value is not *sql.DB" = ErrValueIsNotDB
	example: "INSERT INTO table (field, field) VALUES (args,,,)"
*/
func NewInsertStmt(
	ctx context.Context, 
	table string, 
	fields ...interface{}, 
	) (*sql.Stmt, error) {
		db, ok := ctx.Value("db").(*sql.DB)
		if !ok {
			return  nil, ErrValueIsNotDB
		}
		
		formatFunc := func(n int) string {
			str := "("

			for i := 0 ;; i++ {
				if i == n - 1 {
					str += "%s)"
					break
				}
				str += "%s, "
			}
			// now str look like: (field, field, ...., field)

			str += " values ("

			for i := 0 ;; i++ {
				if i == n - 1 {
					str += "$"+fmt.Sprint(i+1)+")"
					break
				}
				str += "$"+fmt.Sprint(i+1)+", "
			}

			return str
		}

		query := fmt.Sprintf("insert into "+table + formatFunc( len(fields) ), fields...)
		return db.Prepare(query)
		
}