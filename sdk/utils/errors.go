package utils

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	InvalidParamErrorCode                 = "InvalidParam"
	InvalidParameterErrorCode             = "InvalidParameter"
	UnauthorizedErrorCode                 = "Unauthorized"
	InvalidParamDateErrorMessage          = "The Param Date is invalid."
	InvalidParamAuthorizationErrorMessage = "The Param Authorization is invalid."
)

var errorCodeMap = map[string]string{
	"Forbidden.KeyNotFound":  "The specified Key is not found.",
	"Forbidden.NoPermission": "This operation is forbidden by permission system.",
	"InternalFailure":        "Internal Failure",
	"Rejected.Throttling":    "QPS Limit Exceeded",
}

type ErrorContent struct {
	HttpCode  int
	RequestId string
	HostId    string
	Code      string
	Message   string
}

func TransferTeaErrorServerError(err error) error {
	switch e := err.(type) {
	case *tea.SDKError:
		errCode := tea.StringValue(e.Code)
		errMessage := tea.StringValue(e.Message)
		switch errCode {
		case InvalidParamErrorCode:
			if errMessage == InvalidParamDateErrorMessage {
				e.Code = tea.String("IllegalTimestamp")
				e.Message = tea.String(`The input parameter "Timestamp" that is mandatory for processing this request is not supplied.`)
			} else if errMessage == InvalidParamAuthorizationErrorMessage {
				e.Code = tea.String("IncompleteSignature")
				e.Message = tea.String("The request signature does not conform to Aliyun standards.")
			}
		case UnauthorizedErrorCode:
			e.Code = tea.String("InvalidAccessKeyId.NotFound")
			e.Message = tea.String("The Access Key ID provided does not exist in our records.")
		default:
			msg, ok := errorCodeMap[errCode]
			if ok {
				e.Message = tea.String(msg)
			}
		}
		return e
	}
	return err
}

func NewMissingParameterError(paramName string) error {
	return tea.NewSDKError(map[string]interface{}{
		"code":    "ParameterMissing",
		"message": fmt.Sprintf("The parameter %s needed but no provided.", paramName),
	})
}
