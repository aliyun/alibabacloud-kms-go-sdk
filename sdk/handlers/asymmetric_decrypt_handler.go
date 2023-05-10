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

type AsymmetricDecryptTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAsymmetricDecryptTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AsymmetricDecryptTransferHandler {
	return &AsymmetricDecryptTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AsymmetricDecryptTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AsymmetricDecryptTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AsymmetricDecryptTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricDecryptReq := request.(*kms20160120.AsymmetricDecryptRequest)
	ciphertextBlob, err := base64.StdEncoding.DecodeString(tea.StringValue(asymmetricDecryptReq.CiphertextBlob))
	if err != nil {
		return nil, err
	}
	result := &dkmssdk.DecryptRequest{
		Headers:        make(map[string]*string),
		KeyId:          asymmetricDecryptReq.KeyId,
		CiphertextBlob: ciphertextBlob,
		Algorithm:      asymmetricDecryptReq.Algorithm,
	}
	if asymmetricDecryptReq.KeyVersionId != nil {
		result.Headers[utils.MigrationKeyVersionIdKey] = asymmetricDecryptReq.KeyVersionId
	}
	return result, nil
}

func (handler *AsymmetricDecryptTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.DecryptResponse)
	var keyVersionId *string
	if dkmsResponse.Headers != nil {
		keyVersionId, _ = dkmsResponse.Headers[utils.MigrationKeyVersionIdKey]
	}
	body := &kms20160120.AsymmetricDecryptResponseBody{
		KeyId:        dkmsResponse.KeyId,
		KeyVersionId: keyVersionId,
		Plaintext:    tea.String(base64.StdEncoding.EncodeToString(dkmsResponse.Plaintext)),
		RequestId:    dkmsResponse.RequestId,
	}
	return &kms20160120.AsymmetricDecryptResponse{
		Body:       body,
		StatusCode: tea.Int32(http.StatusOK),
		Headers:    dkmsResponse.Headers,
	}, nil
}

func (handler *AsymmetricDecryptTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	runtimeOptions.Headers = append(runtimeOptions.Headers, tea.String(utils.MigrationKeyVersionIdKey))
	decryptRequest := request.(*dkmssdk.DecryptRequest)
	return handler.DedicateClient.DecryptWithOptions(decryptRequest, runtimeOptions)
}

func (handler *AsymmetricDecryptTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	asymmetricDecryptReq := request.(*kms20160120.AsymmetricDecryptRequest)
	return handler.ShareClient.AsymmetricDecryptWithOptions(asymmetricDecryptReq, runtime.RuntimeOptions)
}
