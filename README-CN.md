# 阿里云KMS Go SDK

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

阿里云KMS Go SDK可以帮助Golang开发者快速使用KMS。

- [阿里云专属KMS主页](https://help.aliyun.com/document_detail/311016.html)
- [代码示例](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-kms-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-kms-go-sdk/releases)

## 优势

帮助Golang开发者通过本SDK快速使用阿里云KMS产品的所有API:

- 支持通过KMS公共网关访问进行KMS资源管理和密钥运算
- 支持通过KMS实例网关进行密钥运算

## 软件要求

- Golang 1.13及以上。

## 安装

您可以使用`go mod`管理您的依赖：

```
require (
	github.com/aliyun/alibabacloud-kms-go-sdk v1.0.0
)
```

或者，通过go get命令获取远程代码包：

```
$ go get -u github.com/aliyun/alibabacloud-kms-go-sdk
```

## KMS Client介绍

| KMS 客户端结构体 | 简介 | 使用场景 |
| :-----| :---- | :---- |
| Client | 支持KMS资源管理和KMS实例网关的密钥运算| 1.仅通过VPC网关进行密钥运算操作的场景。<br>2.仅通过公共网关对KMS资源管理的场景。 <br>3.既要通过VPC网关进行密钥运算操作又要通过公共网关对KMS资源管理的场景。|
| TransferClient | 支持用户应用简单修改的情况下就可以从KMS 1.0密钥运算迁移到 KMS 3.0密钥运算 | 使用阿里云 SDK访问KMS 1.0密钥运算的用户，需要迁移到KMS 3.0的场景。|

## 示例代码

### 1. 仅通过VPC网关进行密钥运算操作的场景。

#### 参考以下示例代码调用KMS Encrypt API。更多API示例参考 [operation samples](./examples/operation)

```go
package example

import (
	"os"
	console "github.com/alibabacloud-go/tea-console/client"
	env "github.com/alibabacloud-go/darabonba-env/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	kmssdk "github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateKmsInstanceConfig(clientKeyFile *string, password *string, endpoint *string, caFilePath *string) (_result *dedicatedkmsopenapi.Config, _err error) {
	config := &dedicatedkmsopenapi.Config{
		ClientKeyFile: clientKeyFile,
		Password:      password,
		Endpoint:      endpoint,
		CaFilePath:    caFilePath,
	}
	_result = config
	return _result, _err
}

func CreateClient(kmsInstanceConfig *dedicatedkmsopenapi.Config) (_result *kmssdk.Client, _err error) {
	_result = &kmssdk.Client{}
	_result, _err = kmssdk.NewClient(kmsInstanceConfig, nil)
	return _result, _err
}

func Encrypt(client *kmssdk.Client, paddingMode *string, aad []byte, keyId *string, plaintext []byte, iv []byte, algorithm *string) (_result *dedicatedkmssdk.EncryptResponse, _err error) {
	request := &dedicatedkmssdk.EncryptRequest{
		PaddingMode: paddingMode,
		Aad:         aad,
		KeyId:       keyId,
		Plaintext:   plaintext,
		Iv:          iv,
		Algorithm:   algorithm,
	}
	_result = &dedicatedkmssdk.EncryptResponse{}
	return client.Encrypt(request)
}

func _main(args []*string) (_err error) {
	kmsInstanceConfig, _err := CreateKmsInstanceConfig(env.GetEnv(tea.String("your client key file path env")), env.GetEnv(tea.String("your client key password env")), tea.String("your kms instance endpoint env"), tea.String("your ca file path"))
	if _err != nil {
		return _err
	}

	client, _err := CreateClient(kmsInstanceConfig)
	if _err != nil {
		return _err
	}

	paddingMode := tea.String("your paddingMode")
	aad := util.ToBytes(tea.String("your aad"))
	keyId := tea.String("your keyId")
	plaintext := util.ToBytes(tea.String("your plaintext"))
	iv := util.ToBytes(tea.String("your iv"))
	algorithm := tea.String("your algorithm")
	response, _err := Encrypt(client, paddingMode, aad, keyId, plaintext, iv, algorithm)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(response))
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
```

### 2. 仅通过公共网关对KMS资源管理的场景。

#### 参考以下示例代码调用KMS CreateKey API。更多API示例参考 [manage samples](./examples/manage)

```go
package example

import (
	"os"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	env "github.com/alibabacloud-go/darabonba-env/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	kmssdk "github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateOpenApiConfig(accessKeyId *string, accessKeySecret *string, regionId *string) (_result *openapi.Config, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
	}
	_result = config
	return _result, _err
}

func CreateClient(openApiConfig *openapi.Config) (_result *kmssdk.Client, _err error) {
	_result = &kmssdk.Client{}
	_result, _err = kmssdk.NewClient(nil, openApiConfig)
	return _result, _err
}

func CreateKey(client *kmssdk.Client, enableAutomaticRotation *bool, rotationInterval *string, keyUsage *string, origin *string, description *string, DKMSInstanceId *string, protectionLevel *string, keySpec *string) (_result *kms20160120.CreateKeyResponse, _err error) {
	request := &kms20160120.CreateKeyRequest{
		EnableAutomaticRotation: enableAutomaticRotation,
		RotationInterval:        rotationInterval,
		KeyUsage:                keyUsage,
		Origin:                  origin,
		Description:             description,
		DKMSInstanceId:          DKMSInstanceId,
		ProtectionLevel:         protectionLevel,
		KeySpec:                 keySpec,
	}
	_result = &kms20160120.CreateKeyResponse{}
	_body, _err := client.CreateKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func _main(args []*string) (_err error) {
	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378657.html
	openApiConfig, _err := CreateOpenApiConfig(env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_ID")), env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_SECRET")), tea.String("your region id"))
	if _err != nil {
		return _err
	}

	client, _err := CreateClient(openApiConfig)
	if _err != nil {
		return _err
	}

	enableAutomaticRotation := tea.Bool(false)
	rotationInterval := tea.String("your rotationInterval")
	keyUsage := tea.String("your keyUsage")
	origin := tea.String("your origin")
	description := tea.String("your description")
	dKMSInstanceId := tea.String("your dKMSInstanceId")
	protectionLevel := tea.String("your protectionLevel")
	keySpec := tea.String("your keySpec")
	response, _err := CreateKey(client, enableAutomaticRotation, rotationInterval, keyUsage, origin, description, dKMSInstanceId, protectionLevel, keySpec)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(response))
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}

```

### 3. 既要通过VPC网关进行密钥运算操作又要通过公共网关对KMS资源管理的场景。

#### 参考以下示例代码调用KMS CreateKey API 和 Encrypt API。更多API示例参考 [operation samples](./examples/operation) 和 [manage samples](./examples/manage)

```go
package main

import (
	"os"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	env "github.com/alibabacloud-go/darabonba-env/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	kmssdk "github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
	"github.com/alibabacloud-go/tea/tea"
)

func CreateOpenApiConfig(accessKeyId *string, accessKeySecret *string, regionId *string) (_result *openapi.Config, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
	}
	_result = config
	return _result, _err
}
func CreateKmsInstanceConfig(clientKeyFile *string, password *string, endpoint *string, caFilePath *string) (_result *dedicatedkmsopenapi.Config, _err error) {
	config := &dedicatedkmsopenapi.Config{}
	config.ClientKeyFile = clientKeyFile
	config.Password = password
	config.Endpoint = endpoint
	config.CaFilePath = caFilePath
	_result = config
	return _result, _err
}
func CreateClient(kmsInstanceConfig *dedicatedkmsopenapi.Config, openApiConfig *openapi.Config) (_result *kmssdk.Client, _err error) {
	_result = &kmssdk.Client{}
	_result, _err = kmssdk.NewClient(kmsInstanceConfig, openApiConfig)
	return _result, _err
}

func CreateKey(client *kmssdk.Client, enableAutomaticRotation *bool, rotationInterval *string, keyUsage *string, origin *string, description *string, DKMSInstanceId *string, protectionLevel *string, keySpec *string) (_result *kms20160120.CreateKeyResponse, _err error) {
	request := &kms20160120.CreateKeyRequest{
		EnableAutomaticRotation: enableAutomaticRotation,
		RotationInterval:        rotationInterval,
		KeyUsage:                keyUsage,
		Origin:                  origin,
		Description:             description,
		DKMSInstanceId:          DKMSInstanceId,
		ProtectionLevel:         protectionLevel,
		KeySpec:                 keySpec,
	}
	_result = &kms20160120.CreateKeyResponse{}
	_body, _err := client.CreateKey(request)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func Encrypt(client *kmssdk.Client, paddingMode *string, aad []byte, keyId *string, plaintext []byte, iv []byte, algorithm *string) (_result *dedicatedkmssdk.EncryptResponse, _err error) {
	request := &dedicatedkmssdk.EncryptRequest{
		PaddingMode: paddingMode,
		Aad:         aad,
		KeyId:       keyId,
		Plaintext:   plaintext,
		Iv:          iv,
		Algorithm:   algorithm,
	}
	_result = &dedicatedkmssdk.EncryptResponse{}
	return client.Encrypt(request)
}

func _main(args []*string) (_err error) {
	// 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378657.html
	openApiConfig, _err := CreateOpenApiConfig(env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_ID")), env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_SECRET")), tea.String("your region id"))
	if _err != nil {
		return _err
	}
  kmsInstanceConfig, _err := CreateKmsInstanceConfig(env.GetEnv(tea.String("your client key file path env")), env.GetEnv(tea.String("your client key password env")), tea.String("your kms instance endpoint env"), tea.String("your ca file path"))
	if _err != nil {
		return _err
	}
	client, _err := CreateClient(kmsInstanceConfig, openApiConfig)
	if _err != nil {
		return _err
	}

	enableAutomaticRotation := tea.Bool(false)
	rotationInterval := tea.String("your rotationInterval")
	keyUsage := tea.String("your keyUsage")
	origin := tea.String("your origin")
	description := tea.String("your description")
	dKMSInstanceId := tea.String("your dKMSInstanceId")
	protectionLevel := tea.String("your protectionLevel")
	keySpec := tea.String("your keySpec")
	createKeyResponse, _err := CreateKey(client, enableAutomaticRotation, rotationInterval, keyUsage, origin, description, dKMSInstanceId, protectionLevel, keySpec)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(createKeyResponse))

	paddingMode := tea.String("your paddingMode")
	aad := util.ToBytes(tea.String("your aad"))
	keyId := tea.String("your keyId")
	plaintext := util.ToBytes(tea.String("your plaintext"))
	iv := util.ToBytes(tea.String("your iv"))
	algorithm := tea.String("your algorithm")
	encryptResponse, _err := Encrypt(client, paddingMode, aad, keyId, plaintext, iv, algorithm)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(encryptResponse))
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}

```

### 使用阿里云 SDK访问KMS 1.0密钥运算的用户，需要迁移到KMS 3.0的场景。

#### 参考以下示例代码调用KMS API。更多API示例参考 [kms transfer samples](./examples/transfer)

```go
package example

import (
	"fmt"
	"io/ioutil"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
)

func main() {
	// 创建kms共享网关config并设置相应参数
	config := &openapi.Config{
		// 设置地域Id
		RegionId: tea.String("your-region-id"),
		// 设置访问凭证AccessKeyId
		AccessKeyId: tea.String(os.Getenv("ACCESS_KEY_ID")),
		// 设置访问凭证AccessKeySecret
		AccessKeySecret: tea.String(os.Getenv("ACCESS_KEY_SECRET")),
	}
	// 创建kms实例网关config并设置相应参数
	kmsConfig := &dkmsopenapi.Config{
		// 设置请求协议为https
		Protocol: tea.String("https"),
		// 设置client key文件地址
		ClientKeyFile: tea.String("your-client-key-file-path"),
		// 设置client key密码
		Password: tea.String(os.Getenv("your-client-key-password-env")),
		// 设置kms实例服务地址
		Endpoint: tea.String("your-kms-instance-endpoint"),
	}

	client, err := sdk.NewTransferClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	// 创建密钥调用KMS共享网关
	createKey(client)

	// 生成数据密钥调用KMS实例网关
	generateDataKey(client)
}

// 创建密钥调用KMS共享网关示例，CreateKey请求发送KMS共享网关
func createKey(client *sdk.TransferClient) {
	// 创建密钥请求，设置DKMSInstanceId参数指定KMS实例
	request := &kms20160120.CreateKeyRequest{
		KeySpec:  tea.String("your-key-spec"),
		KeyUsage: tea.String("your-key-usage"),
		// 设置KMS实例ID
		DKMSInstanceId: tea.String("your-kms-instance-id"),
	}

	result, err := client.CreateKey(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}

// 生成数据密钥调用KMS实例网关示例，GenerateDataKey请求默认发送KMS实例网关
func generateDataKey(client *sdk.TransferClient) {
	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:   tea.String("your-key-id"),
		KeySpec: tea.String("your-key-spec"),
		//NumberOfBytes: tea.Int32(32),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &teautil.RuntimeOptions{
		Ca: tea.String(string(ca)),
	}

	result, err := client.GenerateDataKeyWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
```
## KMS实例性能测试

如果需要使用KMS实例SDK进行KMS实例性能测试，请参考benchmarks目录下的压力测试工具示例代码，编译成可执行程序以后使用如下命令运行:

```shell
nohup ./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600 -log_path=./log > aes_256_enc.out 2>&1&
```

压力测试工具如何编译以及使用请参考[文档](README-benchmark-CN.md)。

## 许可证

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
