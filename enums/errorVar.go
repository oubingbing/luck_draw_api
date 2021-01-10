package enums

import "errors"

//登录授权
var (
	LoginTypeErr					= errors.New("登录类型错误")
	LoginRequestSessionErr			= errors.New("获取用户信息失败")
	LoginFail						= errors.New("登录失败")  						//请求微信数据异常
	LoginParseUserJsonErr			= errors.New("解析数据异常，请重试")  			//解析数据异常，请重试
	LoginQueryUserErr				= errors.New("用户查询错误")  					//用户查询错误
	LoginSaveUserDbErr				= errors.New("新增用户异常")  					//新增用户数据库异常
	LoginInsertUserErr				= errors.New("用户数据保存失败")  				//用户数据保存失败
	UpdateNicknameAvatarErr			= errors.New("用户数据更新失败")  				//用户数据更新失败
)

//读取配置
var (
	ReadConfigErr					= errors.New("配置信息错误")
)
