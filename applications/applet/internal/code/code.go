package code

import "github.com/sybapp/infoflow/pkg/xcode"

var (
	RegisterPhoneEmpty     = xcode.New(10001, "注册手机号不能为空")
	RegisterUsernameEmpty  = xcode.New(10002, "用户名不能为空")
	VerifiyCodeEmpty       = xcode.New(10003, "验证码不能为空")
	PhoneHasRegistered     = xcode.New(10004, "手机号已经注册")
	UsernameHasRegistered  = xcode.New(10005, "用户名已经注册")
	RegisterPhoneNotNumber = xcode.New(10006, "手机号必须是数字")
	LoginError             = xcode.New(10007, "登录失败")
)
