package tools

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func ExportSwaData() {
	var err error
	ex := SwaConfig.Export
	tb := SwaConfig.Tables
	var dsn string
	switch ex.DbType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ex.Username, ex.Password, ex.Ip, ex.Port, ex.Dbname)
	case "pgsql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", ex.Ip, ex.Username, ex.Password, ex.Dbname, ex.Port, ex.Config)
	default:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", ex.Username, ex.Password, ex.Ip, ex.Port, ex.Dbname, ex.Config)
	}

	var db *sql.DB
	switch ex.DbType {
	case "mysql":
		db, err = sql.Open("mysql", dsn)
	case "pgsql":
		db, err = sql.Open("postgres", dsn)
	default:
		db, err = sql.Open("mysql", dsn)
	}
	if err != nil {
		fmt.Printf("获取数据库sqlDB错误：%v\n", err)
		return
	}
	if db == nil {
		fmt.Printf("获取数据库sqlDB失败:%v\n", db)
		return
	} else {
		fmt.Printf("获取数据库sqlDB成功\n")
	}

	err = db.Ping()
	if err != nil {
		println("失败! 数据库ping测试失败")
	} else {
		println("数据库ping测试成功")
	}

	fmt.Printf("需导出的表名称：%s\n", tb.TbString)
	tbList := strings.Split(tb.TbString, ",")

	cwdPath, _ := os.Getwd()
	savePath := fmt.Sprintf("%s/tools/export_file", cwdPath)

	_, err = os.Stat(savePath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(savePath, os.ModePerm)
		if err != nil {
			fmt.Printf("失败! 创建文件夹失败：%s\n", savePath)
			return
		}
	}

	for _, oneTb := range tbList {
		fmt.Printf("导出表：%s\n", oneTb)
		rows, err := db.Query(fmt.Sprintf("SELECT * from %s", oneTb))

		if err != nil {
			panic(err)
		}

		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error())
		}

		values := make([]sql.RawBytes, len(columns))

		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		totalValues := make([][]string, 0)
		for rows.Next() {

			var s []string

			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err.Error())
			}

			var rowMap []string
			for idx, dict := range values {
				if columns[idx] != "id" &&
					columns[idx] != "created_at" &&
					columns[idx] != "updated_at" &&
					columns[idx] != "deleted_at" {
					oneD := fmt.Sprintf("%s:\"%s\"", columns[idx], string(dict))
					rowMap = append(rowMap, fmt.Sprintf("%s", oneD))
				}
			}
			for _, v := range values {
				s = append(s, string(v))
			}
			totalValues = append(totalValues, s)
		}

		if err = rows.Err(); err != nil {
			panic(err.Error())
		}
		filePath := fmt.Sprintf("%s/%s.csv", savePath, oneTb)
		fmt.Printf("保存路径：%s\n", filePath)
		writeToCSV(filePath, columns, totalValues)
		fmt.Printf("%s处理完毕\n\n", oneTb)
	}

}

func writeToCSV(file string, columns []string, totalValues [][]string) {
	f, err := os.Create(file)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	for i, row := range totalValues {
		if i == 0 {
			err := w.Write(columns)
			if err != nil {
				return
			}
			err = w.Write(row)
			if err != nil {
				return
			}
		} else {
			err := w.Write(row)
			if err != nil {
				return
			}
		}
	}
	w.Flush()
}
