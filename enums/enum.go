package enums

const ACTIVITY_QUEUE						= "luck_activity_queue"  	 			//参加活动队列
const INBOX_QUEUE							= "luck_inbox_queue"  	 				//保存消息盒子
const WX_NOTIFY_QUEUE						= "luck_wx_notify_queue"  	 			//微信通知队列
const ACTIVITY_HANDLE_PHONE_BILL_QUEUE		= "luck_handle_phone_bill_queue"  	 	//话费充值队列
const ACTIVITY_HANDLE_REA_PAK_QUEUE			= "luck_handle_red_pak_queue"  	 		//发送红包队列
const ACTIVITY_HANDLE_GIFT_QUEUE			= "luck_handle_gift_queue"  	 		//抽礼品队列
const ACTIVITY_QUEUE_TRY					= "luck_activity_queue_try"  			//参加活动重试队列
const SOCKET_USER_TOKEN  					= "luck_api_socket_token"	 			//用户socket token
const DATE_FORMAT  							= "2006-01-02 15:04:05"
const DATE_ONLY_FORMAT  					= "2006-01-02 00:00:00"
const DATE_DAY_FORMAT  						= "2006-01-02"
const DATE_FORMAT_STRING  					= "20060102150405"
const WX_ACCESS_TOKEN_CACHE_KEY  			= "luck_draw_wx_access_token"
const WX_TEMPLATE_SEND_SUCCESS  			= "GYJrbEJfKSFWIKcakFc03dm8F27IcBVoz8OUf2aawQI"
const WX_TEMPLATE_DRAW_FINISH  				= "HHOHnkh0UmYr-bifPvf1o0LWUHpBynwbxLbfPVMDQoA"

//业务不相关错误
const (
	SUCCESS					= 0
	FAIL 					= 1
	DB_CONNECT_ERR 			= 2
	READ_CONFIG_ERR			= 3
	DECODE_ARR_ERR			= 4
	NOT_FOUND				= 5
	SYSTEM_ERR				= 6					//系统异常
	WX_NOTIFY_IO_ERR		= 7					//模板消息读取数据流错误
)

//授权相关 1000 ~ 1999
//参数错误
const (
	AUTH_PARAMS_ERROR 				= 1000		//登录参数错误
	AUTH_LOGIN_TYPE_ERR				= 1001 		//登录类型错误
	AUTH_REQUEST_SESSION_ERR		= 1002		//请求微信session出错
	AUTH_REQUEST_SESSION_RESP_ERR	= 1003		//请求微信session返回错误异常
	AUTH_PARSE_JSON_ERR				= 1004		//解析用户json失败
	AUTH_USER_QUERY_ERR				= 1005		//用户查询错误
	AUTH_USER_SAVE_ERR				= 1006		//新增用户数据库异常
	AUTH_USER_UPDATE_ERR			= 1007		//用户数据更新失败
	AUTH_USER_PARSE_JWT_ERR			= 1008		//解析json失败
	AUTH_USER_CREATE_JWT_ERR		= 1009		//生成jwt异常
	AUTH_TOKEN_EXPIRED				= 1010		//token已过有效期
	AUTH_TOKEN_NULL					= 1011		//token为空
	AUTH_NOT_LOGIN					= 1012		//未登录
	Auth_TRANS_UID_ERR				= 1013		//userId类型转化失败
	AUTH_USER_NOT_FOUND				= 1014		//用户不存在
	AUTH_USER_UPDATE_PHONE_ERR		= 1015		//用户更新手机号失败
	AUTH_WX_ACCESSTOKEN_READ_IO_ERR = 1016		//获取Access token读取数据流错误
	AUTH_PARSE_ACCESS_TOKEN_ERR 	= 1017		//解析Access token失败
	AUTH_PARSE_ACCESS_REQUEST_ERR 	= 1018		//请求Access token失败
)

