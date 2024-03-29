package util

import (
	"bufio"
	"fmt"
	"io"
	"luck_draw/enums"
	"os"
	"strings"
)

func GetMysqlConfig() (string,*enums.ErrorInfo) {
	configs,errInfo := GetConfig()
	if errInfo != nil {
		return "",errInfo
	}

	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(configs["DB_USERNAME"]))
	builder.WriteString(":")
	builder.WriteString(strings.TrimSpace(configs["DB_PASSWORD"]))
	builder.WriteString("@(")
	builder.WriteString(strings.TrimSpace(configs["DB_HOST"]))
	builder.WriteString(":")
	builder.WriteString(strings.TrimSpace(configs["DB_PORT"]))
	builder.WriteString(")/")
	builder.WriteString(strings.TrimSpace(configs["DB_DATABASE"]))
	builder.WriteString("?charset=utf8")
	builder.WriteString("&parseTime=True&loc=Local")
	return builder.String(),nil
}

func GetAppConfig() map[string]string {
	dir, _ := os.Getwd()
	f,err := os.OpenFile(dir+"/app.conf",os.O_RDONLY,0777)
	if err != nil {
		Error(fmt.Sprintf("获取配置文件失败：%v\n",err.Error()))
	}

	reader := bufio.NewReader(f)
	configMp := make(map[string]string)
	for {
		line, err := reader.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		if strings.Index(line,"=") == -1 {
			continue
		}

		configKey := line[0:strings.Index(line,"=")]
		configValue := line[strings.Index(line,"=")+1:]
		configValue = strings.Replace(configValue,"\r","",-1)
		configValue = strings.Replace(configValue,"\n","",-1)
		configMp[configKey] = configValue
	}

	return configMp
}

func GetConfig() ( map[string]string,*enums.ErrorInfo) {
	configs := GetAppConfig()
	for key,_ := range configs {
		_,ok := configs[key]
		if !ok {
			return nil,&enums.ErrorInfo{Code:enums.READ_CONFIG_ERR,Err:enums.ReadConfigErr}
		}
	}

	return configs,nil
}
