package main

import "swa-tools/tools"

var ExactType string

func main() {
	err := tools.InitSwaTools()
	if err != nil {
		panic("初始化swa-tools的配置数据出错")
	}
	ExactType = "export"

	if ExactType == "export" {
		tools.ExportSwaData()
	}
	if ExactType == "initSwa" {
		tools.InitSwaData()
	}
}
