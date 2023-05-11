package handlers

import (
	"encoding/base64"
	"fmt"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	dkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/utils"
	"net/http"
)

type EncryptTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewEncryptTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *EncryptTransferHandler {
	return &EncryptTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *EncryptTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *EncryptTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *EncryptTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	encryptReq := request.(*kms20160120.EncryptRequest)
	var aad []byte
	if encryptReq.EncryptionContext != nil {
		var err error
		aad, err = EncodeUserEncryptionContext(encryptReq.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
	var encoding *string
	if runtime != nil && runtime.Encoding != nil {
		encoding = runtime.Encoding
	} else if handler.Config != nil && handler.Config.Encoding != nil {
		encoding = handler.Config.Encoding
	}
	var plaintext []byte
	if encoding == nil {
		plaintext = []byte(tea.StringValue(encryptReq.Plaintext))
	} else {
		var err error
		plaintext, err = utils.EncoderStringToBytes(tea.StringValue(encryptReq.Plaintext), tea.StringValue(encoding))
		if err != nil {
			return nil, err
		}
	}
	result := &dkmssdk.EncryptRequest{
		KeyId:     encryptReq.KeyId,
		Plaintext: plaintext,
		Aad:       aad,
	}
	return result, nil
}

func (handler *EncryptTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.EncryptResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, ok := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	if !ok {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": fmt.Sprintf("Can not found response headers parameter[%s]", utils.MigrationKeyVersionIdKey),
		})
	}
	mkvId := []byte(tea.StringValue(keyVersionId))

	var ciphertextBlob []byte
	ciphertextBlob = append(ciphertextBlob, mkvId...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.Iv...)
	ciphertextBlob = append(ciphertextBlob, dkmsResponse.CiphertextBlob...)

	body := &kms20160120.EncryptResponseBody{
		KeyId:          dkmsResponse.KeyId,
		KeyVersionId:   keyVersionId,
		CiphertextBlob: tea.String(base64.StdEncoding.EncodeToString(ciphertextBlob)),
		RequestId:      dkmsResponse.RequestId,
	}
	return &kms20160120.EncryptResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *EncryptTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	encryptRequest := request.(*dkmssdk.EncryptRequest)
	return handler.DedicateClient.EncryptWithOptions(encryptRequest, runtimeOptions)
}

func (handler *EncryptTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	encryptReq := request.(*kms20160120.EncryptRequest)
	return handler.ShareClient.EncryptWithOptions(encryptReq, runtime.RuntimeOptions)
}
