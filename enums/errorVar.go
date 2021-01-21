package enums

import "errors"

//公共错误
var (
	DecodeErr						= errors.New("数据解析失败")
)

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
	//JwtParseErr						= errors.New("解析数据失败")  						//
	UnKownSignMethod				= errors.New("授权异常")  						//Unexpected signing method
	LoginCreateTokenErr				= errors.New("授权错误")  						//生成token出错
	TokenNotValid					= errors.New("token非法")  						//生成token出错
	TokenExpired					= errors.New("token已过期")  						//token已过期
	TokenNull						= errors.New("token不能为空")  					//token不能为空
	UserIdTransErr					= errors.New("系统异常")  						//userId转换异常
	UserNotFound					= errors.New("用户不存在")  						//userId转换异常
)

//读取配置
var (
	ReadConfigErr					= errors.New("配置信息错误")
)
