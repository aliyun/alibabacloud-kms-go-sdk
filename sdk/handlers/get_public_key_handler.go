package handlers

import (
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	dkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/utils"
	"net/http"
)

type GetPublicKeyTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewGetPublicKeyTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *GetPublicKeyTransferHandler {
	return &GetPublicKeyTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *GetPublicKeyTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *GetPublicKeyTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *GetPublicKeyTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	getPublicKeyReq := request.(*kms20160120.GetPublicKeyRequest)
	result := &dkmssdk.GetPublicKeyRequest{
		Headers: make(map[string]*string),
		KeyId:   getPublicKeyReq.KeyId,
	}
	if getPublicKeyReq.KeyVersionId != nil {
		result.Headers[utils.MigrationKeyVersionIdKey] = getPublicKeyReq.KeyVersionId
	}
	return result, nil
}

func (handler *GetPublicKeyTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.GetPublicKeyResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, _ := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	body := &kms20160120.GetPublicKeyResponseBody{
		KeyId:        dkmsResponse.KeyId,
		KeyVersionId: keyVersionId,
		PublicKey:    dkmsResponse.PublicKey,
		RequestId:    dkmsResponse.RequestId,
	}
	return &kms20160120.GetPublicKeyResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *GetPublicKeyTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	getPublicKeyRequest := request.(*dkmssdk.GetPublicKeyRequest)
	return handler.DedicateClient.GetPublicKeyWithOptions(getPublicKeyRequest, runtimeOptions)
}

func (handler *GetPublicKeyTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	getPublicKeyReq := request.(*kms20160120.GetPublicKeyRequest)
	return handler.ShareClient.GetPublicKeyWithOptions(getPublicKeyReq, runtime.RuntimeOptions)
}
