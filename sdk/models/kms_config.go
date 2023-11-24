package models

import (
	"github.com/alibabacloud-go/tea/tea"
	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
)

type KmsConfig struct {
	*dkmsopenapi.Config
	//默认使用KMS共享网关的接口API Name列表
	DefaultKmsApiNames []string
	//指定所有接口使用到的字符集编码
	Encoding *string
	// 高级接口开关 默认使用高级接口
	AdvanceSwitch bool
	// 强制使用低级转换接口 默认是false,即使用高级转换接口
	ForceLowVersionCryptoTransfer bool
}

func TransferKmsConfig(config interface{}) (*KmsConfig, error) {
	kmsConfig := &KmsConfig{}
	switch c := config.(type) {
	case *KmsConfig:
		kmsConfig = c
	case KmsConfig:
		kmsConfig = &c
	case *dkmsopenapi.Config:
		kmsConfig.Config = c
	case dkmsopenapi.Config:
		kmsConfig.Config = &c
	default:
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Not support config param type.",
		})
	}

	return kmsConfig, nil
}
