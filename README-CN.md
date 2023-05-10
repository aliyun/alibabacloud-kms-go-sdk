[English](README.md) | 简体中文

# 阿里云KMS Go SDK

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

阿里云KMS Go SDK可以帮助Golang开发者快速使用KMS。

- [阿里云专属KMS主页](https://help.aliyun.com/document_detail/311016.html)
- [代码示例](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-kms-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-kms-go-sdk/releases)

## 优势

帮助Golang开发者通过本SDK快速使用阿里云KMS产品的所有API:

- 通过KMS共享网关访问KMS OpenAPI
- 通过KMS实例网关访问KMS实例提供的API

## 软件要求

- Golang 1.13及以上。

## 安装

您可以使用`go mod`管理您的依赖：

```
require (
	github.com/aliyun/alibabacloud-kms-go-sdk v1.0.0
)
```

或者，通过`go get`命令获取远程代码包：

```
$ go get -u github.com/aliyun/alibabacloud-kms-go-sdk
```

## 客户端机制

阿里云KMS Go SDK支持调用KMS共享网关和KMS实例网关提供的API。
阿里云KMS Go SDK默认情况下会将以下API的请求转发到KMS实例网关，其它KMS API则发送到KMS共享网关。

* Encrypt
* Decrypt
* GenerateDataKey
* GenerateDataKeyWithoutPlaintext
* GetPublicKey
* AsymmetricEncrypt
* AsymmetricDecrypt
* AsymmetricSign
* AsymmetricVerify
* GetSecretValue

阿里云KMS Go SDK也支持将以上的API请求发送到KMS共享网关，具体方法请参考使用示例-[特殊使用场景](#特殊使用场景)

## 使用示例

### 常规使用场景

#### 场景一 可以参考下面的代码调用KMS共享网关和KMS实例网关的服务。

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

	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	// 创建密钥调用KMS共享网关
	createKey(client)

	// 生成数据密钥调用KMS实例网关
	generateDataKey(client)
}

// 创建密钥调用KMS共享网关示例，CreateKey请求发送KMS共享网关
func createKey(client *sdk.Client) {
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
func generateDataKey(client *sdk.Client) {
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

#### 场景二 可以参考下面的代码仅调用KMS实例网关的服务。

```go
package example

import (
	"fmt"
	"io/ioutil"
	"os"

	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
)

func main() {
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

	client, err := sdk.NewClient(nil, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &teautil.RuntimeOptions{
		Ca: tea.String(string(ca)),
	}

	result, err := client.EncryptWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
```

#### 场景三 可以参考下面的代码仅调用KMS共享网关的服务。

```go
package example

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"

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

	client, err := sdk.NewClient(config, nil)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	result, err := client.Encrypt(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
```

### 特殊使用场景

#### 场景一 参考如下代码可以将所有接口的调用请求(包括默认转发到KMS实例网关的)发送到KMS共享网关。

```go
package example

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
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

	// 创建KmsClient
	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	// 设置参数IsUseKmsShareGateway为True, 将所有接口调用请求发送到KMS共享网关
	client.SetIsUseKmsShareGateway(true)

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}

	result, err := client.GetSecretValue(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
```

#### 参考如下代码将调用GetSecretValue接口请求发送到KMS共享网关。

```go
package example

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
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
	// 创建KmsConfig配置并设置相应参数
	kmsConfig := &models.KmsConfig{
		Config: &dkmsopenapi.Config{
			// 设置请求协议为https
			Protocol: tea.String("https"),
			// 设置client key文件地址
			ClientKeyFile: tea.String("your-client-key-file-path"),
			// 设置client key密码
			Password: tea.String(os.Getenv("your-client-key-password-env")),
			// 设置kms实例服务地址
			Endpoint: tea.String("your-kms-instance-endpoint"),
		},
		// 设置指定的API接口转发到KMS共享网关
		DefaultKmsApiNames: []string{"GetSecretValue"},
	}

	// 创建KmsClient
	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}

	result, err := client.GetSecretValue(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
```

#### 场景三 参考如下代码将单独一次调用请求发送到KMS共享网关。

```go
package example

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
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
	// 创建KmsConfig配置并设置相应参数
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

	// 创建KmsClient
	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}

	// 如果忽略ssl验证，可以在RuntimeOptions设置IgnoreSSL为true
	runtime := &models.KmsRuntimeOptions{
		RuntimeOptions: &teautil.RuntimeOptions{
			//IgnoreSSL: tea.Bool(true),
		},
		// 如果您设定IsUseKmsShareGateway为true，那么这次调用将被转发到共享kms网关
		IsUseKmsShareGateway: tea.Bool(true),
	}

	result, err := client.GetSecretValueWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
```

## 字符编码设置说明(默认为UTF-8)

- 您可以参考以下代码示例，设置全局的字符集编码。

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
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
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
	// 创建KmsConfig配置并设置相应参数
	kmsConfig := &models.KmsConfig{
		Config: &dkmsopenapi.Config{
			// 设置请求协议为https
			Protocol: tea.String("https"),
			// 设置client key文件地址
			ClientKeyFile: tea.String("your-client-key-file-path"),
			// 设置client key密码
			Password: tea.String(os.Getenv("your-client-key-password-env")),
			// 设置kms实例服务地址
			Endpoint: tea.String("your-kms-instance-endpoint"),
		},
		// 设置字符集编码为utf-8
		Encoding: tea.String("utf-8"),
	}

	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &teautil.RuntimeOptions{
		Ca: tea.String(string(ca)),
	}

	result, err := client.EncryptWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
```

- 您可以参考以下代码示例，设置单独一次请求的字符集编码。

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
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
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
	// 创建KmsConfig配置并设置相应参数
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

	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &models.KmsRuntimeOptions{
		RuntimeOptions: &teautil.RuntimeOptions{
			Ca: tea.String(string(ca)),
		},
		// 设置字符集编码为utf-8
		Encoding: tea.String("utf-8"),
	}

	result, err := client.EncryptWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
```

## SSL证书验证开关设置(默认验证SSL证书)

您可以参考以下代码示例，设置不验证HTTPS SSL证书，例如在开发测试时，以简化程序。

```go
package example

import (
	"fmt"
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
	// 创建KmsConfig配置并设置相应参数
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

	// 创建KMS client
	client, err := sdk.NewClient(config, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}

	// 忽略ssl验证，在RuntimeOptions设置IgnoreSSL为true
	runtime := &teautil.RuntimeOptions{
		IgnoreSSL: tea.Bool(true),
	}

	result, err := client.GetSecretValueWithOptions(request, runtime)
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
