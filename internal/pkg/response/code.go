package response

type ResCode int64

const (
	CodeSuccess ResCode = 20000

	CodeInvalidParam = 49999 + iota
	CodeServerBusy
	CodeInvalidRouterRequested

	CodeInvalidToken
	CodeNeedLogin
	CodeGenericError
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess: "Success",

	CodeInvalidParam:           "请求参数错误",
	CodeServerBusy:             "系统繁忙，请稍候再试",
	CodeInvalidRouterRequested: "请求无效路由",

	CodeInvalidToken: "无效的Token",
	CodeNeedLogin:    "未登录",
	CodeGenericError: "Error",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
