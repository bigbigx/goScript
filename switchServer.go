// switchServer
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {
		log.Println("Start .. Net dial server .. ")
		time.Sleep(2 * time.Second)
		dsNum := DialServer()
		log.Println("Start .. Ping server .. ")
		psNum := PingServer()
		if dsNum == 0 && psNum == 0 {
			log.Println("server is down, Begin switch .. ")
			body := "<h3>OA Hardware load balancing is down !<h3><h4>Solution :</h4>1: Server 10.0.1.159 is standby application,modification ip to 10.0.0.43 and restart network.<br>2: repair 10.0.0.43 system server and restart .<br><font size=\"4\" color=\"red\">Please manual handle troubleshooting, right now!!</font>"
			SwitchServer("OA Faild ", "html", body)
		}
		time.Sleep(10 * time.Second)
	}
}

// value of contentType is  "plain" or "html"
func SwitchServer(subject, contentType, body string) {
	//需要自定义发件人 邮箱 密码 邮件服务器
	auth := smtp.PlainAuth("", "username@xx.com", "password", "smtp.xxx.net")
	nickname := "DBA-zhifang.Tang"
	to := []string{"username@xx.com"}
	user := "username@xx.com"
	content_type := "Content-Type: text/" + contentType + "; charset=UTF-8"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.263.net:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}

}
//实现ping功能，通过调用系统命令ping实现。
func PingServer() int {
	cmd := "ping -c 1 10.0.0.43 | grep 'Host Unreachable' | wc -l"
	bo, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return 0
	}
	return BytesToInt(bo)
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
//net dial端口，通过net包Dial功能实现。
func DialServer() int {
	_, err := net.Dial("tcp", "10.0.0.43:80")
	if err != nil {
		return 0
	}
	return 1
}
