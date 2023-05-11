package main

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

// 新接入用户可以参考此方法调用KMS实例网关服务。
func NewUserGenerateDataKeySample() {
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

	GenerateDataKey(client)
}

// 密钥迁移前示例代码。
func BeforeMigrateGenerateDataKeySample() {
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

	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:   tea.String("your-key-id"),
		KeySpec: tea.String("your-key-spec"),
		//NumberOfBytes: tea.Int32(32),
	}

	result, err := client.GenerateDataKey(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.String())
}

// 密钥迁移后示例代码。
func AfterMigrateGenerateDataKeySample() {
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

	GenerateDataKey(client)
}

func GenerateDataKey(client *sdk.Client) {
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

func main() {
	NewUserGenerateDataKeySample()
	BeforeMigrateGenerateDataKeySample()
	AfterMigrateGenerateDataKeySample()
}
