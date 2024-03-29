package enums

import "errors"

//公共错误
var (
	DecodeErr						= errors.New("数据解析失败")
	SystemErr						= errors.New("系统繁忙，请稍后重试")
	NetErr							= errors.New("网络繁忙")
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
	TokenNotValid					= errors.New("请重新登录")  					//生成token出错
	TokenExpired					= errors.New("自动登录中...")  				//token已过期
	TokenNull						= errors.New("请先登录")  						//token不能为空
	UserIdTransErr					= errors.New("系统异常")  						//userId转换异常
	UserNotFound					= errors.New("用户不存在")  					//userId转换异常
	UserUpdatePhoneErr				= errors.New("更新用户手机号失败，请重试")  		//userId转换异常
)

//活动
var (
	ActivityPushQueueErr			= errors.New("参加活动写入队列失败")
	ActivityQueueNotFound			= errors.New("处理的参加活动队列记录不存在")
	ActivityQueueANotFound			= errors.New("处理的参加活动队列的活动记录不存在")
	ActivityJoinLogUpdateFailErr	= errors.New("更新活动参与记录因为加入活动因为人数已满失败出错")
	ActivityUpdateJoinNumFailErr	= errors.New("更新活动参与人数出错")
	ActivityStatusErr				= errors.New("活动不是可参加状态")
	ActivityJoinLimit				= errors.New("您今天参与的活动次数已经达到限制了，明天再来吧")
	ActivityJoinLimitShare			= errors.New("分享给好友可获得参与次数")
)

//socket
var (
	SocketEncreyErr					= errors.New("签名异常")
)

//读取配置
var (
	ReadConfigErr					= errors.New("配置信息错误")
)

//地址
var (
	AddressStoreErr					= errors.New("地址保存失败")
	AddressNotFound					= errors.New("地址记录不存在")
	AddressUpdateFail				= errors.New("地址记录更新失败")
	AddressListQueryERr				= errors.New("地址资料查询失败")
	AddressPageQueryERr				= errors.New("地址数据查询失败")
	AddressDetailQueryERr			= errors.New("地址数据查询失败")
	AddressDeleteDbErr  			= errors.New("地址删除失败失败")
	AddressPhoneErr					= errors.New("手机号格式错误")
)

//消息盒子
var (
	InboxStoreErr					= errors.New("消息保存失败")
	InboxUpdateReadErr				= errors.New("更新阅读时间错误")
	InboxPageQueryErr				= errors.New("查询失败")
	InboxCountQueryErr				= errors.New("查询失败")
)

//奖品发放
var (
	GiftPhoneBillSendIOErr			= errors.New("话费发放接口返回数据读取失败")
	GiftPhoneBillParseErr			= errors.New("话费发放接口返回数据解析失败")
	GiftPhoneBillRequestErr			= errors.New("话费发放请求网络错误")
)

//微信相关
var (
	WxAccessTokenRequestErr			= errors.New("请求Access token失败")
)
