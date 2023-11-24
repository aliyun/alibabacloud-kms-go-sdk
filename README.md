English | [简体中文](README-CN.md)

# Alibaba Cloud KMS Go SDK

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

Alibaba Cloud KMS Go SDK can help Golang developers to use the KMS.

- [Alibaba Cloud Dedicated KMS Homepage](https://www.alibabacloud.com/help/zh/doc-detail/311016.htm)
- [Sample Code](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-kms-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-kms-go-sdk/releases)

## Advantage

Alibaba Cloud KMS SDK helps Java developers quickly use all APIs of Alibaba Cloud KMS products:
- KMS resource management and key operations can be performed through KMS public gateway access
- You can perform key operations through KMS instance gateway

## Requirements

- Golang 1.13 or later.

## Installation

If you use `go mod` to manage your dependence, You can declare the dependency on AlibabaCloud KMS Go SDK in the `go.mod`
file:

```text
require (
	github.com/aliyun/alibabacloud-kms-go-sdk v1.0.1
)
```

Or, Run the following command to get the remote code package:

```shell
$ go get -u github.com/aliyun/alibabacloud-kms-go-sdk
```

## Introduction to KMS Client

| KMS client struct | Introduction | Usage scenarios |
| :-----| :---- | :---- |
| Client | KMS resource management and key operations for KMS instance gateways are supported | 1. Scenarios where key operations are performed only through VPC gateways. <br>2. KMS resource management scenarios that only use public gateways. <br>3. Scenarios where you want to perform key operations through VPC gateways and manage KMS resources through public gateways.|
| TransferClient | Users can migrate from KMS 1.0 key operations to KMS 3.0 key operations | Users who use Alibaba Cloud SDK to access KMS 1.0 key operations need to migrate to KMS 3.0 |


## Sample code

### 1. Scenarios where key operations are performed only through VPC gateways.
#### Refer to the following sample code to call the KMS Encrypt API. For more API examples, see [operation samples](./examples/operation)
```go
package example

import (
	"os"
	console  "github.com/alibabacloud-go/tea-console/client"
	env  "github.com/alibabacloud-go/darabonba-env/client"
	util  "github.com/alibabacloud-go/tea-utils/v2/service"
	kmssdk  "github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	dedicatedkmsopenapi  "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dedicatedkmssdk  "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
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

### 2. KMS resources are managed only through public gateways.
#### Refer to the following sample code to call the KMS CreateKey API. For more API examples, see [manage samples](./examples/manage)

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

### 3. You must not only perform key operations through a VPC gateway, but also manage KMS resources through a public gateway.
#### Refer to the following sample code to call the KMS CreateKey API and the Encrypt API. For more API examples, see [operation samples](./examples/operation) and [manage samples](./examplesmanage)

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

func Encrypt (client *kmssdk.Client, paddingMode *string, aad []byte, keyId *string, plaintext []byte, iv []byte, algorithm *string) (_result *dedicatedkmssdk.EncryptResponse, _err error) {
	request := &dedicatedkmssdk.EncryptRequest{
		PaddingMode: paddingMode,
		Aad: aad,
		KeyId: keyId,
		Plaintext: plaintext,
		Iv: iv,
		Algorithm: algorithm,
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

### Users who use Alibaba Cloud SDK to access KMS 1.0 key operations need to migrate to KMS 3.0.
#### Refer to the following sample code to call the KMS API. For more API examples, see [kms transfer samples](./examples/transfer)

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
	// set config
	config := &openapi.Config{
		// set region id
		RegionId: tea.String("your-region-id"),
		// set access key id
		AccessKeyId: tea.String(os.Getenv("ACCESS_KEY_ID")),
		// set access key secret
		AccessKeySecret: tea.String(os.Getenv("ACCESS_KEY_SECRET")),
	}
	// set kms config
	kmsConfig := &dkmsopenapi.Config{
		// set the request protocol to https
		Protocol: tea.String("https"),
		// set client key file path
		ClientKeyFile: tea.String("your-client-key-file-path"),
		// set client key password
		Password: tea.String(os.Getenv("your-client-key-password-env")),
		// set kms instance endpoint
		Endpoint: tea.String("your-kms-instance-endpoint"),
	}

	// create KMS client
	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	// Create a key and invoke the KMS shared gateway
	createKey(client)

	// Generate a data key and call the KMS instance gateway
	generateDataKey(client)
}

// In the example of creating a key to call a KMS shared gateway, CreateKey requests the KMS shared gateway
func createKey(client *sdk.TransferClient) {
	// Create a key request and set the DKMSInstanceId parameter to specify the KMS instance
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

// Generate a data key to call the KMS instance gateway example, and GenerateDataKey requests to send the KMS instance gateway by default
func generateDataKey(client *sdk.TransferClient) {
	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:   tea.String("your-key-id"),
		KeySpec: tea.String("your-key-spec"),
		//NumberOfBytes: tea.Int32(32),
	}

	// Verify the server certificate and set the ca certificate in RuntimeOptions
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

## KMS instance performance testing

If you need to use the KMS instance SDK for KMS instance performance testing, please refer to the sample code of the pressure measurement tools in the directory named benchmarks , compile it into an executable program and run it with the following command:

```shell
nohup ./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600 -log_path=./log > aes_256_enc.out 2>&1&
```

How to compile and use the stress test tool, please refer to [the document](README-benchmark.md).

## License

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
 
