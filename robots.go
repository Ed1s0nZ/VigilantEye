package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func Robots(mentionAll bool, hashvale string, monitorType string) {
	mentionAllStr := "@xxx"
	a := ""
	switch monitorType {
	case "phone":
		a = "手机号"
	case "ID":
		a = "身份证号"
	// 可以增加更多的case
	default:
		a = "未知类型"
	}

	if mentionAll {
		mentionAllStr = "@all"
	}
	dataJsonStr := fmt.Sprintf(`{
		"msgtype": "markdown", "markdown": {
			"content": "<font color=\"red\">[测试流量敏感信息监控告警] </font>\n发现hashvalue为 ` + hashvale + ` 响应包中包含明文<font color=\"warning\">` + a + `</font>, 请登录数据库查看。` + `\n时间 : ` + timenow() + `", 
			"mentioned_list": ["` + mentionAllStr + `"]}}`)
	fmt.Println(dataJsonStr)
	resp, err := http.Post(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
		"application/json",
		bytes.NewBuffer([]byte(dataJsonStr)))
	if err != nil {
		// fmt.Println("weworkAlarm request error")
		return
	}
	defer resp.Body.Close()
	// webAlive()
}
func robotsAtAll(mentionAll bool) {
	mentionAllStr := "@xxx"

	if mentionAll {
		mentionAllStr = "@all"
	}
	dataJsonStr := fmt.Sprintf(`{
		"msgtype": "text", "text": {
			"mentioned_list": ["` + mentionAllStr + `"]}}`)
	fmt.Println(dataJsonStr)
	resp, err := http.Post(
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
		"application/json",
		bytes.NewBuffer([]byte(dataJsonStr)))
	if err != nil {
		// fmt.Println("weworkAlarm request error")
		return
	}
	defer resp.Body.Close()
}
func timenow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
