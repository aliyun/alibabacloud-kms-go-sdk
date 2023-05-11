English | [简体中文](README-CN.md)

# Alibaba Cloud KMS Go SDK

![](https://aliyunsdk-pages.alicdn.com/icons/AlibabaCloud.svg)

Alibaba Cloud KMS Go SDK can help Golang developers to use the KMS.

- [Alibaba Cloud Dedicated KMS Homepage](https://www.alibabacloud.com/help/zh/doc-detail/311016.htm)
- [Sample Code](/examples)
- [Issues](https://github.com/aliyun/alibabacloud-kms-go-sdk/issues)
- [Release](https://github.com/aliyun/alibabacloud-kms-go-sdk/releases)

## Feature

- Access KMS services through the KMS shared gateway.
- Access KMS services through the KMS instance gateway.

## Requirements

- Golang 1.13 or later.

## Installation

If you use `go mod` to manage your dependence, You can declare the dependency on AlibabaCloud KMS Go SDK in the `go.mod`
file:

```text
require (
	github.com/aliyun/alibabacloud-kms-go-sdk v1.0.0
)
```

Or, Run the following command to get the remote code package:

```shell
$ go get -u github.com/aliyun/alibabacloud-kms-go-sdk
```

## Client Mechanism

Alibaba Cloud KMS Go SDK supports to call APIs provided by the KMS shared gateway and KMS instance gateway.
By default, Alibaba Cloud KMS Go SDK sends requests for the following APIs to the KMS instance gateway, and the other
APIs to the KMS shared gateway.

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

Alibaba Cloud KMS Go SDK also supports sending API requests from the above APIs to the KMS shared gateway. For details,
see Sample code - [Special Usage scenarios](#Special usage scenarios).

## Sample code

### General usage scenarios

#### Scenario 1 You can refer to the following code to call APIs of the KMS shared gateway and the KMS instance gateway in the default way.

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

	// create key
	createKey(client)

	// generate data key
	generateDataKey(client)
}

func createKey(client *sdk.Client) {
	request := &kms20160120.CreateKeyRequest{
		KeySpec:  tea.String("your-key-spec"),
		KeyUsage: tea.String("your-key-usage"),
		// set kms instance id
		DKMSInstanceId: tea.String("your-kms-instance-id"),
	}

	result, err := client.CreateKey(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}

func generateDataKey(client *sdk.Client) {
	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:   tea.String("your-key-id"),
		KeySpec: tea.String("your-key-spec"),
		//NumberOfBytes: tea.Int32(32),
	}

	// If verify server CA certificate,you can set CA certificate file path with RuntimeOptions
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

#### Scenario 2 You can refer to the following code to call only the APIs of the KMS instance gateway.

```go
package example

import (
	"fmt"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"io/ioutil"
	"os"

	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
)

func main() {
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
	client, err := sdk.NewClient(nil, kmsConfig)
	if err != nil {
		panic(err)
	}

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	// If verify server CA certificate,you can set CA certificate file path with RuntimeOptions
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

#### Scenario 3 You can refer to the following code to call only the KMS shared gateway API.

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
	// set config
	config := &openapi.Config{
		// set region id
		RegionId: tea.String("your-region-id"),
		// set access key id
		AccessKeyId: tea.String(os.Getenv("ACCESS_KEY_ID")),
		// set access key secret
		AccessKeySecret: tea.String(os.Getenv("ACCESS_KEY_SECRET")),
	}

	// create KMS client
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

### Special usage scenarios

#### Scenario 1 Refer to the following code to forward calls from all of these API to the KMS shared gateway.

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

	// set parameter IsUseKmsShareGateway with true, and forward all interfaces to the KMS shared gateway
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

#### Scenario 2 Refer to the following code to forward the call to a specific API (GetSecretValue) to the KMS shared gateway.

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
	kmsConfig := &models.KmsConfig{
		Config: &dkmsopenapi.Config{
			// set the request protocol to https
			Protocol: tea.String("https"),
			// set client key file path
			ClientKeyFile: tea.String("your-client-key-file-path"),
			// set client key password
			Password: tea.String(os.Getenv("your-client-key-password-env")),
			// set kms instance endpoint
			Endpoint: tea.String("your-kms-instance-endpoint"),
		},
		// set the specified API interface to forward to KMS shared gateway
		DefaultKmsApiNames: []string{"GetSecretValue"},
	}

	// create KMS client
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

#### Scenario 3 Refer to the following code to forward a single call to the KMS shared gateway.

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

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}

	// If you ignore ssl verification，you can set IgnoreSSL with true related to the RuntimeOptions parameter
	runtime := &models.KmsRuntimeOptions{
		RuntimeOptions: &teautil.RuntimeOptions{
			//IgnoreSSL: tea.Bool(true),
		},
		// If you set IsUseKmsShareGateway with true,the request must be sent to the shared KMS gateway
		IsUseKmsShareGateway: tea.Bool(true),
	}

	result, err := client.GetSecretValueWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
```

## Character encoding setting instructions (default UTF-8)

- You can refer to the following code example to set the global character set encoding.

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

	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk/models"
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
	kmsConfig := &models.KmsConfig{
		Config: &dkmsopenapi.Config{
			// set the request protocol to https
			Protocol: tea.String("https"),
			// set client key file path
			ClientKeyFile: tea.String("your-client-key-file-path"),
			// set client key password
			Password: tea.String(os.Getenv("your-client-key-password-env")),
			// set kms instance endpoint
			Endpoint: tea.String("your-kms-instance-endpoint"),
		},
		// set encoding to utf-8
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

	// If verify server CA certificate,you can set CA certificate file path with RuntimeOptions
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

- You can refer to the following code example to set the character set encoding for a single request.

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

	request := &kms20160120.EncryptRequest{
		KeyId:     tea.String("your-key-id"),
		Plaintext: tea.String("your-plaintext"),
	}

	// If verify server CA certificate,you can set CA certificate file path with RuntimeOptions
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &models.KmsRuntimeOptions{
		RuntimeOptions: &teautil.RuntimeOptions{
			Ca: tea.String(string(ca)),
		},
		// set encoding to utf-8
		Encoding: tea.String("utf-8"),
	}
	
	result, err := client.EncryptWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
```

## SSL certificate validation switch setting (default validation of SSL certificate)

You can refer to the following code example to set the HTTPS SSL certificate not to be validated, for example when
developing tests, to simplify the program.

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

	request := &kms20160120.GetSecretValueRequest{
		SecretName:   tea.String("your-secret-name"),
		VersionId:    tea.String("your-version-id"),
		VersionStage: tea.String("your-version-stage"),
	}
	
	// If you ignore ssl verification，you can set IgnoreSSL with true related to the RuntimeOptions parameter
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

## KMS instance performance testing

If you need to use the KMS instance SDK for KMS instance performance testing, please refer to the sample code of the pressure measurement tools in the directory named benchmarks , compile it into an executable program and run it with the following command:

```shell
nohup ./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600 -log_path=./log > aes_256_enc.out 2>&1&
```

How to compile and use the stress test tool, please refer to [the document](README-benchmark.md). 

## License

[Apache-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Copyright (c) 2009-present, Alibaba Cloud All rights reserved.
 
