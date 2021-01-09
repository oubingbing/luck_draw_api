package enums

//业务不相关错误
const (
	SUCCESS					= 0
	FAIL 					= 1
	DB_CONNECT_ERR 			= 2
)

//授权相关 1000 ~ 1999
//参数错误
const AUTH_PARAMS_ERROR =  1000

//活动相关 2000 ~ 2999
const (
	ACTIVITY_PARAM_ERR 				= 2000 		//参数错误
	ACTIVITY_START_DATE_ERR 		= 2001 		//活动开始如期解析错误
	ACTIVITY_END_DATE_ERR 			= 2002 		//活动截止日期解析错误
	ACTIVITY_RUN_DATE_ERR 			= 2003 		//活动开奖日期解析错误
	ACTIVITY_SAVE_ERR 				= 2004 		//活动保存失败
	ACTIVITY_PAGE_ERR				= 2005		//分页查询错误
)

//礼品相关 3000 ~ 3999
const (
	GIFT_SAVE_ERR					= 3000 		//礼品保存失败
	GIFT_FIRST_ERR					= 3001 		//礼品查询出错
	GIFT_NOT_FOUND					= 3002 		//礼品不存在
)


