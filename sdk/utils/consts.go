package utils

const (
	GcmIvLength                            = 12
	EktIdLength                            = 36
	EncryptApiName                         = "Encrypt"
	AsymmetricEncryptApiName               = "AsymmetricEncrypt"
	DecryptApiName                         = "Decrypt"
	AsymmetricDecryptApiName               = "AsymmetricDecrypt"
	AsymmetricSignApiName                  = "AsymmetricSign"
	AsymmetricVerifyApiName                = "AsymmetricVerify"
	GenerateDataKeyApiName                 = "GenerateDataKey"
	GenerateDataKeyWithoutPlaintextApiName = "GenerateDataKeyWithoutPlaintext"
	GetPublicKeyApiName                    = "GetPublicKey"
	GetSecretValueApiName                  = "GetSecretValue"
	DigestMessageType                      = "DIGEST"
	KMSKeySpecAES256                       = "AES_256"
	KMSKeySpecAES128                       = "AES_128"
	MigrationKeyVersionIdKey               = "x-kms-migrationkeyversionid"
	NumberOfBytesAes256                    = 32
	NumberOfBytesAes128                    = 16
	MagicNum                               = '$'
	MagicNumLength                         = 1
	CipherVerAndPaddingModeLength          = 1
	AlgorithmLength                        = 1
	CipherVer                              = 0
	AlgAesGcm                              = 2
)
