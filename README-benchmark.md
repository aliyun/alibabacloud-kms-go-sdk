# Pressure Measurement Tools

The pressure measurement tools could be used to test KMS instance performance.

## Compile

Login to a Linux ECS instance, the ECS can access access the KMS instance. 

- 1.Download code

```shell
$ git clone https://github.com/aliyun/alibabacloud-kms-go-sdk.git
```
- 2.Go to the project directory named benchmarks,and execute the following command to compile

```shell
$ cd alibabacloud-kms-go-sdk/benchmarks
$ go mod init
$ go mod tidy
$ go build -o benchmark .
```

## Usage

After the above compilation steps, an executable file will be generated in the current directory named **benchmark**, and execute the following command to test, command line parameters refer to the following [configurable parameters](#configurable parameters).

Runtime environment: KMS instance computing performance option is 2000, and the configuration of the pressure measurement clients may be 16 cores*1.

Example: Use the key specification Aliyun_AES_256 to perform encryption operation (encrypt) pressure test, the data size is 32 bytes, the number of threads is 32, and the pressure test duration is 600 seconds. The command is as follows:

```shell
nohup ./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600 -log_path=./log > aes_256_enc.out 2>&1&
```

## Configurable parameters

| Name                | Description                                                                                                                                                                                                                                                                                                                                                                |
|---------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| case                | Test case name, currently supported test case names are as follows: <br/> encrypt: encryption operation, call the Encrypt <br/> decrypt: decryption operation, call the Decrypt <br/> sign: signature operation, call the Sign <br/> verify: verification operation, call the Verify <br/> get_secret_value: to obtain the secret value operation, call the GetSecretValue |
| endpoint            | KMS instance endpoint                                                                                                                                                                                                                                                                                                                                                      |
| client_key_path     | the file path of Client Key                                                                                                                                                                                                                                                                                                                                                     |
| client_key_password | the password of Client Key                                                                                                                                                                                                                                                                                                                                                         |
| concurrence_nums    | concurrence goroutines number, default 32                                                                                                                                                                                                                                                                                                                                               |
| duration            | testing duration,default 600s                                                                                                                                                                                                                                                                                                                                                  |
| period              | The result output period, default 1 per second                                                                                                                                                                                                                                                                                                                             |
| log_path            | Log output path, if this parameter is not configured, output to the console                                                                                                                                                                                                                                                                                                |
| key_id              | The kms cmk id, which needs to be set for the encrypt and decrypt test, and ignored for the get secret value test                                                                                                                                                                                                                                                          |
| data_size           | Test data size, unit byte, default 32                                                                                                                                                                                                                                                                                                                                      |
| secret_name         | Secret name, which needs to be set for get secret value test, and ignored for the encrypt and decrypt test                                                                                                                                                                                                                                                                 |
| ca_path             | the path of CA certificate , Ignore verifying the server certificate by default                                                                                                                                                                                                                                                                                       |

Test case instructions:

- encrypt: test Encrypt api performance.

Example: The data size is 32 bytes, the number of threads is 32, the pressure test time is 600 seconds, and the output is output to the console.

```shell
./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- decrypt: test Decrypt api performance.

Example: The data size is 32 bytes, the number of threads is 32, the pressure test time is 600 seconds, and the output is output to the console.

```shell
./benchmark -case=decrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- sign：test Sign api performance.

Example: The data size is 32 bytes, the number of threads is 32, the pressure test time is 600 seconds, and the output is output to the console.

```shell
./benchmark -case=sgin -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- verify：test Verify api performance.

Example: The data size is 32 bytes, the number of threads is 32, the pressure test time is 600 seconds, and the output is output to the console.

```shell
./benchmark -case=verify -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- get_secret_value：test GetSecretValue api performance.

Example: The data size is 32 bytes, the number of threads is 32, the pressure test time is 600 seconds, and the output is output to the console.

```shell
./benchmark -case=get_secret_value -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -secret_name=**** -data_size=32 -concurrence_nums=32 -duration=600
```

## Reference configuration parameters for different performances of KMS instances

| Computing Performance Options | Client Machine Configuration | Client Machine Load(%CPU) | Key Specification |  case   | concurrence_nums | data_size |
|:-----------------------------:|:----------------------------:|:-------------------------:|:-----------------:|:-------:|:----------------:|:---------:|
|             2000              |         8 cores * 1          |            65             |  Aliyun_AES_256   | encrypt |        64        |    32     |
|             3000              |         8 cores * 2          |            65             |  Aliyun_AES_256   | encrypt |    64(Single)    |    32     |
|             4000              |         16 cores * 1         |            65             |  Aliyun_AES_256   | encrypt |        64        |    32     |
|             6000              |         16 cores * 2         |            65             |  Aliyun_AES_256   | encrypt |    64(Single)    |    32     |
|             8000              |         16 cores * 4         |            65             |  Aliyun_AES_256   | encrypt |    64(Single)    |    32     |

Description of configuration options:

- Client machine load: The CPU usage of the client machine is recommended to be 65%.

- Concurrency nums: The recommended 64 for 8-core machines. If the performance of the pressure test cannot meet expectations under this concurrency, you can appropriately increase the client machine configuration, and then increase the concurrency.

- Data size: The recommended data size is 32. The larger the data, the lower the pressure measurement performance.

- If the business side calls the SDK for stress testing and the number of connections continues to increase, you can set MaxIdleConns to be the same as the number of concurrent connections.

## Result output

```text
INFO simple 2023-04-27 16:19:28.614499 kms_instance_bench.go:250: ----------------- Statistics: [2023-04-27T16:19:28+08:00]--------------
INFO simple 2023-04-27 16:19:28.614511 kms_instance_bench.go:436: [Benchmark-Detail]    RequestCount: 20168     ResponseCount: 20168    TPS: 1617       AvgTPS: 1680
MaxOnceTimeCost: 0.122744 s     MinOnceTimeCost: 0.009237 s     AvgOnceTimeCost: 0.019395 s
ClientErrorCount: 0     LimitExceededErrorCount: 0      TimeoutErrorCount: 0
```
Explanation of output parameters:

RequestCount: total number of requests

ResponseCount: total number of responses

TPS：transactions per second

AvgTPS：average transactions per second

MaxOnceTimeCost：maximum time spent on a single request

MinOnceTimeCost：minimum time spent on a single request

AvgOnceTimeCost：average time spent on a single request

ClientErrorCount：client error count

LimitExceededErrorCount：limit exceeded error count

TimeoutErrorCount：timeout error count

