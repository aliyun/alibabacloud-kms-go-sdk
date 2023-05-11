package models

import (
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	dkmsutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
)

type KmsRuntimeOptions struct {
	*teautil.RuntimeOptions
	IsUseKmsShareGateway *bool
	Encoding             *string
}

func TransferKmsRuntimeOptions(runtime interface{}) (*KmsRuntimeOptions, error) {
	kmsRuntime := &KmsRuntimeOptions{}
	switch r := runtime.(type) {
	case *KmsRuntimeOptions:
		kmsRuntime = r
	case KmsRuntimeOptions:
		kmsRuntime = &r
	case *teautil.RuntimeOptions:
		kmsRuntime.RuntimeOptions = r
	case teautil.RuntimeOptions:
		kmsRuntime.RuntimeOptions = &r
	default:
		return nil, tea.NewSDKError(map[string]interface{}{
			"message": "Not support runtime param type.",
		})
	}
	return kmsRuntime, nil
}

func (kro *KmsRuntimeOptions) ToDKmsRuntimeOptions() *dkmsutil.RuntimeOptions {
	return &dkmsutil.RuntimeOptions{
		Autoretry:      kro.Autoretry,
		IgnoreSSL:      kro.IgnoreSSL,
		MaxAttempts:    kro.MaxAttempts,
		BackoffPolicy:  kro.BackoffPolicy,
		BackoffPeriod:  kro.BackoffPeriod,
		ReadTimeout:    kro.ReadTimeout,
		ConnectTimeout: kro.ConnectTimeout,
		HttpProxy:      kro.HttpProxy,
		HttpsProxy:     kro.HttpsProxy,
		NoProxy:        kro.NoProxy,
		MaxIdleConns:   kro.MaxIdleConns,
		Socks5Proxy:    kro.Socks5Proxy,
		Socks5NetWork:  kro.Socks5NetWork,
		Verify:         kro.Ca,
	}
}
