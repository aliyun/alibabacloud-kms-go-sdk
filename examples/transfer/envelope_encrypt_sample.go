package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
)

const gcmIvLength = 12

// The envelope cipher text may be stored
type EnvelopeCipherText struct {
	EncryptedDataKey string
	CipherText       string
	Iv               string
}

// 新接入用户可以参考此方法调用KMS实例网关服务。
func NewUserEnvelopeEncryptSample() {
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

	client, err := sdk.NewTransferClient(nil, kmsConfig)
	if err != nil {
		panic(err)
	}

	EnvelopeEncrypt(client)
}

// 密钥迁移前示例代码。
func BeforeMigrateEnvelopeEncryptSample() {
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

	// 本地待加密数据
	data := []byte("your plaintext data")

	// 获取数据密钥，下面以<your-key-id>密钥为例进行说明，数据密钥长度32字节
	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:         tea.String("your-key-id"),
		NumberOfBytes: tea.Int32(32),
	}

	// 调用生成数据密钥接口
	generateDataKeyResponse, err := client.GenerateDataKey(request)
	if err != nil {
		panic(err)
	}

	// KMS返回的数据密钥密文，解密本地数据密文时，先将数据密钥密文解密后使用
	encryptedDataKey := generateDataKeyResponse.Body.CiphertextBlob
	// KMS返回的数据密钥明文, 加密本地数据使用
	plainDataKey, err := base64.StdEncoding.DecodeString(tea.StringValue(generateDataKeyResponse.Body.Plaintext))
	if err != nil {
		panic(err)
	}

	// 计算本地加密初始向量，解密时需要传入
	iv := make([]byte, gcmIvLength)
	rand.Read(iv)

	// 本地加密数据，下面是以AES-256 GCM模式为例
	block, err := aes.NewCipher(plainDataKey)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	cipherText := gcm.Seal(nil, iv, data, nil)

	// 输出密文，密文输出或持久化由用户根据需要进行处理，下面示例仅展示将密文输出到一个对象的情况
	// 假如envelopeCipherText是需要输出的密文对象，至少需要包括以下三个内容:
	// (1) encryptedDataKey: 专属KMS返回的数据密钥密文
	// (2) cipherText: 密文数据
	// (3) iv: 加密初始向量
	envelopeCipherText := &EnvelopeCipherText{
		EncryptedDataKey: tea.StringValue(encryptedDataKey),
		CipherText:       base64.StdEncoding.EncodeToString(cipherText),
		Iv:               base64.StdEncoding.EncodeToString(iv),
	}

	fmt.Println(envelopeCipherText)
}

// 密钥迁移后示例代码。
func AfterMigrateEnvelopeEncryptSample() {
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

	EnvelopeEncrypt(client)
}

func EnvelopeEncrypt(client *sdk.TransferClient) {
	// 本地待加密数据
	data := []byte("your plaintext data")

	// 获取数据密钥，下面以<your-key-id>密钥为例进行说明，数据密钥长度32字节
	request := &kms20160120.GenerateDataKeyRequest{
		KeyId:         tea.String("your-key-id"),
		NumberOfBytes: tea.Int32(32),
	}

	// 验证服务器证书，在RuntimeOptions设置ca证书
	ca, err := ioutil.ReadFile("your-ca-certificate-file-path")
	if err != nil {
		panic(err)
	}
	runtime := &teautil.RuntimeOptions{
		Ca: tea.String(string(ca)),
	}

	// 调用生成数据密钥接口
	generateDataKeyResponse, err := client.GenerateDataKeyWithOptions(request, runtime)
	if err != nil {
		panic(err)
	}

	// KMS返回的数据密钥密文，解密本地数据密文时，先将数据密钥密文解密后使用
	encryptedDataKey := generateDataKeyResponse.Body.CiphertextBlob
	// KMS返回的数据密钥明文, 加密本地数据使用
	plainDataKey, err := base64.StdEncoding.DecodeString(tea.StringValue(generateDataKeyResponse.Body.Plaintext))
	if err != nil {
		panic(err)
	}

	// 计算本地加密初始向量，解密时需要传入
	iv := make([]byte, gcmIvLength)
	rand.Read(iv)

	// 本地加密数据，下面是以AES-256 GCM模式为例
	block, err := aes.NewCipher(plainDataKey)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	cipherText := gcm.Seal(nil, iv, data, nil)

	// 输出密文，密文输出或持久化由用户根据需要进行处理，下面示例仅展示将密文输出到一个对象的情况
	// 假如envelopeCipherText是需要输出的密文对象，至少需要包括以下三个内容:
	// (1) encryptedDataKey: 专属KMS返回的数据密钥密文
	// (2) cipherText: 密文数据
	// (3) iv: 加密初始向量
	envelopeCipherText := &EnvelopeCipherText{
		EncryptedDataKey: tea.StringValue(encryptedDataKey),
		CipherText:       base64.StdEncoding.EncodeToString(cipherText),
		Iv:               base64.StdEncoding.EncodeToString(iv),
	}

	fmt.Println(envelopeCipherText)

}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	NewUserEnvelopeEncryptSample()
	BeforeMigrateEnvelopeEncryptSample()
	AfterMigrateEnvelopeEncryptSample()

}
