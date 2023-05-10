package handlers

import (
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	dkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/utils"
	"net/http"
)

type GetSecretValueTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewGetSecretValueTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *GetSecretValueTransferHandler {
	return &GetSecretValueTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *GetSecretValueTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *GetSecretValueTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *GetSecretValueTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	getSecretValueReq := request.(*kms20160120.GetSecretValueRequest)
	result := &dkmssdk.GetSecretValueRequest{
		SecretName:          getSecretValueReq.SecretName,
		VersionStage:        getSecretValueReq.VersionStage,
		VersionId:           getSecretValueReq.VersionId,
		FetchExtendedConfig: getSecretValueReq.FetchExtendedConfig,
	}
	return result, nil
}

func (handler *GetSecretValueTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.GetSecretValueResponse)
	body := &kms20160120.GetSecretValueResponseBody{
		SecretName:        dkmsResponse.SecretName,
		VersionId:         dkmsResponse.VersionId,
		CreateTime:        dkmsResponse.CreateTime,
		SecretData:        dkmsResponse.SecretData,
		SecretDataType:    dkmsResponse.SecretDataType,
		AutomaticRotation: dkmsResponse.AutomaticRotation,
		RotationInterval:  dkmsResponse.RotationInterval,
		NextRotationDate:  dkmsResponse.NextRotationDate,
		ExtendedConfig:    dkmsResponse.ExtendedConfig,
		LastRotationDate:  dkmsResponse.LastRotationDate,
		SecretType:        dkmsResponse.SecretType,
		VersionStages: &kms20160120.GetSecretValueResponseBodyVersionStages{
			VersionStage: []*string{},
		},
		RequestId: dkmsResponse.RequestId,
	}
	for _, state := range dkmsResponse.VersionStages {
		body.VersionStages.VersionStage = append(body.VersionStages.VersionStage, state)
	}
	return &kms20160120.GetSecretValueResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *GetSecretValueTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	getSecretValueRequest := request.(*dkmssdk.GetSecretValueRequest)
	return handler.DedicateClient.GetSecretValueWithOptions(getSecretValueRequest, runtimeOptions)
}

func (handler *GetSecretValueTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	getSecretValueReq := request.(*kms20160120.GetSecretValueRequest)
	return handler.ShareClient.GetSecretValueWithOptions(getSecretValueReq, runtime.RuntimeOptions)
}
