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
	ACTIVITY_DETAIL_PARAM_ERR		= 2006		//详情id不能为空
	ACTIVITY_DETAIL_QUERY_ERR		= 2007		//详情查询错误
	ACTIVITY_DETAIL_NOT_FOUND		= 2008		//详情不存在
	ACTIVITY_JOIN_PARAM_ERR			= 2009		//参团参数失败，id为空
	ACTIVITY_JOIN_LIMIT				= 2010		//活动参与人数达到限制啦
	ACTIVITY_JOIN_SAVE_LOG_FAIL		= 2011		//参加活动失败
	ACTIVITY_JOIN_REPEAT			= 2012		//您已参加该活动，不可重复参加
	ACTIVITY_JOIN_QUERY_ERR			= 2013		//查询参与日志出错
)

//礼品相关 3000 ~ 3999
const (
	GIFT_SAVE_ERR					= 3000 		//礼品保存失败
	GIFT_FIRST_ERR					= 3001 		//礼品查询出错
	GIFT_NOT_FOUND					= 3002 		//礼品不存在
	GIFT_GET_DETAIL_ERR				= 3003 		//礼品详情查询错误
)