//活动相关 2000 ~ 2999
const (
	ACTIVITY_PARAM_ERR 					= 2000 		//参数错误
	ACTIVITY_START_DATE_ERR 			= 2001 		//活动开始如期解析错误
	ACTIVITY_END_DATE_ERR 				= 2002 		//活动截止日期解析错误
	ACTIVITY_RUN_DATE_ERR 				= 2003 		//活动开奖日期解析错误
	ACTIVITY_SAVE_ERR 					= 2004 		//活动保存失败
	ACTIVITY_PAGE_ERR					= 2005		//分页查询错误
	ACTIVITY_DETAIL_PARAM_ERR			= 2006		//详情id不能为空
	ACTIVITY_DETAIL_QUERY_ERR			= 2007		//详情查询错误
	ACTIVITY_DETAIL_NOT_FOUND			= 2008		//详情不存在
	ACTIVITY_JOIN_PARAM_ERR				= 2009		//参团参数失败，id为空
	ACTIVITY_JOIN_LIMIT					= 2010		//活动参与人数达到限制啦
	ACTIVITY_JOIN_SAVE_LOG_FAIL			= 2011		//参加活动失败
	ACTIVITY_JOIN_REPEAT				= 2012		//您已参加该活动，不可重复参加
	ACTIVITY_JOIN_QUERY_ERR				= 2013		//查询参与日志出错
	ACTIVITY_PUSH_QUEUE_ERR				= 2014		//参加活动写入队列失败
	ACTIVITY_DEAL_QUEUE_NOT_FOUND		= 2015		//处理参加活动队列的记录不存在
	ACTIVITY_DEAL_QUEUE_A_NOT_FOUND		= 2016		//处理参加活动队列的活动记录不存在
	ACTIVITY_DEAL_QUEUE_UPDATE_LOG_ERR	= 2017		//更新活动参与记录因为加入活动因为人数已满失败出错
	ACTIVITY_DEAL_QUEUE_UPDATE_A_ERR	= 2018		//更新活动参与人数出错
	ACTIVITY_STATUS_NOT_RUNNING	    	= 2019		//活动不是可参加状态
	ACTIVITY_MEMBER_ENOUTH		    	= 2020		//活动参加人数已满
	ACTIVITY_JOIN_LOG_QUERY_ERR		    = 2021		//活动参加记录查询错误
	ACTIVITY_JOIN_LOG_QUERY_MEMBER_ERR	= 2022		//活动参与会员查询错误
	ACTIVITY_Id_EMPYT			     	= 2023		//activity_id为空
	ACTIVITY_GET_RUNNING_ERR			= 2024		//完结活动时获取数据错误
	ACTIVITY_DEAL_NOT_HANDLE			= 2025		//未找到对应的活动处理函数
	ACTIVITY_PUSH_BILL_QUEUE_ERR		= 2026		//推送到话费发货队列失败
	ACTIVITY_FINDISH_DB_ERR				= 2027		//活动变更为已完成数据库出错
	ACTIVITY_UPDATE_CONSUME_DB_ERR		= 2028		//活动更新实际消耗奖品数量出错
	ACTIVITY_PUSH_GIFT_QUEUE_ERR		= 2029		//推送到物品发货队列失败
	ACTIVITY_UPDATE_JL_ERR				= 2030		//跟新join log错误
	ACTIVITY_JOIN_LIMIT_TIME			= 2031		//您今天参与的活动次数已经达到五次了，明天再来吧
)

//礼品相关 3000 ~ 3999
const (
	GIFT_SAVE_ERR					= 3000 		//礼品保存失败
	GIFT_FIRST_ERR					= 3001 		//礼品查询出错
	GIFT_NOT_FOUND					= 3002 		//礼品不存在
	GIFT_GET_DETAIL_ERR				= 3003 		//礼品详情查询错误
)

//socket相关 4000 ~ 4999
const (
	SOCKET_ENCRYPE_ERR				= 4000		//生成签名出错
	SOCKET_SIGN_ENCODE_ERR			= 4001		//sign转成json异常
	SOCKET_POST_SIGN_ERR			= 4002		//请求授权网络出错
	SOCKET_AUTH_ERR					= 4003		//请求授权返回失败
)

//地址相关 5000 ~ 5999
const (
	ADDRESS_USER_ID_ASSERT_ERR		= 5000		//保存地址中user id断言失败
	ADDRESS_STORE_ERR	 			= 5001		//保存地址数据库异常
	ADDRESS_STORE_AFFECT_ERR	 	= 5002		//保存地址数据库异常
	ADDRESS_NOT_FOUND		 		= 5003		//地址记录不存在
	ADDRESS_FIND_ERR		 		= 5004		//地址记录查询错误
	ADDRESS_UPDATE_ERR		 		= 5005		//地址记录更新失败
	ADDRESS_LIST_QUERY_ERR		 	= 5006		//地址资料查询失败
	ADDRESS_PAGE_QUERY_ERR		 	= 5007		//地址分页查询错误
	ADDRESS_DETAIL_QUERY_ERR		= 5008		//地址详情查询错误
	ADDRESS_DELETE_DB_ERR			= 5009		//地址删除错误
	ADDRESS_FORMAT_ERR				= 5010		//手机号格式错误
)

//消息盒子相关 6000 ~ 6999
const (
	INBOX_CREATE_FAIL				= 6000		//消息保存失败
	INBOX_UPDATE_READ_FAIL			= 6001		//更新阅读时间错误
	INBOX_PAGE_QUERY_FAIL			= 6002		//消息列表分页查询失败
	INBOX_COUNT_QUERY_FAIL			= 6003		//消息盒子统计未读出错
)

//礼品发放 7000 ~ 7999
const (
	GIFT_PHONE_BILL_SEND_IO_ERR		= 7000		//话费发放接口返回数据读取失败
	GIFT_PHONE_BILL_PARSE_ERR		= 7001		//话费发放接口返回数据解析失败
	GIFT_PHONE_BILL_SEND_FAIL		= 7002		//话费发放失败
	GIFT_PHONE_BILL_REQUEST_ERR		= 7003		//话费发放请求网络错误
)

