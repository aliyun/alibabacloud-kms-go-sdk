# 压力测试工具

压力测试工具可用来测试KMS实例性能。

## 工具编译
登录Linux ECS实例中， ECS需要处在可以访问KMS实例的VPC中
- 1.下载工具代码

```shell
$ git clone https://github.com/aliyun/alibabacloud-kms-go-sdk.git
```
- 2.切换到项目目录benchmarks下，执行如下命令进行编译

```shell
$ cd alibabacloud-kms-go-sdk/benchmarks
$ go mod init
$ go mod tidy
$ go build -o benchmark .
```

## 使用方法

完成上面编译工具步骤后，在当前目录下将生成工具可执行文件: **benchmark**，执行如下示例命令进行测试，命令行参数参考下面[可配置参数](#可配置参数)项。

运行环境：KMS实例计算性能选项2000，客户端机器配置16核*1台。

示例：使用密钥规格Aliyun_AES_256，进行加密操作(encrypt)压测，数据大小32字节，线程数32，压测时间600秒，命令如下:

```shell
nohup ./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600 -log_path=./log > aes_256_enc.out 2>&1&
```

## 可配置参数

| 参数名称                | 参数说明                                                                                                                                                                                            |
|---------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| case                | 测试case名称，当前支持的测试case名称如下：<br/> encrypt： 加密操作，调用Encrypt接口<br/> decrypt： 解密操作，调用Decrypt接口<br/> sign： 签名操作，调用Sign接口<br/> verify： 验签操作，调用Verify接口<br/> get_secret_value： 获取凭据值操作，调用GetSecretValue接口 |
| endpoint            | KMS实例地址                                                                                                                                                                                         |
| client_key_path     | Client Key文件路径                                                                                                                                                                                  |
| client_key_password | Client Key口令                                                                                                                                                                                    |
| concurrence_nums    | 并发线程数，默认32                                                                                                                                                                                      |
| duration            | 测试时间, 默认600秒                                                                                                                                                                                    |
| period              | 结果输出周期，默认1秒输出一次结果                                                                                                                                                                               |
| log_path            | 日志输出路径，不配置输出到控制台                                                                                                                                                                                |
| key_id              | 测试密钥的Id，加解密测试需要设置此项，获取凭据测试忽略此项                                                                                                                                                                  |
| data_size           | 测试数据大小，单位字节，默认32，测试数据越大性能越低                                                                                                                                                                     |
| secret_name         | 测试的凭据名称，获取凭据测试需要设置此项，加解密测试忽略此项                                                                                                                                                                  |
| ca_path             | CA证书路径，默认为空表示忽略验证服务端证书                                                                                                                                                                          |

测试case使用说明：
- encrypt：测试加密接口性能。

示例：数据大小32字节，线程数32，压测时间600秒，输出到控制台。

```shell
./benchmark -case=encrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- decrypt：测试解密接口性能。

示例：数据大小32字节，线程数32，压测时间600秒，输出到控制台。

```shell
./benchmark -case=decrypt -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- sign：测试签名接口性能。

示例：数据大小32字节，线程数32，压测时间600秒，输出到控制台。

```shell
./benchmark -case=sgin -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- verify：测试验签接口性能。

示例：数据大小32字节，线程数32，压测时间600秒，输出到控制台。

```shell
./benchmark -case=verify -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -key_id=key-**** -data_size=32 -concurrence_nums=32 -duration=600
```
- get_secret_value：测试获取凭据接口性能。

示例：数据大小32字节，线程数32，压测时间600秒，输出到控制台。

```shell
./benchmark -case=get_secret_value -client_key_file=./ClientKey_****.json -client_key_password=**** -endpoint=kst-****.cryptoservice.kms.aliyuncs.com -secret_name=**** -data_size=32 -concurrence_nums=32 -duration=600
```

## KMS实例不同性能的参考配置参数

| 计算性能选项 | 客户端机器配置 | 客户端机器负载(%CPU) |      密钥规格      |  case   | concurrence_nums | data_size |
|:------:|:-------:|:-------------:|:--------------:|:-------:|:----------------:|:---------:|
|  2000  |  8核*1台  |      65       | Aliyun_AES_256 | encrypt |        64        |    32     |
|  3000  |  8核*2台  |      65       | Aliyun_AES_256 | encrypt |      64(单台)      |    32     |
|  4000  | 16核*1台  |      65       | Aliyun_AES_256 | encrypt |        64        |    32     |
|  6000  | 16核*2台  |      65       | Aliyun_AES_256 | encrypt |      64(单台)      |    32     |
|  8000  | 16核*4台  |      65       | Aliyun_AES_256 | encrypt |      64(单台)      |    32     |

配置选择说明：

- 客户端机器负载：客户端机器CPU使用率推荐65%。如果在此负载下压测性能不能符合预期，可以适当增加客户端机器配置。

- 并发数量：客户端为8核机器推荐并发数64。如果在此并发数下压测性能不能符合预期，可以适当增加客户端机器配置，然后提高并发数。

- 数据大小: 数据大小推荐为32。数据越大的压测性能降低。

- 如果业务侧自己通过调用SDK进行压测，出现连接数持续增加，可以将MaxIdleConns设置为与并发数相同。


## 结果输出

```text
INFO simple 2023-04-27 16:19:28.614499 kms_instance_bench.go:250: ----------------- Statistics: [2023-04-27T16:19:28+08:00]--------------
INFO simple 2023-04-27 16:19:28.614511 kms_instance_bench.go:436: [Benchmark-Detail]    RequestCount: 20168     ResponseCount: 20168    TPS: 1617       AvgTPS: 1680
MaxOnceTimeCost: 0.122744 s     MinOnceTimeCost: 0.009237 s     AvgOnceTimeCost: 0.019395 s
ClientErrorCount: 0     LimitExceededErrorCount: 0      TimeoutErrorCount: 0
```
输出参数解释：

RequestCount：总请求数

ResponseCount：总响应数

TPS：每秒处理的事务数

AvgTPS：每秒处理的事务数均值

MaxOnceTimeCost：单次请求最大耗时

MinOnceTimeCost：单次请求最小耗时

AvgOnceTimeCost：单次请求平均耗时

ClientErrorCount：客户端错误次数

LimitExceededErrorCount：限流错误次数

TimeoutErrorCount：超时错误次数

