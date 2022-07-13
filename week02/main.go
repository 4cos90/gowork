package main

import (
	"database/sql"
	"errors"
	"fmt"
)

//1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
var mockError = errors.New("sql.ErrNoRows")

func main() {
	_, err := GetInfoFromDB("")
	if err != nil {
		fmt.Printf("GetInfoFromDB error: %+v\n", err)
		return
	}
	fmt.Printf("GetRows Success")
}

/*
在业务代码中，查询数据库查不到数据是一种常见的现象，从业务逻辑的角度来说，进行查询操作成功了，但是查不到数据，具体要怎么处理这种情况可视业务逻辑来处理。
如果将次错误抛到上层业务代码中，那么业务代码中err != nil 时，每次还要额外进行判断是否是查不到行的错误，那么处理起来就非常繁琐了。
因此我考虑当错误是sql.ErrNoRows时，我在dao层中处理此err ，并在返回给上层时返回成功，不再把错误抛到上层。
*/
func GetInfoFromDB(sql string) (*sql.Rows, error) {
	//db := GetDbContext()
	//defer db.Close()
	rows, err := mockQuery(sql)
	//defer rows.Close()
	if err != nil {
		if errors.Is(err, mockError) {
			fmt.Printf("mockQuery error: %+v\n", err)
			return rows, nil
		}
	}
	return rows, err
}

func mockQuery(sql string) (*sql.Rows, error) {
	return nil, mockError
}

//没有数据库环境，此处假设链接成功
func GetDbContext() *sql.DB {
	db, err := sql.Open("", "")
	if err != nil {
		panic("DBOpen Fail")
	}
	err = db.Ping()
	if err != nil {
		panic("DBPing Fail")
	}
	return db
}
