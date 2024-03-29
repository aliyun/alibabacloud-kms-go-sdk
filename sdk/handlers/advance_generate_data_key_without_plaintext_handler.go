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

type AdvanceGenerateDataKeyWithoutPlaintextTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAdvanceGenerateDataKeyWithoutPlaintextTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler {
	return &AdvanceGenerateDataKeyWithoutPlaintextTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
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
	result := &dkmssdk.AdvanceGenerateDataKeyRequest{
		KeyId:         generateDataKeyWithoutPlaintextReq.KeyId,
		NumberOfBytes: numberOfBytesInteger,
		Aad:           aad,
	}
	return result, nil
}

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.AdvanceGenerateDataKeyResponse)
	keyVersionId := dkmsResponse.KeyVersionId
	from := utils.MagicNumLength + utils.CipherVerAndPaddingModeLength + utils.AlgorithmLength
	ciphertextBlob := dkmsResponse.CiphertextBlob[from:len(dkmsResponse.CiphertextBlob)]

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

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	advanceGenerateDataKeyRequest := request.(*dkmssdk.AdvanceGenerateDataKeyRequest)
	dkmsResponse, err := handler.DedicateClient.AdvanceGenerateDataKeyWithOptions(advanceGenerateDataKeyRequest, runtimeOptions)
	if err != nil {
		return nil, err
	}
	dkmsAdvanceEncryptRequest := &dkmssdk.AdvanceEncryptRequest{
		KeyId:     advanceGenerateDataKeyRequest.KeyId,
		Plaintext: []byte(base64.StdEncoding.EncodeToString(dkmsResponse.Plaintext)),
		Aad:       advanceGenerateDataKeyRequest.Aad,
	}
	advanceEncryptResponse, err := handler.DedicateClient.AdvanceEncryptWithOptions(dkmsAdvanceEncryptRequest, runtimeOptions)
	if err != nil {
		return nil, err
	}
	dkmsResponse.Iv = advanceEncryptResponse.Iv
	dkmsResponse.CiphertextBlob = advanceEncryptResponse.CiphertextBlob
	dkmsResponse.Headers = advanceEncryptResponse.Headers
	return dkmsResponse, nil
}

func (handler *AdvanceGenerateDataKeyWithoutPlaintextTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	generateDataKeyWithoutPlaintextReq := request.(*kms20160120.GenerateDataKeyWithoutPlaintextRequest)
	return handler.ShareClient.GenerateDataKeyWithoutPlaintextWithOptions(generateDataKeyWithoutPlaintextReq, runtime.RuntimeOptions)
}
