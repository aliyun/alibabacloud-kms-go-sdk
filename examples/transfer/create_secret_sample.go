package main

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

	client, err := sdk.NewTransferClient(config, nil)
	if err != nil {
		panic(err)
	}

	// 创建凭据请求
	request := &kms20160120.CreateSecretRequest{
		SecretName: tea.String("your-secret-name"),
		SecretData: tea.String("your-secret-data"),
		VersionId:  tea.String("your-version-id"),
		// 设置KMS实例ID
		DKMSInstanceId: tea.String("your-kms-instance-id"),
		// 设置加密保护凭据值的KMS用户主密钥ID
		EncryptionKeyId: tea.String("your-encryption-key-id"),
	}

	result, err := client.CreateSecret(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}
