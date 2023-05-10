package models

import (
	"github.com/alibabacloud-go/tea/tea"
	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
)

type KmsConfig struct {
	*dkmsopenapi.Config
	DefaultKmsApiNames []string
	Encoding           *string
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
