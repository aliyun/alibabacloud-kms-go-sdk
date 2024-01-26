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

type AdvanceDecryptTransferHandler struct {
	ShareClient    *kms20160120.Client
	DedicateClient *dkmssdk.Client
	Config         *models.KmsConfig
	Action         string
}

func NewAdvanceDecryptTransferHandler(shareClient *kms20160120.Client, dedicateClient *dkmssdk.Client, action string, config *models.KmsConfig) *AdvanceDecryptTransferHandler {
	return &AdvanceDecryptTransferHandler{
		ShareClient:    shareClient,
		DedicateClient: dedicateClient,
		Config:         config,
		Action:         action,
	}
}

func (handler *AdvanceDecryptTransferHandler) GetClient() interface{} {
	return handler.DedicateClient
}

func (handler *AdvanceDecryptTransferHandler) GetAction() string {
	return handler.Action
}

func (handler *AdvanceDecryptTransferHandler) BuildKmsRequest(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
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
	ivBytes := ciphertext[utils.EktIdLength : utils.EktIdLength+utils.GcmIvLength]
	cipherVerAndPaddingMode := utils.CipherVer<<4 | 0
	var aad []byte
	if decryptReq.EncryptionContext != nil {
		aad, err = EncodeUserEncryptionContext(decryptReq.EncryptionContext)
		if err != nil {
			return nil, err
		}
	}
	var arr = []byte{utils.MagicNum, byte(cipherVerAndPaddingMode), utils.AlgAesGcm}
	ciphertextBlob := append(arr, ciphertext...)
	result := &dkmssdk.AdvanceDecryptRequest{
		Headers:        make(map[string]*string),
		CiphertextBlob: ciphertextBlob,
		Iv:             ivBytes,
		Aad:            aad,
	}
	return result, nil
}

func (handler *AdvanceDecryptTransferHandler) TransferResponse(response interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	dkmsResponse := response.(*dkmssdk.AdvanceDecryptResponse)
	keyVersionId := dkmsResponse.KeyVersionId
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

func (handler *AdvanceDecryptTransferHandler) DedicateGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	runtimeOptions := runtime.ToDKmsRuntimeOptions()
	decryptRequest := request.(*dkmssdk.AdvanceDecryptRequest)
	return handler.DedicateClient.AdvanceDecryptWithOptions(decryptRequest, runtimeOptions)
}

func (handler *AdvanceDecryptTransferHandler) ShareGatewayApi(request interface{}, runtime *models.KmsRuntimeOptions) (interface{}, error) {
	decryptReq := request.(*kms20160120.DecryptRequest)
	return handler.ShareClient.DecryptWithOptions(decryptReq, runtime.RuntimeOptions)
}
