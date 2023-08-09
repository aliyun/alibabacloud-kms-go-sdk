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

	// 创建密钥请求
	request := &kms20160120.CreateKeyRequest{
		KeySpec:        tea.String("your-key-spec"),
		KeyUsage:       tea.String("your-key-usage"),
		DKMSInstanceId: tea.String("your-kms-instance-id"),
	}

	result, err := client.CreateKey(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())

}
