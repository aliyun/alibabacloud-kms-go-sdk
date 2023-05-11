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

type GenerateDataKeyWithoutPlaintextTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewGenerateDataKeyWithoutPlaintextTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *GenerateDataKeyWithoutPlaintextTransferHandler {
	return &GenerateDataKeyWithoutPlaintextTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	generateDataKeyWithoutPlaintextReq := request.(*kms20160120.GenerateDataKeyWithoutPlaintextRequest)
	numberOfBytesInteger := generateDataKeyWithoutPlaintextReq.NumberOfBytes
	if generateDataKeyWithoutPlaintextReq.NumberOfBytes == nil {
		if tea.StringValue(generateDataKeyWithoutPlaintextReq.KeySpec) == "" {
			numberOfBytesInteger = tea.Int32(utils.NumberOfBytesAes256)
		} else if tea.StringValue(generateDataKeyWithoutPlaintextReq.KeySpec) == utils.KMSKeySpecAES256 {
			numberOfBytesInteger = tea.Int32(utils.NumberOfBytesAes256)
		} else if tea.StringValue(generateDataKeyWithoutPlaintextReq.KeySpec) == utils.KMSKeySpecAES128 {
			numberOfBytesInteger = tea.Int32(utils.NumberOfBytesAes128)
		} else {
			return nil, tea.NewSDKError(map[string]interface{}{
				"code":    utils.InvalidParameterErrorCode,
				"message": "The specified parameter KeySpec is not valid",
			})
		}
	}
	var aad []byte
	if generateDataKeyWithoutPlaintextReq.EncryptionContext != nil {
		var err error
		aad, err = EncodeUserEncryptionContext(generateDataKeyWithoutPlaintextReq.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
	result := &dkmssdk.GenerateDataKeyRequest{
		KeyId:         generateDataKeyWithoutPlaintextReq.KeyId,
		NumberOfBytes: numberOfBytesInteger,
		Aad:           aad,
	}
	return result, nil
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.GenerateDataKeyResponse)
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

	body := &kms20160120.GenerateDataKeyWithoutPlaintextResponseBody{
		KeyId:          dkmsResponse.KeyId,
		KeyVersionId:   keyVersionId,
		CiphertextBlob: tea.String(base64.StdEncoding.EncodeToString(ciphertextBlob)),
		RequestId:      dkmsResponse.RequestId,
	}
	return &kms20160120.GenerateDataKeyWithoutPlaintextResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	generateDataKeyRequest := request.(*dkmssdk.GenerateDataKeyRequest)
	dkmsResponse, err := handler.DedicateClient.GenerateDataKeyWithOptions(generateDataKeyRequest, runtimeOptions)
	if err != nil {
		return nil, err
	}
	dkmsEncryptRequest := &dkmssdk.EncryptRequest{
		KeyId:     generateDataKeyRequest.KeyId,
		Plaintext: []byte(base64.StdEncoding.EncodeToString(dkmsResponse.Plaintext)),
		Aad:       generateDataKeyRequest.Aad,
	}
	encryptResponse, err := handler.DedicateClient.EncryptWithOptions(dkmsEncryptRequest, runtimeOptions)
	if err != nil {
		return nil, err
	}
	dkmsResponse.Iv = encryptResponse.Iv
	dkmsResponse.CiphertextBlob = encryptResponse.CiphertextBlob
	dkmsResponse.Headers = encryptResponse.Headers
	return dkmsResponse, nil
}

func (handler *GenerateDataKeyWithoutPlaintextTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	generateDataKeyWithoutPlaintextReq := request.(*kms20160120.GenerateDataKeyWithoutPlaintextRequest)
	return handler.ShareClient.GenerateDataKeyWithoutPlaintextWithOptions(generateDataKeyWithoutPlaintextReq, runtime.RuntimeOptions)
}
