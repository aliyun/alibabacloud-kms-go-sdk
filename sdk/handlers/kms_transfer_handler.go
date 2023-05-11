package handlers

import (
	"bytes"
	"crypto/sha256"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
	"sort"
)

type KmsTransferHandler interface {
	CallKmsDedicateGateway(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
	CallKmsShareGateway(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
}

type TransferHandler interface {
	BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
	TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
	DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
	ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error)
	GetClient() interface{}
	GetAction() string
}

type WrappingKmsTransferHandler struct {
	TransferHandler
}

func NewWrappingKmsTransferHandler(handler TransferHandler) *WrappingKmsTransferHandler {
	return &WrappingKmsTransferHandler{
		TransferHandler: handler,
	}
}

func (transfer *WrappingKmsTransferHandler) CallKmsDedicateGateway(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	kmsRequest, err := transfer.BuildKmsRequest(request, runtime)
	if err != nil {
		return nil, err
	}
	response, err := transfer.DedicateGatewayApi(kmsRequest, runtime)
	if err != nil {
		return nil, err
	}
	return transfer.TransferResponse(response, runtime)
}

func (transfer *WrappingKmsTransferHandler) CallKmsShareGateway(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	return transfer.ShareGatewayApi(request, runtime)
}

func EncodeUserEncryptionContext(ctxInterfaceMap map[string]interface{}) ([]byte, error) {
	var keys []string
	ctxStrMap := map[string]string{}
	for key, value := range ctxInterfaceMap {
		strVal, ok := value.(string)
		if !ok {
			return nil, tea.NewSDKError(map[string]interface{}{
				"message": "The specified parameter EncryptContext is not valid",
			})
		}
		keys = append(keys, key)
		ctxStrMap[key] = strVal
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	keysLen := len(keys)
	for i, key := range keys {
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(ctxStrMap[key])
		if i < keysLen-1 {
			buf.WriteByte('&')
		}
	}
	encData := sha256.Sum256(buf.Bytes())
	return encData[:], nil
}
