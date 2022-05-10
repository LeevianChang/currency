package model

const (
	CodeErrTime    = 100
	CodeErrDecrypt = 101
	CodeLimit      = 102
	CodeErrEncrypt = 103

	XTimeStamp = "X-Timestamp"
)

const (
	CodeErrTimeMessage    = "params error"
	CodeErrDecryptMessage = "decrypt message error"
	CodeErrEncryptMessage = "encrypt message error"
	CodeErrLimitMessage   = "limit message error"
)

// Reply ...
type Reply struct {
	Errno   int64       `json:"errno"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
