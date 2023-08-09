// This file is auto-generated, don't edit it. Thanks.
package main

import (
  "os"
  openapi  "github.com/alibabacloud-go/darabonba-openapi/v2/client"
  console  "github.com/alibabacloud-go/tea-console/client"
  env  "github.com/alibabacloud-go/darabonba-env/client"
  util  "github.com/alibabacloud-go/tea-utils/v2/service"
  kmssdk  "github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
  kms20160120  "github.com/alibabacloud-go/kms-20160120/v3/client"
  "github.com/alibabacloud-go/tea/tea"
)


func CreateOpenApiConfig (accessKeyId *string, accessKeySecret *string, regionId *string) (_result *openapi.Config, _err error) {
  config := &openapi.Config{
    AccessKeyId: accessKeyId,
    AccessKeySecret: accessKeySecret,
    RegionId: regionId,
  }
  _result = config
  return _result , _err
}

func CreateClient (openApiConfig *openapi.Config) (_result *kmssdk.Client, _err error) {
  _result = &kmssdk.Client{}
  _result, _err = kmssdk.NewClient(nil, openApiConfig)
  return _result, _err
}

func GetParametersForImport (client *kmssdk.Client, wrappingKeySpec *string, keyId *string, wrappingAlgorithm *string) (_result *kms20160120.GetParametersForImportResponse, _err error) {
  request := &kms20160120.GetParametersForImportRequest{
    WrappingKeySpec: wrappingKeySpec,
    KeyId: keyId,
    WrappingAlgorithm: wrappingAlgorithm,
  }
  _result = &kms20160120.GetParametersForImportResponse{}
  _body, _err := client.GetParametersForImport(request)
  if _err != nil {
    return _result, _err
  }
  _result = _body
  return _result, _err
}

func _main (args []*string) (_err error) {
  // 请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID 和 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
  // 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例使用环境变量获取 AccessKey 的方式进行调用，仅供参考，建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378657.html
  openApiConfig, _err := CreateOpenApiConfig(env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_ID")), env.GetEnv(tea.String("ALIBABA_CLOUD_ACCESS_KEY_SECRET")), tea.String("your region id"))
  if _err != nil {
    return _err
  }

  client, _err := CreateClient(openApiConfig)
  if _err != nil {
    return _err
  }

  wrappingKeySpec := tea.String("your wrappingKeySpec")
  keyId := tea.String("your keyId")
  wrappingAlgorithm := tea.String("your wrappingAlgorithm")
  response, _err := GetParametersForImport(client, wrappingKeySpec, keyId, wrappingAlgorithm)
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
