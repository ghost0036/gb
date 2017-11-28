package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
)

type sendInfo struct {
	host     string
	port     string
	username string
	password string
}

var (
	h bool
	f string

	mailInfo sendInfo
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&f, "f", "con2fig.ini", "set configuration `file`")
	flag.Usage = usage

}

func main() {

	//初始化参数
	flag.Parse() //解析输入的参数

	cfg, err := goconfig.LoadConfigFile(f)
	if err != nil {
		log.Printf("读取配置文件失败[%s]", f)
		return
	}

	host, _ := cfg.GetValue("send", "host")
	port, _ := cfg.GetValue("send", "port")
	username, _ := cfg.GetValue("send", "username")
	password, _ := cfg.GetValue("send", "password")

	mailInfo = sendInfo{host, port, username, password}
	to := "shengpeng_lu@163.com"

	subject := "使用Golang发送邮件"

	body := `
		<html>
		<body>
		<h3>
		"Test send to email"
		</h3>
		</body>
		</html>
		`
	err = SendToMail(subject, body, to, "html")
	if err != nil {
		log.Printf("%s", err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `nginx version: nginx/1.10.0
Usage: server [-f filename]

Options:
`)
	flag.PrintDefaults()
}

/**
SendToMail：发送邮件方法
*/
func SendToMail(subject, body, to, mailtype string) error {
	auth := smtp.PlainAuth("", mailInfo.username, mailInfo.password, mailInfo.host)
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + mailInfo.username + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(mailInfo.host+":"+mailInfo.port, auth, mailInfo.username, sendTo, msg)
	return err
}
