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

type AsymmetricEncryptTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAsymmetricEncryptTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AsymmetricEncryptTransferHandler {
	return &AsymmetricEncryptTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AsymmetricEncryptTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AsymmetricEncryptTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AsymmetricEncryptTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricEncryptReq := request.(*kms20160120.AsymmetricEncryptRequest)
	plaintext, err := base64.StdEncoding.DecodeString(tea.StringValue(asymmetricEncryptReq.Plaintext))
	if err != nil {
		return nil, err
	}
	result := &dkmssdk.EncryptRequest{
		Headers:   make(map[string]*string),
		KeyId:     asymmetricEncryptReq.KeyId,
		Plaintext: plaintext,
		Algorithm: asymmetricEncryptReq.Algorithm,
	}
	if asymmetricEncryptReq.KeyVersionId != nil {
		result.Headers[utils.MigrationKeyVersionIdKey] = asymmetricEncryptReq.KeyVersionId
	}
	return result, nil
}

func (handler *AsymmetricEncryptTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.EncryptResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, _ := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	body := &kms20160120.AsymmetricEncryptResponseBody{
		KeyId:          dkmsResponse.KeyId,
		KeyVersionId:   keyVersionId,
		CiphertextBlob: tea.String(base64.StdEncoding.EncodeToString(dkmsResponse.CiphertextBlob)),
		RequestId:      dkmsResponse.RequestId,
	}
	return &kms20160120.AsymmetricEncryptResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *AsymmetricEncryptTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	encryptRequest := request.(*dkmssdk.EncryptRequest)
	return handler.DedicateClient.EncryptWithOptions(encryptRequest, runtimeOptions)
}

func (handler *AsymmetricEncryptTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricEncryptReq := request.(*kms20160120.AsymmetricEncryptRequest)
	return handler.ShareClient.AsymmetricEncryptWithOptions(asymmetricEncryptReq, runtime.RuntimeOptions)
}
