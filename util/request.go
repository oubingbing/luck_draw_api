package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {

}

type beforeRequestHandle func(req *http.Request)

type afterRequestHandle func(resp *http.Response)

/**
 * get请求
 */
func (h HttpClient) Get(url string ,beforeHandle beforeRequestHandle,afterHandle afterRequestHandle) error {
	client := &http.Client{}
	var req *http.Request

	urlArr := strings.Split(url,"?")
	if len(urlArr)  == 2 {
		url = urlArr[0] + "?" + getParseParam(urlArr[1])
	}
	req, _ = http.NewRequest("GET", url, nil)

	if beforeHandle != nil{
		beforeHandle(req)
	}

	resp, err := client.Do(req)

	if afterHandle != nil {
		afterHandle(resp)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return  err
}

/**
 * Post请求
 */
func (h HttpClient) Post(urlVal string,data string,beforeHandle beforeRequestHandle,afterHandle afterRequestHandle) error {
	method  := "POST"
	client := &http.Client{}

	req, createErr := http.NewRequest(method, urlVal, bytes.NewBuffer([]byte(data)))
	if createErr != nil {
		fmt.Printf("创建失败:%v\n",createErr)
	}

	req.Header.Set("Content-Type", "application/json")

	if beforeHandle != nil{
		beforeHandle(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if afterHandle != nil{
		afterHandle(resp)
	}

	defer resp.Body.Close()

	return err
}

func getParseParam(param string) string  {
	return url.PathEscape(param)
}

func Input(ctx *gin.Context,key string) (interface{},bool) {
	var data map[string]interface{}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, &data)
	value,ok := data[key]
	return value,ok
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetClientIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		return "127.0.0.1"
	}

	return ip
}
