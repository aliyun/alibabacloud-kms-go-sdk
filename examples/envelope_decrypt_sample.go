package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

// The envelope cipher text may be stored
type EnvelopeCipherPersistObject struct {
	EncryptedDataKey string
	CipherText       string
	Iv               string
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	NewUserEnvelopeDecryptSample()
	BeforeMigrateEnvelopeDecryptSample()
	AfterMigrateEnvelopeDecryptSample()

}

// 新接入用户可以参考此方法调用KMS实例网关服务。
func NewUserEnvelopeDecryptSample() {
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

	EnvelopeDecrypt(client)
}

// 密钥迁移前示例代码。
func BeforeMigrateEnvelopeDecryptSample() {
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

	// 读取封信加密持久化对象
	envelopeCipherPersistObject := getEnvelopeCipherPersistObject()

	// 数据密钥密文
	encryptedDataKey := envelopeCipherPersistObject.EncryptedDataKey

	// 从封信加密持久化对象获取数据密钥密文，调用KMS在线解密
	request := &kms20160120.DecryptRequest{
		CiphertextBlob: tea.String(encryptedDataKey),
	}

	decryptResponse, err := client.Decrypt(request)
	if err != nil {
		panic(err)
	}

	// 数据密钥明文
	plainDataKey, err := base64.StdEncoding.DecodeString(tea.StringValue(decryptResponse.Body.Plaintext))
	if err != nil {
		panic(err)
	}
	// 本地加密时使用的初始向量, 在本地解密数据时需要传入
	iv, err := base64.StdEncoding.DecodeString(envelopeCipherPersistObject.Iv)
	if err != nil {
		panic(err)
	}
	// 待解密的本地数据密文
	ciphertext, err := base64.StdEncoding.DecodeString(envelopeCipherPersistObject.CipherText)
	if err != nil {
		panic(err)
	}

	// 使用数据密钥明文在本地进行解密, 下面是以AES-256 GCM模式为例
	block, err := aes.NewCipher(plainDataKey)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	decryptedData, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(decryptedData))
}

// 密钥迁移后示例代码。
func AfterMigrateEnvelopeDecryptSample() {
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

	EnvelopeDecrypt(client)
}

func EnvelopeDecrypt(client *sdk.Client) {
	// 读取封信加密持久化对象
	envelopeCipherPersistObject := getEnvelopeCipherPersistObject()

	// 数据密钥密文
	encryptedDataKey := envelopeCipherPersistObject.EncryptedDataKey

	// 从封信加密持久化对象获取数据密钥密文，调用KMS在线解密
	request := &kms20160120.DecryptRequest{
		CiphertextBlob: tea.String(encryptedDataKey),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &teautil.RuntimeOptions{
		Ca: tea.String(string(ca)),
	}

	decryptResponse, err := client.DecryptWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	// 数据密钥明文
	plainDataKey, err := base64.StdEncoding.DecodeString(tea.StringValue(decryptResponse.Body.Plaintext))
	if err != nil {
		panic(err)
	}
	// 本地加密时使用的初始向量, 在本地解密数据时需要传入
	iv, err := base64.StdEncoding.DecodeString(envelopeCipherPersistObject.Iv)
	if err != nil {
		panic(err)
	}
	// 待解密的本地数据密文
	ciphertext, err := base64.StdEncoding.DecodeString(envelopeCipherPersistObject.CipherText)
	if err != nil {
		panic(err)
	}

	// 使用数据密钥明文在本地进行解密, 下面是以AES-256 GCM模式为例
	block, err := aes.NewCipher(plainDataKey)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	decryptedData, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(decryptedData))
}

func getEnvelopeCipherPersistObject() *EnvelopeCipherPersistObject {
	// TODO 用户需要在此处代码进行替换，从存储中读取封信加密持久化对象
	return &EnvelopeCipherPersistObject{}
}
