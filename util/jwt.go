package util

import (
	"fmt"
	"luck_draw/enums"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	Id  				uint
	OpenId 				string
	jwt.StandardClaims
}

/**
 * 解析token
 */
func ParseToken(token string) (map[string]interface{},error) {
	config ,configErr := GetConfig()
	if configErr != nil{
		Info(fmt.Sprintf("获取数据失败：%v\n",configErr.Err.Error()))
		return nil,configErr.Err
	}

	key := []byte(config["JWT_SECRET_KEY"])

	tokenPoint,err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			Error(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			return nil, enums.UnKownSignMethod
		}

		return key, nil
	})

	if tokenPoint != nil {
		if tokenPoint.Valid {
			//pass
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("That's not even a token")
				return nil,enums.TokenNotValid
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil,enums.TokenExpired
			} else {
				return nil,enums.TokenNotValid
			}
		} else {
			return nil,enums.TokenNotValid
		}
	}

	if err != nil {
		return nil,err
	}

	if c, ok := tokenPoint.Claims.(jwt.MapClaims); ok && tokenPoint.Valid {
		return c,nil
	} else {
		return  nil,err
	}
}

/**
 * 创建token
 */
func CreateToken(id uint,OpenId string) (string,error) {
	config ,configErr := GetConfig()
	if configErr != nil{
		Error(fmt.Sprintf("获取数据失败：%v\n",configErr.Err.Error()))
		return "",configErr.Err
	}

	key := []byte(config["JWT_SECRET_KEY"])

	claims := MyCustomClaims{
		id,
		OpenId,
		jwt.StandardClaims{
			ExpiresAt:time.Now().Unix()+(7200),//过期时间，两小时
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)//指定签名方法
	tokenString,err := token.SignedString(key)
	if err != nil{
		Error(fmt.Sprintf("授权错误：%v,key:%v",err.Error(),key))
		return  "",err
	}else{
		return tokenString,nil
	}
}