package sdk

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/utils"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"

	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/handlers"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
)

type TransferClient struct {
	*kms20160120.Client
	handlers             map[string]handlers.KmsTransferHandler
	kmsClient            *dkmssdk.Client
	kmsConfig            *models.KmsConfig
	IsUseKmsShareGateway bool
}

func NewTransferClient(config *openapi.Config, kmsConfig interface{}) (*TransferClient, error) {
	client := new(TransferClient)
	err := client.Init(config, kmsConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewTransferClientWithInstanceGatewayConfig(config *dkmsopenapi.Config) (*TransferClient, error) {
	return NewTransferClient(nil, config)
}

func NewTransferClientWithShareGatewayConfig(config *openapi.Config) (*TransferClient, error) {
	return NewTransferClient(config, nil)
}

func (client *TransferClient) Init(config *openapi.Config, kmsConfig interface{}) error {
	client.handlers = make(map[string]handlers.KmsTransferHandler)
	if config != nil && kmsConfig == nil {
		shareClient, err := kms20160120.NewClient(config)
		if err != nil {
			return err
		}
		client.Client = shareClient
		client.IsUseKmsShareGateway = true
		client.initKmsTransferHandlers()
		return nil
	} else if kmsConfig != nil {
		c, err := models.TransferKmsConfig(kmsConfig)
		if err != nil {
			return err
		}
		client.kmsConfig = c
		if config == nil {
			config = &openapi.Config{Endpoint: c.Config.Endpoint}
		}
		shareClient, err := kms20160120.NewClient(config)
		if err != nil {
			return err
		}
		client.Client = shareClient
		kmsClient, err := dkmssdk.NewClient(c.Config)
		if err != nil {
			return err
		}
		client.kmsClient = kmsClient
		client.initKmsTransferHandlers()
		return nil
	}
	return tea.NewSDKError(map[string]interface{}{
		"message": "The parameter config and kmsConfig can not be both nil.",
	})
}

func (client *TransferClient) initKmsTransferHandlers() {
	client.handlers[utils.AsymmetricDecryptApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewAsymmetricDecryptTransferHandler(client.Client, client.kmsClient, utils.AsymmetricDecryptApiName, client.kmsConfig))
	client.handlers[utils.AsymmetricEncryptApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewAsymmetricEncryptTransferHandler(client.Client, client.kmsClient, utils.AsymmetricEncryptApiName, client.kmsConfig))
	client.handlers[utils.AsymmetricSignApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewAsymmetricSignTransferHandler(client.Client, client.kmsClient, utils.AsymmetricSignApiName, client.kmsConfig))
	client.handlers[utils.AsymmetricVerifyApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewAsymmetricVerifyTransferHandler(client.Client, client.kmsClient, utils.AsymmetricVerifyApiName, client.kmsConfig))
	if !client.kmsConfig.AdvanceSwitch {
		client.handlers[utils.DecryptApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewAdvanceDecryptTransferHandler(client.Client, client.kmsClient, utils.DecryptApiName, client.kmsConfig))
		client.handlers[utils.EncryptApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewAdvanceEncryptTransferHandler(client.Client, client.kmsClient, utils.EncryptApiName, client.kmsConfig))
		client.handlers[utils.GenerateDataKeyApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewAdvanceGenerateDataKeyTransferHandler(client.Client, client.kmsClient, utils.GenerateDataKeyApiName, client.kmsConfig))
		client.handlers[utils.GenerateDataKeyWithoutPlaintextApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewAdvanceGenerateDataKeyWithoutPlaintextTransferHandler(client.Client, client.kmsClient, utils.GenerateDataKeyWithoutPlaintextApiName, client.kmsConfig))
	} else {
		client.handlers[utils.DecryptApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewDecryptTransferHandler(client.Client, client.kmsClient, utils.DecryptApiName, client.kmsConfig))
		client.handlers[utils.EncryptApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewEncryptTransferHandler(client.Client, client.kmsClient, utils.EncryptApiName, client.kmsConfig))
		client.handlers[utils.GenerateDataKeyApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewGenerateDataKeyTransferHandler(client.Client, client.kmsClient, utils.GenerateDataKeyApiName, client.kmsConfig))
		client.handlers[utils.GenerateDataKeyWithoutPlaintextApiName] = handlers.NewWrappingKmsTransferHandler(
			handlers.NewGenerateDataKeyWithoutPlaintextTransferHandler(client.Client, client.kmsClient, utils.GenerateDataKeyWithoutPlaintextApiName, client.kmsConfig))
	}
	client.handlers[utils.GetPublicKeyApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewGetPublicKeyTransferHandler(client.Client, client.kmsClient, utils.GetPublicKeyApiName, client.kmsConfig))
	client.handlers[utils.GetSecretValueApiName] = handlers.NewWrappingKmsTransferHandler(
		handlers.NewGetSecretValueTransferHandler(client.Client, client.kmsClient, utils.GetSecretValueApiName, client.kmsConfig))
}

func (client *TransferClient) SetIsUseKmsShareGateway(isUseKmsShareGateway bool) {
	client.IsUseKmsShareGateway = isUseKmsShareGateway
}

func (client *TransferClient) dispatchApi(action string, request interface{}, runtime interface{}) (interface{}, error) {
	kmsRuntime, err := models.TransferKmsRuntimeOptions(runtime)
	if err != nil {
		return nil, err
	}
	handler, ok := client.handlers[action]
	if !ok || client.judgeCallDefaultKms(action, kmsRuntime) {
		return handler.CallKmsShareGateway(request, kmsRuntime)
	}
	result, err := handler.CallKmsDedicateGateway(request, kmsRuntime)
	if err != nil {
		return nil, utils.TransferTeaErrorServerError(err)
	}
	return result, nil
}

func (client *TransferClient) judgeCallDefaultKms(action string, runtime *models.KmsRuntimeOptions) bool {
	if runtime != nil && runtime.IsUseKmsShareGateway != nil {
		return tea.BoolValue(runtime.IsUseKmsShareGateway)
	}
	if client.kmsConfig != nil && client.kmsConfig.DefaultKmsApiNames != nil {
		for _, name := range client.kmsConfig.DefaultKmsApiNames {
			if action == name {
				return true
			}
		}
	}
	return client.IsUseKmsShareGateway
}

func (client *TransferClient) AsymmetricDecryptWithOptions(request *kms20160120.AsymmetricDecryptRequest, runtime interface{}) (_result *kms20160120.AsymmetricDecryptResponse, _err error) {
	response, _err := client.dispatchApi(utils.AsymmetricDecryptApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.AsymmetricDecryptResponse)
	return
}

func (client *TransferClient) AsymmetricDecrypt(request *kms20160120.AsymmetricDecryptRequest) (_result *kms20160120.AsymmetricDecryptResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.AsymmetricDecryptWithOptions(request, runtime)
}

func (client *TransferClient) AsymmetricEncryptWithOptions(request *kms20160120.AsymmetricEncryptRequest, runtime interface{}) (_result *kms20160120.AsymmetricEncryptResponse, _err error) {
	response, _err := client.dispatchApi(utils.AsymmetricEncryptApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.AsymmetricEncryptResponse)
	return
}

func (client *TransferClient) AsymmetricEncrypt(request *kms20160120.AsymmetricEncryptRequest) (_result *kms20160120.AsymmetricEncryptResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.AsymmetricEncryptWithOptions(request, runtime)
}

func (client *TransferClient) AsymmetricSignWithOptions(request *kms20160120.AsymmetricSignRequest, runtime interface{}) (_result *kms20160120.AsymmetricSignResponse, _err error) {
	response, _err := client.dispatchApi(utils.AsymmetricSignApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.AsymmetricSignResponse)
	return
}

func (client *TransferClient) AsymmetricSign(request *kms20160120.AsymmetricSignRequest) (_result *kms20160120.AsymmetricSignResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.AsymmetricSignWithOptions(request, runtime)
}

func (client *TransferClient) AsymmetricVerifyWithOptions(request *kms20160120.AsymmetricVerifyRequest, runtime interface{}) (_result *kms20160120.AsymmetricVerifyResponse, _err error) {
	response, _err := client.dispatchApi(utils.AsymmetricVerifyApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.AsymmetricVerifyResponse)
	return
}

func (client *TransferClient) AsymmetricVerify(request *kms20160120.AsymmetricVerifyRequest) (_result *kms20160120.AsymmetricVerifyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.AsymmetricVerifyWithOptions(request, runtime)
}

func (client *TransferClient) DecryptWithOptions(request *kms20160120.DecryptRequest, runtime interface{}) (_result *kms20160120.DecryptResponse, _err error) {
	response, _err := client.dispatchApi(utils.DecryptApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.DecryptResponse)
	return
}

func (client *TransferClient) Decrypt(request *kms20160120.DecryptRequest) (_result *kms20160120.DecryptResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.DecryptWithOptions(request, runtime)
}

func (client *TransferClient) EncryptWithOptions(request *kms20160120.EncryptRequest, runtime interface{}) (_result *kms20160120.EncryptResponse, _err error) {
	response, _err := client.dispatchApi(utils.EncryptApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.EncryptResponse)
	return
}

func (client *TransferClient) Encrypt(request *kms20160120.EncryptRequest) (_result *kms20160120.EncryptResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.EncryptWithOptions(request, runtime)
}

func (client *TransferClient) GenerateDataKeyWithOptions(request *kms20160120.GenerateDataKeyRequest, runtime interface{}) (_result *kms20160120.GenerateDataKeyResponse, _err error) {
	response, _err := client.dispatchApi(utils.GenerateDataKeyApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.GenerateDataKeyResponse)
	return
}

func (client *TransferClient) GenerateDataKey(request *kms20160120.GenerateDataKeyRequest) (_result *kms20160120.GenerateDataKeyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.GenerateDataKeyWithOptions(request, runtime)
}

func (client *TransferClient) GenerateDataKeyWithoutPlaintextWithOptions(request *kms20160120.GenerateDataKeyWithoutPlaintextRequest, runtime interface{}) (_result *kms20160120.GenerateDataKeyWithoutPlaintextResponse, _err error) {
	response, _err := client.dispatchApi(utils.GenerateDataKeyWithoutPlaintextApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.GenerateDataKeyWithoutPlaintextResponse)
	return
}

func (client *TransferClient) GenerateDataKeyWithoutPlaintext(request *kms20160120.GenerateDataKeyWithoutPlaintextRequest) (_result *kms20160120.GenerateDataKeyWithoutPlaintextResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.GenerateDataKeyWithoutPlaintextWithOptions(request, runtime)
}

func (client *TransferClient) GetPublicKeyWithOptions(request *kms20160120.GetPublicKeyRequest, runtime interface{}) (_result *kms20160120.GetPublicKeyResponse, _err error) {
	response, _err := client.dispatchApi(utils.GetPublicKeyApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.GetPublicKeyResponse)
	return
}

func (client *TransferClient) GetPublicKey(request *kms20160120.GetPublicKeyRequest) (_result *kms20160120.GetPublicKeyResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.GetPublicKeyWithOptions(request, runtime)
}

func (client *TransferClient) GetSecretValueWithOptions(request *kms20160120.GetSecretValueRequest, runtime interface{}) (_result *kms20160120.GetSecretValueResponse, _err error) {
	response, _err := client.dispatchApi(utils.GetSecretValueApiName, request, runtime)
	if _err != nil {
		return
	}
	_result = response.(*kms20160120.GetSecretValueResponse)
	return
}

func (client *TransferClient) GetSecretValue(request *kms20160120.GetSecretValueRequest) (_result *kms20160120.GetSecretValueResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	return client.GetSecretValueWithOptions(request, runtime)
}
