package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
)

const (
	dbConnection = "username:password@tcp(host:port)/dbname"
	checkPeriod  = 5 * time.Minute
)

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", dbConnection)
	if err != nil {
		fmt.Println("打开数据库失败,err:", err)
		return
	}
	fmt.Println("打开数据库成功")
	//初始化数据库
	database.SetConnMaxLifetime(60 * time.Second)
	database.SetMaxOpenConns(2)
	database.SetMaxIdleConns(2)
	Db = database
}

func monitorPhoneNumbers() {
	query := `SELECT respbody, hashvalue FROM xxx WHERE timestamp > ? AND respbody REGEXP '\\b1[3456789][0-9]{9}\\b'`
	monitor(query, "phone")

}

func monitorIDNumbers() {
	query := `SELECT respbody, hashvalue FROM xxx WHERE timestamp > ? AND respbody REGEXP '\\b[1-9]\\d{5}((19[2-9]\\d)|20[0-1]\\d)((0[1-9])|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X|x)\\b'`
	monitor(query, "ID")
}

func monitor(query, monitorType string) {
	lastChecked := time.Now().Add(-checkPeriod)
	fmt.Println(lastChecked.Format("2006-01-02 15:04:05"))
	rows, err := Db.Query(query, lastChecked.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("Error during '%s' number database query: %v\n", monitorType, err)
		return
	}
	defer rows.Close()

	foundNew := false
	var respBody, hashvalue string
	for rows.Next() {
		err := rows.Scan(&respBody, &hashvalue)
		if err != nil {
			log.Printf("Error scanning rows for '%s' number: %v\n", monitorType, err)
			continue
		}
		foundNew = true

	}
	if foundNew {
		sendAlert(hashvalue, monitorType)
	}
}

func sendAlert(hashvalue, monitorType string) {
	log.Printf("Alert for '%s', hashvalue: %s\n", monitorType, hashvalue)
	Robots(true, hashvalue, monitorType)
	robotsAtAll(true)
}

func main() {
	c := cron.New()
	c.AddFunc("*/5 * * * *", func() { //每五分钟执行一次
		go monitorPhoneNumbers()
		go monitorIDNumbers()
	})

	c.Start()
	select {}
}
