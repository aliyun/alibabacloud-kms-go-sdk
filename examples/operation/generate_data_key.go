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
  number  "github.com/alibabacloud-go/darabonba-number/client"
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

func GenerateDataKey (client *kmssdk.Client, aad []byte, keyId *string, numberOfBytes *int32, algorithm *string) (_result *dedicatedkmssdk.GenerateDataKeyResponse, _err error) {
  request := &dedicatedkmssdk.GenerateDataKeyRequest{
    Aad: aad,
    KeyId: keyId,
    NumberOfBytes: numberOfBytes,
    Algorithm: algorithm,
  }
  _result = &dedicatedkmssdk.GenerateDataKeyResponse{}
  return client.GenerateDataKey(request)
}

func _main (args []*string) (_err error) {
  kmsInstanceConfig, _err := CreateKmsInstanceConfig(env.GetEnv(tea.String("your client key file path env")), env.GetEnv(tea.String("your client key password env")), tea.String("your kms instance endpoint env"), tea.String("your ca file path"))
  if _err != nil {
    return _err
  }

  client, _err := CreateClient(kmsInstanceConfig)
  if _err != nil {
    return _err
  }

  aad := util.ToBytes(tea.String("your aad"))
  keyId := tea.String("your keyId")
  parseIntTmp, err := util.AssertAsString(tea.String("your numberOfBytes"))
  if err != nil {
    _err = err
    return _err
  }
  numberOfBytes := number.ParseInt(parseIntTmp)
  algorithm := tea.String("your algorithm")
  response, _err := GenerateDataKey(client, aad, keyId, tea.ToInt32(numberOfBytes), algorithm)
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
