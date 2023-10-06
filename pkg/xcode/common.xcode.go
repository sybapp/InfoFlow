package xcode

var (
	OK                 = add(0, "OK")
	NoLogin            = add(101, "NOT_LOGIN")
	RequestErr         = add(400, "INVALID_ARGUMENT")
	Unauthorized       = add(401, "UNAUTHENTICATED")
	AccessDenied       = add(403, "PERMISSION_DENIED")
	NotFound           = add(404, "NOT_FOUND")
	MethodNotAllowed   = add(405, "METHOD_NOT_ALLOWED")
	Canceled           = add(498, "CANCELED")
	ServerErr          = add(500, "INTERNAL_ERROR")
	ServiceUnavailable = add(503, "UNAVAILABLE")
	Deadline           = add(504, "DEADLINE_EXCEEDED")
	LimitExceed        = add(509, "RESOURCE_EXHAUSTED")
)

var (
	DatabaseError         = 10001
	RedisError            = 10002
	ParameterError        = 10003
	RegisterError         = 10011
	UserHasRegistered     = 10012
	VerificationCodeError = 10013
	LoginError            = 10021
	LoginPasswordError    = 10022
	LoginUserNotExist     = 10023
)
