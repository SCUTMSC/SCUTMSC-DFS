package util

type ErrorCode int32

const (
	_ int32 = iota + 9999
	StatusOK
	StatusParamInvalid
	StatusServerError
	StatusRegisterFailed
	StatusLoginFailed
	StatusInvalidToken
)
