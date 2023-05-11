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

type AsymmetricSignTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAsymmetricSignTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AsymmetricSignTransferHandler {
	return &AsymmetricSignTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AsymmetricSignTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AsymmetricSignTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AsymmetricSignTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	signReq := request.(*kms20160120.AsymmetricSignRequest)
	message, err := base64.StdEncoding.DecodeString(tea.StringValue(signReq.Digest))
	if err != nil {
		return nil, err
	}
	result := &dkmssdk.SignRequest{
		Headers:     make(map[string]*string),
		KeyId:       signReq.KeyId,
		Message:     message,
		MessageType: tea.String(utils.DigestMessageType),
		Algorithm:   signReq.Algorithm,
	}
	if signReq.KeyVersionId != nil {
		result.Headers[utils.MigrationKeyVersionIdKey] = signReq.KeyVersionId
	}
	return result, nil
}

func (handler *AsymmetricSignTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.SignResponse)
	if dkmsResponse.Headers == nil {
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Can not found response headers",
		})
	}
	keyVersionId, _ := dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	body := &kms20160120.AsymmetricSignResponseBody{
		KeyId:        dkmsResponse.KeyId,
		KeyVersionId: keyVersionId,
		Value:        tea.String(base64.StdEncoding.EncodeToString(dkmsResponse.Signature)),
		RequestId:    dkmsResponse.RequestId,
	}
	return &kms20160120.AsymmetricSignResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *AsymmetricSignTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	signRequest := request.(*dkmssdk.SignRequest)
	return handler.DedicateClient.SignWithOptions(signRequest, runtimeOptions)
}

func (handler *AsymmetricSignTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricSignReq := request.(*kms20160120.AsymmetricSignRequest)
	return handler.ShareClient.AsymmetricSignWithOptions(asymmetricSignReq, runtime.RuntimeOptions)
}
