package tools

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)


func InitSwaData() {
	var err error
	im := SwaConfig.Import
	tbs := SwaConfig.Tables



	var dsn string
	switch im.DbType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", im.Username, im.Password, im.Ip, im.Port, im.Dbname)
	case "pgsql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", im.Ip, im.Username, im.Password, im.Dbname, im.Port, im.Config)
	default:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", im.Username, im.Password, im.Ip, im.Port, im.Dbname, im.Config)
	}

	var db *sql.DB
	switch im.DbType {
	case "mysql":
		db, err = sql.Open("mysql", dsn)
	case "pgsql":
		db, err = sql.Open("postgres", dsn)
	default:
		db, err = sql.Open("mysql", dsn)
	}

	if err != nil {
		fmt.Printf("失败! 获取数据库sqlDB错误：%v\n", err)
		return
	}

	if db == nil {
		fmt.Printf("失败! 获取数据库sqlDB失败:%v\n", db)
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
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	fmt.Printf("需导出的表名称：%s\n", tbs.TbString)
	tbList := strings.Split(tbs.TbString, ",")

	cwdPath, _ := os.Getwd()
	cwdPath = strings.Replace(cwdPath, "\\", "/", -1)
	savePath := fmt.Sprintf("%s/tools/export_file", cwdPath)
	_, err = os.Stat(savePath)
	if os.IsNotExist(err) {
		fmt.Printf("失败! 数据文件夹不存在,请检查数据：%s\n", savePath)
		return
	}

	for _, oneTb := range tbList {
		fmt.Printf("初始化表：%s\n", oneTb)
		filePath := fmt.Sprintf("%s/%s.csv", savePath, oneTb)
		fmt.Printf("数据文件路径：%s\n", filePath)
		tbData, err := ReadCsv(filePath)
		if err != nil {
			fmt.Println("失败! 读取csv数据错误：", err)
		}

		if len(tbData) > 1 {
			cols := tbData[0]
			colString := strings.Join(cols, ",")
			rowsData := tbData[1:]

			var wString string
			for i := 0; i < len(cols); i++ {
				wString += "?"
				if i != len(cols)-1 {
					wString += ","
				}
			}

			SqlInsert := fmt.Sprintf("insert into %s (%s) VALUES(%s)", oneTb, colString, wString)
			stmt, err1 := db.Prepare(SqlInsert)
			if err1 != nil {
				fmt.Printf("失败! 预处理%s插入sql失败:%v\n", oneTb, err1)
				continue
			}

			var a interface{}
			for _, oneRow := range rowsData {
				sit := make([]interface{}, len(oneRow))
				for idx, it := range oneRow {
					if it != "" {
						sit[idx] = it
					} else {
						sit[idx] = a
					}

				}
				_, err2 := stmt.Exec(sit...)
				if err2 != nil {
					fmt.Printf("失败! 写入%s数据%s失败:%v\n", oneTb, oneRow, err2)
				}
			}
		}

	}

}

func ReadCsv(filepath string) ([][]string, error) {
	opencast, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("csv文件打开失败:%v", err)
	}
	defer func(opencast *os.File) {
		err := opencast.Close()
		if err != nil {

		}
	}(opencast)

	ReadCsv := csv.NewReader(opencast)


	ReadAll, err := ReadCsv.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv文件读取所有行失败:%v", err)
	}
	return ReadAll, nil
	/*
	  说明：
	   1、读取csv文件返回的内容为切片类型，可以通过遍历的方式使用或Slicer[0]方式获取具体的值。
	   2、同一个函数或线程内，两次调用Read()方法时，第二次调用时得到的值为每二行数据，依此类推。
	   3、大文件时使用逐行读取，小文件直接读取所有然后遍历，两者应用场景不一样，需要注意。
	*/

}
