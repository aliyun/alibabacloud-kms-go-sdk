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

type AsymmetricVerifyTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAsymmetricVerifyTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AsymmetricVerifyTransferHandler {
	return &AsymmetricVerifyTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AsymmetricVerifyTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AsymmetricVerifyTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AsymmetricVerifyTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	verifyReq := request.(*kms20160120.AsymmetricVerifyRequest)
	message, err := base64.StdEncoding.DecodeString(tea.StringValue(verifyReq.Digest))
	if err != nil {
		return nil, err
	}
	signature, err := base64.StdEncoding.DecodeString(tea.StringValue(verifyReq.Value))
	if err != nil {
		return nil, err
	}
	result := &dkmssdk.VerifyRequest{
		Headers:     make(map[string]*string),
		KeyId:       verifyReq.KeyId,
		Signature:   signature,
		Message:     message,
		MessageType: tea.String(utils.DigestMessageType),
		Algorithm:   verifyReq.Algorithm,
	}
	if verifyReq.KeyVersionId != nil {
		result.Headers[utils.MigrationKeyVersionIdKey] = verifyReq.KeyVersionId
	}
	return result, nil
}

func (handler *AsymmetricVerifyTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.VerifyResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, _ := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	body := &kms20160120.AsymmetricVerifyResponseBody{
		KeyId:        dkmsResponse.KeyId,
		KeyVersionId: keyVersionId,
		Value:        dkmsResponse.Value,
		RequestId:    dkmsResponse.RequestId,
	}
	return &kms20160120.AsymmetricVerifyResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *AsymmetricVerifyTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	verifyRequest := request.(*dkmssdk.VerifyRequest)
	return handler.DedicateClient.VerifyWithOptions(verifyRequest, runtimeOptions)
}

func (handler *AsymmetricVerifyTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricVerifyReq := request.(*kms20160120.AsymmetricVerifyRequest)
	return handler.ShareClient.AsymmetricVerifyWithOptions(asymmetricVerifyReq, runtime.RuntimeOptions)
}
