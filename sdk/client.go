// This file is auto-generated, don't edit it. Thanks.
package sdk

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
)

type Client struct {
	dedicatedkmssdk.Client
	KmsClient *kms20160120.Client
}

func NewClient(kmsInstanceConfig *dedicatedkmsopenapi.Config, openApiConfig *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(kmsInstanceConfig, openApiConfig)
	return client, err
}

func (client *Client) Init(kmsInstanceConfig *dedicatedkmsopenapi.Config, openApiConfig *openapi.Config) (_err error) {
	if kmsInstanceConfig == nil{
		kmsInstanceConfig = &dedicatedkmsopenapi.Config{Endpoint: tea.String("mock endpoint")}
	}

	_err = client.Client.Init(kmsInstanceConfig)
	if _err != nil {
		return _err
	}
	if openApiConfig == nil{
		openApiConfig = &openapi.Config{RegionId: tea.String("mock regionId")}
	}
	client.KmsClient, _err = kms20160120.NewClient(openApiConfig)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *Client) DoAction(query map[string]interface{}, action *string, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
	req := &openapi.OpenApiRequest{
		Query: openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      action,
		Version:     tea.String("2016-01-20"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		ReqBodyType: tea.String("formData"),
		BodyType:    tea.String("json"),
	}
	_result = make(map[string]interface{})
	_body, _err := client.KmsClient.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用CancelKeyDeletion接口撤销密钥删除
 * @param request
 * @return CancelKeyDeletionResponse
 */
func (client *Client) CancelKeyDeletion(request *kms20160120.CancelKeyDeletionRequest) (_result *kms20160120.CancelKeyDeletionResponse, _err error) {
	_result = &kms20160120.CancelKeyDeletionResponse{}
	_body, _err := client.KmsClient.CancelKeyDeletion(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用CancelKeyDeletion接口撤销密钥删除
 * @param request
 * @param runtime
 * @return CancelKeyDeletionResponse
 */
func (client *Client) CancelKeyDeletionWithOptions(request *kms20160120.CancelKeyDeletionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.CancelKeyDeletionResponse, _err error) {
	_result = &kms20160120.CancelKeyDeletionResponse{}
	_body, _err := client.KmsClient.CancelKeyDeletionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用CreateAlias接口为主密钥（CMK）创建一个别名
 * @param request
 * @return CreateAliasResponse
 */
func (client *Client) CreateAlias(request *kms20160120.CreateAliasRequest) (_result *kms20160120.CreateAliasResponse, _err error) {
	_result = &kms20160120.CreateAliasResponse{}
	_body, _err := client.KmsClient.CreateAlias(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用CreateAlias接口为主密钥（CMK）创建一个别名
 * @param request
 * @param runtime
 * @return CreateAliasResponse
 */
func (client *Client) CreateAliasWithOptions(request *kms20160120.CreateAliasRequest, runtime *util.RuntimeOptions) (_result *kms20160120.CreateAliasResponse, _err error) {
	_result = &kms20160120.CreateAliasResponse{}
	_body, _err := client.KmsClient.CreateAliasWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用CreateKey接口创建一个主密钥
 * @param request
 * @return CreateKeyResponse
 */
func (client *Client) CreateKey(request *kms20160120.CreateKeyRequest) (_result *kms20160120.CreateKeyResponse, _err error) {
	_result = &kms20160120.CreateKeyResponse{}
	_body, _err := client.KmsClient.CreateKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用CreateKey接口创建一个主密钥
 * @param request
 * @param runtime
 * @return CreateKeyResponse
 */
func (client *Client) CreateKeyWithOptions(request *kms20160120.CreateKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.CreateKeyResponse, _err error) {
	_result = &kms20160120.CreateKeyResponse{}
	_body, _err := client.KmsClient.CreateKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用CreateKeyVersion接口为用户主密钥（CMK）创建密钥版本
 * @param request
 * @return CreateKeyVersionResponse
 */
func (client *Client) CreateKeyVersion(request *kms20160120.CreateKeyVersionRequest) (_result *kms20160120.CreateKeyVersionResponse, _err error) {
	_result = &kms20160120.CreateKeyVersionResponse{}
	_body, _err := client.KmsClient.CreateKeyVersion(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用CreateKeyVersion接口为用户主密钥（CMK）创建密钥版本
 * @param request
 * @param runtime
 * @return CreateKeyVersionResponse
 */
func (client *Client) CreateKeyVersionWithOptions(request *kms20160120.CreateKeyVersionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.CreateKeyVersionResponse, _err error) {
	_result = &kms20160120.CreateKeyVersionResponse{}
	_body, _err := client.KmsClient.CreateKeyVersionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 创建凭据并存入凭据的初始版本
 * @param request
 * @return CreateSecretResponse
 */
func (client *Client) CreateSecret(request *kms20160120.CreateSecretRequest) (_result *kms20160120.CreateSecretResponse, _err error) {
	_result = &kms20160120.CreateSecretResponse{}
	_body, _err := client.KmsClient.CreateSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数创建凭据并存入凭据的初始版本
 * @param request
 * @param runtime
 * @return CreateSecretResponse
 */
func (client *Client) CreateSecretWithOptions(request *kms20160120.CreateSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.CreateSecretResponse, _err error) {
	_result = &kms20160120.CreateSecretResponse{}
	_body, _err := client.KmsClient.CreateSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DeleteAlias接口删除别名
 * @param request
 * @return DeleteAliasResponse
 */
func (client *Client) DeleteAlias(request *kms20160120.DeleteAliasRequest) (_result *kms20160120.DeleteAliasResponse, _err error) {
	_result = &kms20160120.DeleteAliasResponse{}
	_body, _err := client.KmsClient.DeleteAlias(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DeleteAlias接口删除别名
 * @param request
 * @param runtime
 * @return DeleteAliasResponse
 */
func (client *Client) DeleteAliasWithOptions(request *kms20160120.DeleteAliasRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DeleteAliasResponse, _err error) {
	_result = &kms20160120.DeleteAliasResponse{}
	_body, _err := client.KmsClient.DeleteAliasWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DeleteKeyMaterial接口删除已导入的密钥材料
 * @param request
 * @return DeleteKeyMaterialResponse
 */
func (client *Client) DeleteKeyMaterial(request *kms20160120.DeleteKeyMaterialRequest) (_result *kms20160120.DeleteKeyMaterialResponse, _err error) {
	_result = &kms20160120.DeleteKeyMaterialResponse{}
	_body, _err := client.KmsClient.DeleteKeyMaterial(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DeleteKeyMaterial接口删除已导入的密钥材料
 * @param request
 * @param runtime
 * @return DeleteKeyMaterialResponse
 */
func (client *Client) DeleteKeyMaterialWithOptions(request *kms20160120.DeleteKeyMaterialRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DeleteKeyMaterialResponse, _err error) {
	_result = &kms20160120.DeleteKeyMaterialResponse{}
	_body, _err := client.KmsClient.DeleteKeyMaterialWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DeleteSecret接口删除凭据对象
 * @param request
 * @return DeleteSecretResponse
 */
func (client *Client) DeleteSecret(request *kms20160120.DeleteSecretRequest) (_result *kms20160120.DeleteSecretResponse, _err error) {
	_result = &kms20160120.DeleteSecretResponse{}
	_body, _err := client.KmsClient.DeleteSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DeleteSecret接口删除凭据对象
 * @param request
 * @param runtime
 * @return DeleteSecretResponse
 */
func (client *Client) DeleteSecretWithOptions(request *kms20160120.DeleteSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DeleteSecretResponse, _err error) {
	_result = &kms20160120.DeleteSecretResponse{}
	_body, _err := client.KmsClient.DeleteSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DescribeKey接口查询用户主密钥（CMK）详情
 * @param request
 * @return DescribeKeyResponse
 */
func (client *Client) DescribeKey(request *kms20160120.DescribeKeyRequest) (_result *kms20160120.DescribeKeyResponse, _err error) {
	_result = &kms20160120.DescribeKeyResponse{}
	_body, _err := client.KmsClient.DescribeKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DescribeKey接口查询用户主密钥（CMK）详情
 * @param request
 * @param runtime
 * @return DescribeKeyResponse
 */
func (client *Client) DescribeKeyWithOptions(request *kms20160120.DescribeKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DescribeKeyResponse, _err error) {
	_result = &kms20160120.DescribeKeyResponse{}
	_body, _err := client.KmsClient.DescribeKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DescribeKeyVersion接口查询指定密钥版本信息
 * @param request
 * @return DescribeKeyVersionResponse
 */
func (client *Client) DescribeKeyVersion(request *kms20160120.DescribeKeyVersionRequest) (_result *kms20160120.DescribeKeyVersionResponse, _err error) {
	_result = &kms20160120.DescribeKeyVersionResponse{}
	_body, _err := client.KmsClient.DescribeKeyVersion(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DescribeKeyVersion接口查询指定密钥版本信息
 * @param request
 * @param runtime
 * @return DescribeKeyVersionResponse
 */
func (client *Client) DescribeKeyVersionWithOptions(request *kms20160120.DescribeKeyVersionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DescribeKeyVersionResponse, _err error) {
	_result = &kms20160120.DescribeKeyVersionResponse{}
	_body, _err := client.KmsClient.DescribeKeyVersionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DescribeSecret接口查询凭据的元数据信息
 * @param request
 * @return DescribeSecretResponse
 */
func (client *Client) DescribeSecret(request *kms20160120.DescribeSecretRequest) (_result *kms20160120.DescribeSecretResponse, _err error) {
	_result = &kms20160120.DescribeSecretResponse{}
	_body, _err := client.KmsClient.DescribeSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DescribeSecret接口查询凭据的元数据信息
 * @param request
 * @param runtime
 * @return DescribeSecretResponse
 */
func (client *Client) DescribeSecretWithOptions(request *kms20160120.DescribeSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DescribeSecretResponse, _err error) {
	_result = &kms20160120.DescribeSecretResponse{}
	_body, _err := client.KmsClient.DescribeSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DisableKey接口禁用指定的主密钥（CMK）进行加解密
 * @param request
 * @return DisableKeyResponse
 */
func (client *Client) DisableKey(request *kms20160120.DisableKeyRequest) (_result *kms20160120.DisableKeyResponse, _err error) {
	_result = &kms20160120.DisableKeyResponse{}
	_body, _err := client.KmsClient.DisableKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DisableKey接口禁用指定的主密钥（CMK）进行加解密
 * @param request
 * @param runtime
 * @return DisableKeyResponse
 */
func (client *Client) DisableKeyWithOptions(request *kms20160120.DisableKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.DisableKeyResponse, _err error) {
	_result = &kms20160120.DisableKeyResponse{}
	_body, _err := client.KmsClient.DisableKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用EnableKey接口启用指定的主密钥进行加解密
 * @param request
 * @return EnableKeyResponse
 */
func (client *Client) EnableKey(request *kms20160120.EnableKeyRequest) (_result *kms20160120.EnableKeyResponse, _err error) {
	_result = &kms20160120.EnableKeyResponse{}
	_body, _err := client.KmsClient.EnableKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用EnableKey接口启用指定的主密钥进行加解密
 * @param request
 * @param runtime
 * @return EnableKeyResponse
 */
func (client *Client) EnableKeyWithOptions(request *kms20160120.EnableKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.EnableKeyResponse, _err error) {
	_result = &kms20160120.EnableKeyResponse{}
	_body, _err := client.KmsClient.EnableKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ExportDataKey接口使用传入的公钥加密导出数据密钥
 * @param request
 * @return ExportDataKeyResponse
 */
func (client *Client) ExportDataKey(request *kms20160120.ExportDataKeyRequest) (_result *kms20160120.ExportDataKeyResponse, _err error) {
	_result = &kms20160120.ExportDataKeyResponse{}
	_body, _err := client.KmsClient.ExportDataKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ExportDataKey接口使用传入的公钥加密导出数据密钥
 * @param request
 * @param runtime
 * @return ExportDataKeyResponse
 */
func (client *Client) ExportDataKeyWithOptions(request *kms20160120.ExportDataKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ExportDataKeyResponse, _err error) {
	_result = &kms20160120.ExportDataKeyResponse{}
	_body, _err := client.KmsClient.ExportDataKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用GenerateAndExportDataKey接口随机生成一个数据密钥
 * @param request
 * @return GenerateAndExportDataKeyResponse
 */
func (client *Client) GenerateAndExportDataKey(request *kms20160120.GenerateAndExportDataKeyRequest) (_result *kms20160120.GenerateAndExportDataKeyResponse, _err error) {
	_result = &kms20160120.GenerateAndExportDataKeyResponse{}
	_body, _err := client.KmsClient.GenerateAndExportDataKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用GenerateAndExportDataKey接口随机生成一个数据密钥
 * @param request
 * @param runtime
 * @return GenerateAndExportDataKeyResponse
 */
func (client *Client) GenerateAndExportDataKeyWithOptions(request *kms20160120.GenerateAndExportDataKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.GenerateAndExportDataKeyResponse, _err error) {
	_result = &kms20160120.GenerateAndExportDataKeyResponse{}
	_body, _err := client.KmsClient.GenerateAndExportDataKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用GetParametersForImport接口获取导入主密钥材料的参数
 * @param request
 * @return GetParametersForImportResponse
 */
func (client *Client) GetParametersForImport(request *kms20160120.GetParametersForImportRequest) (_result *kms20160120.GetParametersForImportResponse, _err error) {
	_result = &kms20160120.GetParametersForImportResponse{}
	_body, _err := client.KmsClient.GetParametersForImport(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用GetParametersForImport接口获取导入主密钥材料的参数
 * @param request
 * @param runtime
 * @return GetParametersForImportResponse
 */
func (client *Client) GetParametersForImportWithOptions(request *kms20160120.GetParametersForImportRequest, runtime *util.RuntimeOptions) (_result *kms20160120.GetParametersForImportResponse, _err error) {
	_result = &kms20160120.GetParametersForImportResponse{}
	_body, _err := client.KmsClient.GetParametersForImportWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用GetRandomPassword接口获得一个随机口令字符串
 * @param request
 * @return GetRandomPasswordResponse
 */
func (client *Client) GetRandomPassword(request *kms20160120.GetRandomPasswordRequest) (_result *kms20160120.GetRandomPasswordResponse, _err error) {
	_result = &kms20160120.GetRandomPasswordResponse{}
	_body, _err := client.KmsClient.GetRandomPassword(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用GetRandomPassword接口获得一个随机口令字符串
 * @param request
 * @param runtime
 * @return GetRandomPasswordResponse
 */
func (client *Client) GetRandomPasswordWithOptions(request *kms20160120.GetRandomPasswordRequest, runtime *util.RuntimeOptions) (_result *kms20160120.GetRandomPasswordResponse, _err error) {
	_result = &kms20160120.GetRandomPasswordResponse{}
	_body, _err := client.KmsClient.GetRandomPasswordWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ImportKeyMaterial接口导入密钥材料
 * @param request
 * @return ImportKeyMaterialResponse
 */
func (client *Client) ImportKeyMaterial(request *kms20160120.ImportKeyMaterialRequest) (_result *kms20160120.ImportKeyMaterialResponse, _err error) {
	_result = &kms20160120.ImportKeyMaterialResponse{}
	_body, _err := client.KmsClient.ImportKeyMaterial(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ImportKeyMaterial接口导入密钥材料
 * @param request
 * @param runtime
 * @return ImportKeyMaterialResponse
 */
func (client *Client) ImportKeyMaterialWithOptions(request *kms20160120.ImportKeyMaterialRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ImportKeyMaterialResponse, _err error) {
	_result = &kms20160120.ImportKeyMaterialResponse{}
	_body, _err := client.KmsClient.ImportKeyMaterialWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListAliases接口查询当前用户在当前地域的所有别名
 * @param request
 * @return ListAliasesResponse
 */
func (client *Client) ListAliases(request *kms20160120.ListAliasesRequest) (_result *kms20160120.ListAliasesResponse, _err error) {
	_result = &kms20160120.ListAliasesResponse{}
	_body, _err := client.KmsClient.ListAliases(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListAliases接口查询当前用户在当前地域的所有别名
 * @param request
 * @param runtime
 * @return ListAliasesResponse
 */
func (client *Client) ListAliasesWithOptions(request *kms20160120.ListAliasesRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListAliasesResponse, _err error) {
	_result = &kms20160120.ListAliasesResponse{}
	_body, _err := client.KmsClient.ListAliasesWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListKeys查询调用者在调用地域的所有主密钥ID
 * @param request
 * @return ListKeysResponse
 */
func (client *Client) ListKeys(request *kms20160120.ListKeysRequest) (_result *kms20160120.ListKeysResponse, _err error) {
	_result = &kms20160120.ListKeysResponse{}
	_body, _err := client.KmsClient.ListKeys(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListKeys查询调用者在调用地域的所有主密钥ID
 * @param request
 * @param runtime
 * @return ListKeysResponse
 */
func (client *Client) ListKeysWithOptions(request *kms20160120.ListKeysRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListKeysResponse, _err error) {
	_result = &kms20160120.ListKeysResponse{}
	_body, _err := client.KmsClient.ListKeysWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListKeyVersions接口列出主密钥的所有密钥版本
 * @param request
 * @return ListKeyVersionsResponse
 */
func (client *Client) ListKeyVersions(request *kms20160120.ListKeyVersionsRequest) (_result *kms20160120.ListKeyVersionsResponse, _err error) {
	_result = &kms20160120.ListKeyVersionsResponse{}
	_body, _err := client.KmsClient.ListKeyVersions(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListKeyVersions接口列出主密钥的所有密钥版本
 * @param request
 * @param runtime
 * @return ListKeyVersionsResponse
 */
func (client *Client) ListKeyVersionsWithOptions(request *kms20160120.ListKeyVersionsRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListKeyVersionsResponse, _err error) {
	_result = &kms20160120.ListKeyVersionsResponse{}
	_body, _err := client.KmsClient.ListKeyVersionsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListResourceTags获取用户主密钥的标签
 * @param request
 * @return ListResourceTagsResponse
 */
func (client *Client) ListResourceTags(request *kms20160120.ListResourceTagsRequest) (_result *kms20160120.ListResourceTagsResponse, _err error) {
	_result = &kms20160120.ListResourceTagsResponse{}
	_body, _err := client.KmsClient.ListResourceTags(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListResourceTags获取用户主密钥的标签
 * @param request
 * @param runtime
 * @return ListResourceTagsResponse
 */
func (client *Client) ListResourceTagsWithOptions(request *kms20160120.ListResourceTagsRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListResourceTagsResponse, _err error) {
	_result = &kms20160120.ListResourceTagsResponse{}
	_body, _err := client.KmsClient.ListResourceTagsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListSecrets接口查询当前用户在当前地域创建的所有凭据
 * @param request
 * @return ListSecretsResponse
 */
func (client *Client) ListSecrets(request *kms20160120.ListSecretsRequest) (_result *kms20160120.ListSecretsResponse, _err error) {
	_result = &kms20160120.ListSecretsResponse{}
	_body, _err := client.KmsClient.ListSecrets(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListSecrets接口查询当前用户在当前地域创建的所有凭据
 * @param request
 * @param runtime
 * @return ListSecretsResponse
 */
func (client *Client) ListSecretsWithOptions(request *kms20160120.ListSecretsRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListSecretsResponse, _err error) {
	_result = &kms20160120.ListSecretsResponse{}
	_body, _err := client.KmsClient.ListSecretsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ListSecretVersionIds接口查询凭据的所有版本信息
 * @param request
 * @return ListSecretVersionIdsResponse
 */
func (client *Client) ListSecretVersionIds(request *kms20160120.ListSecretVersionIdsRequest) (_result *kms20160120.ListSecretVersionIdsResponse, _err error) {
	_result = &kms20160120.ListSecretVersionIdsResponse{}
	_body, _err := client.KmsClient.ListSecretVersionIds(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ListSecretVersionIds接口查询凭据的所有版本信息
 * @param request
 * @param runtime
 * @return ListSecretVersionIdsResponse
 */
func (client *Client) ListSecretVersionIdsWithOptions(request *kms20160120.ListSecretVersionIdsRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ListSecretVersionIdsResponse, _err error) {
	_result = &kms20160120.ListSecretVersionIdsResponse{}
	_body, _err := client.KmsClient.ListSecretVersionIdsWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用PutSecretValue接口为凭据存入一个新版本的凭据值
 * @param request
 * @return PutSecretValueResponse
 */
func (client *Client) PutSecretValue(request *kms20160120.PutSecretValueRequest) (_result *kms20160120.PutSecretValueResponse, _err error) {
	_result = &kms20160120.PutSecretValueResponse{}
	_body, _err := client.KmsClient.PutSecretValue(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用PutSecretValue接口为凭据存入一个新版本的凭据值
 * @param request
 * @param runtime
 * @return PutSecretValueResponse
 */
func (client *Client) PutSecretValueWithOptions(request *kms20160120.PutSecretValueRequest, runtime *util.RuntimeOptions) (_result *kms20160120.PutSecretValueResponse, _err error) {
	_result = &kms20160120.PutSecretValueResponse{}
	_body, _err := client.KmsClient.PutSecretValueWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用RestoreSecret接口恢复被删除的凭据
 * @param request
 * @return RestoreSecretResponse
 */
func (client *Client) RestoreSecret(request *kms20160120.RestoreSecretRequest) (_result *kms20160120.RestoreSecretResponse, _err error) {
	_result = &kms20160120.RestoreSecretResponse{}
	_body, _err := client.KmsClient.RestoreSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用RestoreSecret接口恢复被删除的凭据
 * @param request
 * @param runtime
 * @return RestoreSecretResponse
 */
func (client *Client) RestoreSecretWithOptions(request *kms20160120.RestoreSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.RestoreSecretResponse, _err error) {
	_result = &kms20160120.RestoreSecretResponse{}
	_body, _err := client.KmsClient.RestoreSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用RotateSecret接口手动轮转凭据
 * @param request
 * @return RotateSecretResponse
 */
func (client *Client) RotateSecret(request *kms20160120.RotateSecretRequest) (_result *kms20160120.RotateSecretResponse, _err error) {
	_result = &kms20160120.RotateSecretResponse{}
	_body, _err := client.KmsClient.RotateSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用RotateSecret接口手动轮转凭据
 * @param request
 * @param runtime
 * @return RotateSecretResponse
 */
func (client *Client) RotateSecretWithOptions(request *kms20160120.RotateSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.RotateSecretResponse, _err error) {
	_result = &kms20160120.RotateSecretResponse{}
	_body, _err := client.KmsClient.RotateSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用ScheduleKeyDeletion接口申请删除一个指定的主密钥（CMK)
 * @param request
 * @return ScheduleKeyDeletionResponse
 */
func (client *Client) ScheduleKeyDeletion(request *kms20160120.ScheduleKeyDeletionRequest) (_result *kms20160120.ScheduleKeyDeletionResponse, _err error) {
	_result = &kms20160120.ScheduleKeyDeletionResponse{}
	_body, _err := client.KmsClient.ScheduleKeyDeletion(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用ScheduleKeyDeletion接口申请删除一个指定的主密钥（CMK)
 * @param request
 * @param runtime
 * @return ScheduleKeyDeletionResponse
 */
func (client *Client) ScheduleKeyDeletionWithOptions(request *kms20160120.ScheduleKeyDeletionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.ScheduleKeyDeletionResponse, _err error) {
	_result = &kms20160120.ScheduleKeyDeletionResponse{}
	_body, _err := client.KmsClient.ScheduleKeyDeletionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用SetDeletionProtection接口为用户主密钥（CMK）开启或关闭删除保护
 * @param request
 * @return SetDeletionProtectionResponse
 */
func (client *Client) SetDeletionProtection(request *kms20160120.SetDeletionProtectionRequest) (_result *kms20160120.SetDeletionProtectionResponse, _err error) {
	_result = &kms20160120.SetDeletionProtectionResponse{}
	_body, _err := client.KmsClient.SetDeletionProtection(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用SetDeletionProtection接口为用户主密钥（CMK）开启或关闭删除保护
 * @param request
 * @param runtime
 * @return SetDeletionProtectionResponse
 */
func (client *Client) SetDeletionProtectionWithOptions(request *kms20160120.SetDeletionProtectionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.SetDeletionProtectionResponse, _err error) {
	_result = &kms20160120.SetDeletionProtectionResponse{}
	_body, _err := client.KmsClient.SetDeletionProtectionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用TagResource接口为主密钥、凭据或证书绑定标签
 * @param request
 * @return TagResourceResponse
 */
func (client *Client) TagResource(request *kms20160120.TagResourceRequest) (_result *kms20160120.TagResourceResponse, _err error) {
	_result = &kms20160120.TagResourceResponse{}
	_body, _err := client.KmsClient.TagResource(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用TagResource接口为主密钥、凭据或证书绑定标签
 * @param request
 * @param runtime
 * @return TagResourceResponse
 */
func (client *Client) TagResourceWithOptions(request *kms20160120.TagResourceRequest, runtime *util.RuntimeOptions) (_result *kms20160120.TagResourceResponse, _err error) {
	_result = &kms20160120.TagResourceResponse{}
	_body, _err := client.KmsClient.TagResourceWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UntagResource接口为主密钥、凭据或证书解绑标签
 * @param request
 * @return UntagResourceResponse
 */
func (client *Client) UntagResource(request *kms20160120.UntagResourceRequest) (_result *kms20160120.UntagResourceResponse, _err error) {
	_result = &kms20160120.UntagResourceResponse{}
	_body, _err := client.KmsClient.UntagResource(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UntagResource接口为主密钥、凭据或证书解绑标签
 * @param request
 * @param runtime
 * @return UntagResourceResponse
 */
func (client *Client) UntagResourceWithOptions(request *kms20160120.UntagResourceRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UntagResourceResponse, _err error) {
	_result = &kms20160120.UntagResourceResponse{}
	_body, _err := client.KmsClient.UntagResourceWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateAlias接口更新已存在的别名所代表的主密钥（CMK）ID
 * @param request
 * @return UpdateAliasResponse
 */
func (client *Client) UpdateAlias(request *kms20160120.UpdateAliasRequest) (_result *kms20160120.UpdateAliasResponse, _err error) {
	_result = &kms20160120.UpdateAliasResponse{}
	_body, _err := client.KmsClient.UpdateAlias(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateAlias接口更新已存在的别名所代表的主密钥（CMK）ID
 * @param request
 * @param runtime
 * @return UpdateAliasResponse
 */
func (client *Client) UpdateAliasWithOptions(request *kms20160120.UpdateAliasRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateAliasResponse, _err error) {
	_result = &kms20160120.UpdateAliasResponse{}
	_body, _err := client.KmsClient.UpdateAliasWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateKeyDescription接口更新主密钥的描述信息
 * @param request
 * @return UpdateKeyDescriptionResponse
 */
func (client *Client) UpdateKeyDescription(request *kms20160120.UpdateKeyDescriptionRequest) (_result *kms20160120.UpdateKeyDescriptionResponse, _err error) {
	_result = &kms20160120.UpdateKeyDescriptionResponse{}
	_body, _err := client.KmsClient.UpdateKeyDescription(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateKeyDescription接口更新主密钥的描述信息
 * @param request
 * @param runtime
 * @return UpdateKeyDescriptionResponse
 */
func (client *Client) UpdateKeyDescriptionWithOptions(request *kms20160120.UpdateKeyDescriptionRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateKeyDescriptionResponse, _err error) {
	_result = &kms20160120.UpdateKeyDescriptionResponse{}
	_body, _err := client.KmsClient.UpdateKeyDescriptionWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateRotationPolicy接口更新密钥轮转策略
 * @param request
 * @return UpdateRotationPolicyResponse
 */
func (client *Client) UpdateRotationPolicy(request *kms20160120.UpdateRotationPolicyRequest) (_result *kms20160120.UpdateRotationPolicyResponse, _err error) {
	_result = &kms20160120.UpdateRotationPolicyResponse{}
	_body, _err := client.KmsClient.UpdateRotationPolicy(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateRotationPolicy接口更新密钥轮转策略
 * @param request
 * @param runtime
 * @return UpdateRotationPolicyResponse
 */
func (client *Client) UpdateRotationPolicyWithOptions(request *kms20160120.UpdateRotationPolicyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateRotationPolicyResponse, _err error) {
	_result = &kms20160120.UpdateRotationPolicyResponse{}
	_body, _err := client.KmsClient.UpdateRotationPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateSecret接口更新凭据的元数据
 * @param request
 * @return UpdateSecretResponse
 */
func (client *Client) UpdateSecret(request *kms20160120.UpdateSecretRequest) (_result *kms20160120.UpdateSecretResponse, _err error) {
	_result = &kms20160120.UpdateSecretResponse{}
	_body, _err := client.KmsClient.UpdateSecret(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateSecret接口更新凭据的元数据
 * @param request
 * @param runtime
 * @return UpdateSecretResponse
 */
func (client *Client) UpdateSecretWithOptions(request *kms20160120.UpdateSecretRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateSecretResponse, _err error) {
	_result = &kms20160120.UpdateSecretResponse{}
	_body, _err := client.KmsClient.UpdateSecretWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateSecretRotationPolicy接口更新凭据轮转策略
 * @param request
 * @return UpdateSecretRotationPolicyResponse
 */
func (client *Client) UpdateSecretRotationPolicy(request *kms20160120.UpdateSecretRotationPolicyRequest) (_result *kms20160120.UpdateSecretRotationPolicyResponse, _err error) {
	_result = &kms20160120.UpdateSecretRotationPolicyResponse{}
	_body, _err := client.KmsClient.UpdateSecretRotationPolicy(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateSecretRotationPolicy接口更新凭据轮转策略
 * @param request
 * @param runtime
 * @return UpdateSecretRotationPolicyResponse
 */
func (client *Client) UpdateSecretRotationPolicyWithOptions(request *kms20160120.UpdateSecretRotationPolicyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateSecretRotationPolicyResponse, _err error) {
	_result = &kms20160120.UpdateSecretRotationPolicyResponse{}
	_body, _err := client.KmsClient.UpdateSecretRotationPolicyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用UpdateSecretVersionStage接口更新凭据的版本状态
 * @param request
 * @return UpdateSecretVersionStageResponse
 */
func (client *Client) UpdateSecretVersionStage(request *kms20160120.UpdateSecretVersionStageRequest) (_result *kms20160120.UpdateSecretVersionStageResponse, _err error) {
	_result = &kms20160120.UpdateSecretVersionStageResponse{}
	_body, _err := client.KmsClient.UpdateSecretVersionStage(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用UpdateSecretVersionStage接口更新凭据的版本状态
 * @param request
 * @param runtime
 * @return UpdateSecretVersionStageResponse
 */
func (client *Client) UpdateSecretVersionStageWithOptions(request *kms20160120.UpdateSecretVersionStageRequest, runtime *util.RuntimeOptions) (_result *kms20160120.UpdateSecretVersionStageResponse, _err error) {
	_result = &kms20160120.UpdateSecretVersionStageResponse{}
	_body, _err := client.KmsClient.UpdateSecretVersionStageWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用OpenKmsService接口为当前阿里云账号开通密钥管理服务
 * @return OpenKmsServiceResponse
 */
func (client *Client) OpenKmsService() (_result *kms20160120.OpenKmsServiceResponse, _err error) {
	_result = &kms20160120.OpenKmsServiceResponse{}
	_body, _err := client.KmsClient.OpenKmsService()
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用OpenKmsService接口为当前阿里云账号开通密钥管理服务
 * @param runtime
 * @return OpenKmsServiceResponse
 */
func (client *Client) OpenKmsServiceWithOptions(runtime *util.RuntimeOptions) (_result *kms20160120.OpenKmsServiceResponse, _err error) {
	_result = &kms20160120.OpenKmsServiceResponse{}
	_body, _err := client.KmsClient.OpenKmsServiceWithOptions(runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DescribeRegions接口查询当前账号的可用地域列表
 * @return DescribeRegionsResponse
 */
func (client *Client) DescribeRegions() (_result *kms20160120.DescribeRegionsResponse, _err error) {
	_result = &kms20160120.DescribeRegionsResponse{}
	_body, _err := client.KmsClient.DescribeRegions()
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DescribeRegions接口查询当前账号的可用地域列表
 * @param runtime
 * @return DescribeRegionsResponse
 */
func (client *Client) DescribeRegionsWithOptions(runtime *util.RuntimeOptions) (_result *kms20160120.DescribeRegionsResponse, _err error) {
	_result = &kms20160120.DescribeRegionsResponse{}
	_body, _err := client.KmsClient.DescribeRegionsWithOptions(runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用DescribeAccountKmsStatus接口查询当前阿里云账号的密钥管理服务状态
 * @return DescribeAccountKmsStatusResponse
 */
func (client *Client) DescribeAccountKmsStatus() (_result *kms20160120.DescribeAccountKmsStatusResponse, _err error) {
	_result = &kms20160120.DescribeAccountKmsStatusResponse{}
	_body, _err := client.KmsClient.DescribeAccountKmsStatus()
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用DescribeAccountKmsStatus接口查询当前阿里云账号的密钥管理服务状态
 * @param runtime
 * @return DescribeAccountKmsStatusResponse
 */
func (client *Client) DescribeAccountKmsStatusWithOptions(runtime *util.RuntimeOptions) (_result *kms20160120.DescribeAccountKmsStatusResponse, _err error) {
	_result = &kms20160120.DescribeAccountKmsStatusResponse{}
	_body, _err := client.KmsClient.DescribeAccountKmsStatusWithOptions(runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用GetSecretValue接口获取共享网关凭据值
 * @param request
 * @return GetSecretValueResponse
 */
func (client *Client) GetSecretValueBySharedEndpoint(request *kms20160120.GetSecretValueRequest) (_result *kms20160120.GetSecretValueResponse, _err error) {
	_result = &kms20160120.GetSecretValueResponse{}
	_body, _err := client.KmsClient.GetSecretValue(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用GetSecretValue接口获取共享网关凭据值
 * @param request
 * @param runtime
 * @return GetSecretValueResponse
 */
func (client *Client) GetSecretValueWithOptionsBySharedEndpoint(request *kms20160120.GetSecretValueRequest, runtime *util.RuntimeOptions) (_result *kms20160120.GetSecretValueResponse, _err error) {
	_result = &kms20160120.GetSecretValueResponse{}
	_body, _err := client.KmsClient.GetSecretValueWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 调用GetPublicKey接口获取共享网关非对称密钥的公钥
 * @param request
 * @return GetPublicKeyResponse
 */
func (client *Client) GetPublicKeyBySharedEndpoint(request *kms20160120.GetPublicKeyRequest) (_result *kms20160120.GetPublicKeyResponse, _err error) {
	_result = &kms20160120.GetPublicKeyResponse{}
	_body, _err := client.KmsClient.GetPublicKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * 带运行参数调用GetPublicKey接口获取共享网关非对称密钥的公钥
 * @param request
 * @param runtime
 * @return GetPublicKeyResponse
 */
func (client *Client) GetPublicKeyWithOptionsBySharedEndpoint(request *kms20160120.GetPublicKeyRequest, runtime *util.RuntimeOptions) (_result *kms20160120.GetPublicKeyResponse, _err error) {
	_result = &kms20160120.GetPublicKeyResponse{}
	_body, _err := client.KmsClient.GetPublicKeyWithOptions(request, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
