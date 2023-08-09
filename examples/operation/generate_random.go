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

func GenerateRandom (client *kmssdk.Client, length *int32) (_result *dedicatedkmssdk.GenerateRandomResponse, _err error) {
  request := &dedicatedkmssdk.GenerateRandomRequest{
    Length: length,
  }
  _result = &dedicatedkmssdk.GenerateRandomResponse{}
  return client.GenerateRandom(request)
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

  parseIntTmp, err := util.AssertAsString(tea.String("your length"))
  if err != nil {
    _err = err
    return _err
  }
  length := number.ParseInt(parseIntTmp)
  response, _err := GenerateRandom(client, tea.ToInt32(length))
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
