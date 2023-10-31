// This file is auto-generated, don't edit it. Thanks.
package main

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


func CreateKmsInstanceConfig (clientKeyFile *string, password *string, endpoint *string, caFilePath *string) (_result *dedicatedkmsopenapi.Config, _err error) {
  config := &dedicatedkmsopenapi.Config{
    ClientKeyFile: clientKeyFile,
    Password: password,
    Endpoint: endpoint,
    CaFilePath: caFilePath,
  }
  _result = config
  return _result , _err
}

func CreateClient (kmsInstanceConfig *dedicatedkmsopenapi.Config) (_result *kmssdk.Client, _err error) {
  _result = &kmssdk.Client{}
  _result, _err = kmssdk.NewClient(kmsInstanceConfig, nil)
  return _result, _err
}

func AdvanceEncrypt (client *kmssdk.Client, paddingMode *string, aad []byte, keyId *string, plaintext []byte, iv []byte, algorithm *string) (_result *dedicatedkmssdk.AdvanceEncryptResponse, _err error) {
  request := &dedicatedkmssdk.AdvanceEncryptRequest{
    PaddingMode: paddingMode,
    Aad: aad,
    KeyId: keyId,
    Plaintext: plaintext,
    Iv: iv,
    Algorithm: algorithm,
  }
  _result = &dedicatedkmssdk.AdvanceEncryptResponse{}
  return client.AdvanceEncrypt(request)
}

func _main (args []*string) (_err error) {
  kmsInstanceConfig, _err := CreateKmsInstanceConfig(env.GetEnv(tea.String("your client key file path env")), env.GetEnv(tea.String("your client key password env")), tea.String("your kms instance endpoint"), tea.String("your ca file path"))
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
  response, _err := AdvanceEncrypt(client, paddingMode, aad, keyId, plaintext, iv, algorithm)
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
