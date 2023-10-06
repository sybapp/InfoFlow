package code

import "github.com/sybapp/infoflow/pkg/xcode"

var (
	RegisterPhoneEmpty     = xcode.Errorf(10001, "注册手机号不能为空")
	RegisterUsernameEmpty  = xcode.Errorf(10002, "用户名不能为空")
	VerifiyCodeEmpty       = xcode.Errorf(10003, "验证码不能为空")
	PhoneHasRegistered     = xcode.Errorf(10004, "手机号已经注册")
	UsernameHasRegistered  = xcode.Errorf(10005, "用户名已经注册")
	RegisterPhoneNotNumber = xcode.Errorf(10006, "手机号必须是数字")
	LoginError             = xcode.Errorf(10007, "登录失败")
)
