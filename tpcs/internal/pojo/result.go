package pojo

const (
	ResultSuccess_False                      = false
	ResultSuccess_True                       = true
	ResultMsg_TryAgainLater                  = "请稍后再试！"
	ResultMsg_SendVcodeFail                  = "验证码发送失败！"
	ResultMsg_UserExisted                    = "该用户已存在！"
	ResultMsg_UserNotExists                  = "该用户不存在！"
	ResultMsg_UsernameNotNone                = "用户名不能为空！"
	ResultMsg_UsernameLengthLessThan15       = "用户名长度不能超过15个字符！"
	ResultMsg_UsernameOrPasswordErr          = "用户名或密码错误！"
	ResultMsg_UserPendingAudit               = "该用户尚未通过审核！"
	ResultMsg_EmailNotNone                   = "请填写邮箱！"
	ResultMsg_EmailExisted                   = "该邮箱已被其他用户绑定！"
	ResultMsg_EmailVcodeNotNone              = "请填写邮箱验证码！"
	ResultMsg_EmailVcodeIllegal              = "邮箱验证码错误或已过期！"
	ResultMsg_PasswordNotNone                = "密码不能为空！"
	ResultMsg_Password2NotNone               = "确认密码不能为空！"
	ResultMsg_2PasswordNotSame               = "两次输入的密码不一致！"
	ResultMsg_FormParseErr                   = "表单解析异常！"
	ResultMsg_FormErr_NoneEmail              = "表单异常，未获取到邮箱！"
	ResultMsg_FormErr_NoneUsername           = "表单异常，未获取到用户名！"
	ResultMsg_QuestionContentNotNone         = "请输入题目内容！"
	ResultMsg_QuestionScoreMoreThan0         = "题目分值必须大于0！"
	ResultMsg_QuestionTypeNotNone            = "请选择题型！"
	ResultMsg_QuestionDifficultyNotNone      = "请选择题目难度！"
	ResultMsg_QuestionCourseNotNone          = "请选择题目所属课程！"
	ResultMsg_InsufficientPermissions        = "权限不足！"
	ResultMsg_CourseExisted                  = "该课程已存在！"
	ResultMsg_QuestionCombinePlanNameExisted = "该组卷方案名称已存在！"
	ResultMsg_FileNotFound                   = "文件丢失或已损坏，下载失败！"
	ResultMsg_UserFreezed                    = "该用户已被冻结！"

	TPCS_Register_Audit = "TPCS Register Audit"
)

type Result struct {
	Code    interface{}/* int */ `json:"code"`
	Success interface{}/* bool */ `json:"success"`
	Msg     interface{}/* string */ `json:"msg"`
	Count   interface{}/* int */ `json:"count"`
	Data    interface{} `json:"data"`
}

type CombineResult struct {
	Success    bool   `json:"success"`
	PaperHtml  string `json:"paperHtml"`
	AnswerHtml string `json:"answerHtml"`
}
