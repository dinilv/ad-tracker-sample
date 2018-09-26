package main

import (
	db "github.com/adcamie/adserver/db/config"
	"github.com/adcamie/adserver/db/dao"
	model "github.com/adcamie/adserver/db/model/v1"
)

func main() {
	db.InitializeBigQuery()
	dao.CreateBigQueryTable("test", model.BigQueryTrackLogBackUp{})
}
