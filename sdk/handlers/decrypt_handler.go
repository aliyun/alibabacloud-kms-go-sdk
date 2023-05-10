package handlers

import (
	"encoding/base64"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	dkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/utils"
	"net/http"
)

type DecryptTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewDecryptTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *DecryptTransferHandler {
	return &DecryptTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *DecryptTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *DecryptTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *DecryptTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	decryptReq := request.(*kms20160120.DecryptRequest)
	ciphertext, err := base64.StdEncoding.DecodeString(tea.StringValue(decryptReq.CiphertextBlob))
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < utils.EktIdLength+utils.GcmIvLength {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "The specified parameter CiphertextBlob is not valid, ciphertext length too small",
		})
	}
	var aad []byte
	if decryptReq.EncryptionContext != nil {
		aad, err = EncodeUserEncryptionContext(decryptReq.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
	ektId := ciphertext[:utils.EktIdLength]
	iv := ciphertext[utils.EktIdLength : utils.EktIdLength+utils.GcmIvLength]
	ciphertextBlob := ciphertext[utils.EktIdLength+utils.GcmIvLength:]
	result := &dkmssdk.DecryptRequest{
		Headers:        make(map[string]*string),
		CiphertextBlob: ciphertextBlob,
		Iv:             iv,
		Aad:            aad,
	}
	result.Headers[utils.MigrationKeyVersionIdKey] = tea.String(string(ektId))
	return result, nil
}

func (handler *DecryptTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.DecryptResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, _ := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	var encoding *string
	if runtime != nil && runtime.Encoding != nil {
		encoding = runtime.Encoding
	} else if handler.Config != nil && handler.Config.Encoding != nil {
		encoding = handler.Config.Encoding
	}
	var plaintext string
	if encoding == nil {
		plaintext = string(dkmsResponse.Plaintext)
	} else {
		var err error
		plaintext, err = utils.DecoderBytesToString(dkmsResponse.Plaintext, tea.StringValue(encoding))
		if err != nil {
			return nil, err
		}
	}
	body := &kms20160120.DecryptResponseBody{
		KeyId:        dkmsResponse.KeyId,
		KeyVersionId: keyVersionId,
		Plaintext:    tea.String(plaintext),
		RequestId:    dkmsResponse.RequestId,
	}
	return &kms20160120.DecryptResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *DecryptTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	decryptRequest := request.(*dkmssdk.DecryptRequest)
	return handler.DedicateClient.DecryptWithOptions(decryptRequest, runtimeOptions)
}

func (handler *DecryptTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	decryptReq := request.(*kms20160120.DecryptRequest)
	return handler.ShareClient.DecryptWithOptions(decryptReq, runtime.RuntimeOptions)
}
