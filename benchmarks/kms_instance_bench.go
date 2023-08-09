package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	dkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	"github.com/aliyun/alibabacloud-kms-go-sdk/sdk"
	"github.com/cihub/seelog"
)

const seelogConfig string = `
<seelog minlevel="debug">
	<outputs formatid="test">
		<buffered size="1048576" flushperiod="100">
			<rollingfile type="size" filename="LOG_DIR" maxsize="2000000000" maxrolls="10" />
		</buffered>
	</outputs>
	<formats>
		<format id="test" format="%LEVEL simple %Date(2006-01-02 15:04:05.000000) %File:%Line: %Msg%n"/>
	</formats>
</seelog>
`
const seelogConsoleConfig string = `
<seelog minlevel="debug">
	<outputs formatid="test">
		<console/>
	</outputs>
	<formats>
		<format id="test" format="%LEVEL simple %Date(2006-01-02 15:04:05.000000) %File:%Line: %Msg%n"/>
	</formats>
</seelog>
`

const recordFormat string = "[Benchmark-Detail]\tRequestCount: %d\tResponseCount: %d\tTPS: %d\tAvgTPS: %d\n" +
	"MaxOnceTimeCost: %f s\tMinOnceTimeCost: %f s\tAvgOnceTimeCost: %f s\n" +
	"ClientErrorCount: %d\tLimitExceededErrorCount: %d\tTimeoutErrorCount: %d"

const (
	TypeEncryptWorker        = 1
	TypeDecryptWorker        = 2
	TypeSignWorker           = 3
	TypeVerifyWorker         = 4
	TypeGetSecretValueWorker = 5
)

var benchmarkCases = map[string]int{
	"encrypt":          TypeEncryptWorker,
	"decrypt":          TypeDecryptWorker,
	"sign":             TypeSignWorker,
	"verify":           TypeVerifyWorker,
	"get_secret_value": TypeGetSecretValueWorker,
}

type Benchmark struct {
	config       *Config
	countCollect *CountCollect
	worker       Worker
	reportLog    seelog.LoggerInterface
	debugLog     seelog.LoggerInterface
}

type Config struct {
	Case            string
	Endpoint        string
	KeyId           string
	ClientKeyPath   string
	P12Password     string
	DataSize        int
	ConcurrenceNums int
	Duration        int64
	Period          int
	LogPath         string
	EnableDebugLog  bool
	CaFilePath      string
	SecretName      string
	Algorithm       string

	plainText         string
	encryptionContext map[string]interface{}
	digest            string
	ca                string
}

type CountRecorder struct {
	count int64
	lock  sync.RWMutex
}

type CountCollect struct {
	RequestCountList            []*CountRecorder
	ResponseCountList           []*CountRecorder
	ResponseLimitCount          *CountRecorder
	TimeoutErrorCount           *CountRecorder
	ClientErrorCount            *CountRecorder
	TimeCostSumPerGoroutineList []time.Duration
	CountPerGoroutineList       []int64
	RunBenchTimeCost            time.Duration
	MinOnceTimeCostList         []time.Duration
	MaxOnceTimeCostList         []time.Duration
	TpsList                     []*CountRecorder
	AnalysisLastTime            time.Time
}

type Worker interface {
	DoAction() (string, error)
}

func NewBenchmark(config *Config, reportLog, debugLog seelog.LoggerInterface) (*Benchmark, error) {
	benchmark := &Benchmark{
		config:       config,
		countCollect: NewCountCollect(config.ConcurrenceNums),
		reportLog:    reportLog,
		debugLog:     debugLog,
	}
	err := benchmark.InitWorker()
	if err != nil {
		return nil, err
	}
	return benchmark, nil
}

func NewCountCollect(concurrenceNums int) *CountCollect {
	return &CountCollect{
		RequestCountList:            make([]*CountRecorder, 0, concurrenceNums),
		ResponseCountList:           make([]*CountRecorder, 0, concurrenceNums),
		ResponseLimitCount:          &CountRecorder{count: 0},
		TimeoutErrorCount:           &CountRecorder{count: 0},
		ClientErrorCount:            &CountRecorder{count: 0},
		TimeCostSumPerGoroutineList: make([]time.Duration, 0, concurrenceNums),
		CountPerGoroutineList:       make([]int64, 0, concurrenceNums),
		MinOnceTimeCostList:         make([]time.Duration, 0, concurrenceNums),
		MaxOnceTimeCostList:         make([]time.Duration, 0, concurrenceNums),
		TpsList:                     make([]*CountRecorder, 0, concurrenceNums),
	}
}

func (bench *Benchmark) InitWorker() error {
	client, err := sdk.NewTransferClient(nil, &dkmsopenapi.Config{
		Protocol:      tea.String("https"),
		ClientKeyFile: tea.String(bench.config.ClientKeyPath),
		Password:      tea.String(bench.config.P12Password),
		Endpoint:      tea.String(bench.config.Endpoint),
		MaxIdleConns:  tea.Int(bench.config.ConcurrenceNums),
	})
	if err != nil {
		return err
	}
	workerType, _ := benchmarkCases[bench.config.Case]
	switch workerType {
	case TypeEncryptWorker:
		bench.worker = NewEncryptWorker(client, bench.config.plainText, bench.config.KeyId, bench.config.encryptionContext, bench.config.ca)
	case TypeDecryptWorker:
		worker := NewEncryptWorker(client, bench.config.plainText, bench.config.KeyId, bench.config.encryptionContext, bench.config.ca)
		res, err := worker.encrypt()
		if err != nil {
			return err
		}
		bench.worker = NewDecryptWorker(client, tea.StringValue(res.Body.CiphertextBlob), bench.config.encryptionContext, bench.config.ca)
	case TypeSignWorker:
		bench.worker = NewSignWorker(client, bench.config.KeyId, bench.config.digest, bench.config.Algorithm, bench.config.ca)
	case TypeVerifyWorker:
		worker := NewSignWorker(client, bench.config.KeyId, bench.config.digest, bench.config.Algorithm, bench.config.ca)
		res, err := worker.sign()
		if err != nil {
			return err
		}
		bench.worker = NewVerifyWorker(client, bench.config.KeyId, tea.StringValue(res.Body.Value), bench.config.digest, bench.config.Algorithm, bench.config.ca)
	case TypeGetSecretValueWorker:
		bench.worker = NewGetSecretValueWorker(client, bench.config.SecretName, bench.config.ca)
	default:
		return errors.New(fmt.Sprintf("invalid benchmark case worker type: %d", workerType))
	}
	return nil
}

func (bench *Benchmark) Run() {
	finishChan := make(chan struct{})
	finishWg := &sync.WaitGroup{}
	syncWg := &sync.WaitGroup{}
	syncWg.Add(bench.config.ConcurrenceNums)
	c := 0
	// 启动压测
	for ; c < bench.config.ConcurrenceNums; c++ {
		finishWg.Add(1)
		bench.countCollect.RequestCountList = append(bench.countCollect.RequestCountList, &CountRecorder{count: 0})
		bench.countCollect.ResponseCountList = append(bench.countCollect.ResponseCountList, &CountRecorder{count: 0})
		bench.countCollect.TimeCostSumPerGoroutineList = append(bench.countCollect.TimeCostSumPerGoroutineList, time.Duration(0))
		bench.countCollect.CountPerGoroutineList = append(bench.countCollect.CountPerGoroutineList, 0)
		bench.countCollect.MinOnceTimeCostList = append(bench.countCollect.MinOnceTimeCostList, time.Duration(math.MaxInt64))
		bench.countCollect.MaxOnceTimeCostList = append(bench.countCollect.MaxOnceTimeCostList, time.Duration(0))
		bench.countCollect.TpsList = append(bench.countCollect.TpsList, &CountRecorder{count: 0})
		go bench.execute(bench.countCollect, syncWg, finishWg, time.Duration(bench.config.Duration)*time.Second, c, finishChan)
		syncWg.Done()
	}

	beginTime := time.Now()
	bench.countCollect.AnalysisLastTime = beginTime
	closeLogChan := make(chan struct{})
	closeLogWg := &sync.WaitGroup{}
	closeLogWg.Add(1)

	// 定时输出统计日志
	go func(countCollect *CountCollect, begin time.Time) {
		i := 0
		for {
			select {
			case <-closeLogChan:
				closeLogWg.Done()
				return
			case <-time.After(time.Duration(bench.config.Period) * time.Second):
				i++
				bench.reportLog.Infof("----------------- Time_%d: [%s]--------------", i, time.Now().Format(time.RFC3339))
				countCollect.RunBenchTimeCost = time.Since(begin)
				period := time.Since(countCollect.AnalysisLastTime).Seconds()
				countCollect.AnalysisLastTime = time.Now()
				countCollect.Analysis(bench.reportLog, period)
				bench.reportLog.Flush()
			}
		}
	}(bench.countCollect, beginTime)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// 最终输出统计日志
	go func(countCollect *CountCollect, begin time.Time) {
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		close(finishChan)
		closeLogWg.Wait()
		fmt.Println()
		bench.reportLog.Infof("----------------- Statistics: [%s]--------------", time.Now().Format(time.RFC3339))
		countCollect.RunBenchTimeCost = time.Since(begin)
		period := time.Since(countCollect.AnalysisLastTime).Seconds()
		countCollect.AnalysisLastTime = time.Now()
		countCollect.Analysis(bench.reportLog, period)
		bench.reportLog.Flush()
		done <- true
	}(bench.countCollect, beginTime)

	// 等待压测完成
	finishWg.Wait()

	// 关闭定时输出日志
	close(closeLogChan)
	// 关闭最终输出日志
	close(sigs)
	<-done
}

func (bench *Benchmark) execute(countCollect *CountCollect, syncWg, finishWg *sync.WaitGroup, duration time.Duration, index int, finishChan <-chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic err: ", err)
		}
		if finishWg != nil {
			finishWg.Done()
		}
		if bench.debugLog != nil {
			bench.debugLog.Flush()
		}
	}()
	if syncWg != nil {
		syncWg.Wait()
	}
	start := time.Now()
	count := int64(0)
	for {
		select {
		case <-finishChan:
			return
		default:
			countCollect.RequestCountList[index].Add()
			onceTimeStart := time.Now()
			requestId, err := bench.worker.DoAction()
			onceTimeCost := time.Since(onceTimeStart)
			// 更新耗时统计
			countCollect.TimeCostSumPerGoroutineList[index] += onceTimeCost
			countCollect.CountPerGoroutineList[index] += 1
			if countCollect.MaxOnceTimeCostList[index].Nanoseconds() < onceTimeCost.Nanoseconds() {
				countCollect.MaxOnceTimeCostList[index] = onceTimeCost
			}
			if countCollect.MinOnceTimeCostList[index].Nanoseconds() > onceTimeCost.Nanoseconds() {
				countCollect.MinOnceTimeCostList[index] = onceTimeCost
			}
			// 更新统计数据
			countCollect.UpdateCount(index, err)
			if err != nil {
				fmt.Println(fmt.Sprintf("[BenchError_%d]\tError：%v", count, err))
			}
			if bench.debugLog != nil {
				bench.debugLog.Infof("[BenchDebug_%d]\tError: %v\tonceTimeCost: %fs\tRequestId: %s", count, err, onceTimeCost.Seconds(), requestId)
			}
			if duration != -1 && time.Since(start) > duration {
				return
			}
			count++
		}
	}
}

func (config *Config) ApplyFlag() error {
	err := config.checkParams()
	if err != nil {
		return err
	}
	if config.plainText == "" || len(config.plainText) != config.DataSize {
		if len(config.plainText) > 0 {
			config.plainText = ""
		}
		for i := 0; i < config.DataSize; i++ {
			config.plainText += "1"
		}
	}
	sha256Array := sha256.Sum256([]byte(config.plainText))
	config.digest = base64.StdEncoding.EncodeToString(sha256Array[:])
	if config.CaFilePath != "" {
		ca, err := ioutil.ReadFile(config.CaFilePath)
		if err != nil {
			return err
		}
		config.ca = string(ca)
	}
	config.encryptionContext = map[string]interface{}{
		"this":       "is",
		"encryption": "context",
	}
	return nil
}

func (config *Config) checkParams() error {
	if config.Case == "" {
		return errors.New(fmt.Sprintf("'--case' can not be empty"))
	} else {
		if _, ok := benchmarkCases[config.Case]; !ok {
			var caseNamesSlice []string
			for name := range benchmarkCases {
				caseNamesSlice = append(caseNamesSlice, name)
			}
			sort.Strings(caseNamesSlice)
			caseNames := strings.Join(caseNamesSlice, "\n")
			return errors.New(fmt.Sprintf("invalid case name:%s\nvalid case names:\n%s", config.Case, caseNames))
		}
	}
	if config.Endpoint == "" {
		return errors.New(fmt.Sprintf("'--endpoint' can not be empty"))
	}
	if config.KeyId == "" {
		if config.Case != "get_secret_value" {
			return errors.New(fmt.Sprintf("'--key_id' can not be empty"))
		}
	}
	if config.ClientKeyPath == "" {
		return errors.New(fmt.Sprintf("'--client_key_path' can not be empty"))
	}
	if config.P12Password == "" {
		return errors.New(fmt.Sprintf("'--client_key_password' can not be empty"))
	}
	if config.SecretName == "" {
		if config.Case == "get_secret_value" {
			return errors.New(fmt.Sprintf("'--secret_name' can not be empty"))
		}
	}
	return nil
}

func (config *Config) GetLogDirPath() (string, error) {
	if _, err := os.Stat(config.LogPath); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(config.LogPath, 0755)
			if err != nil {
				return "", err
			}
		}
	}
	logDirName := fmt.Sprintf("%s_crn%d_%d", config.Case, config.ConcurrenceNums, config.DataSize)
	return path.Join(config.LogPath, logDirName), nil
}

func (collect *CountCollect) Analysis(logger seelog.LoggerInterface, period float64) {
	var allGoroutineTimeCost time.Duration
	var allGoroutineCount, allRequestCount, allResponseCount, allTps int64

	for i := 0; i < len(collect.RequestCountList); i++ {
		allGoroutineCount += collect.CountPerGoroutineList[i]
		allGoroutineTimeCost += collect.TimeCostSumPerGoroutineList[i]
		allRequestCount += collect.RequestCountList[i].Count()
		allResponseCount += collect.ResponseCountList[i].Count()
		allTps += collect.TpsList[i].CountAndReset()
	}

	maxDuration := time.Duration(0)
	minDuration := time.Duration(math.MaxInt64)

	for _, duration := range collect.MinOnceTimeCostList {
		if minDuration.Nanoseconds() > duration.Nanoseconds() {
			minDuration = duration
		}
	}
	for _, duration := range collect.MaxOnceTimeCostList {
		if maxDuration.Nanoseconds() < duration.Nanoseconds() {
			maxDuration = duration
		}
	}
	record := fmt.Sprintf(recordFormat,
		allRequestCount,
		allResponseCount,
		int64(float64(allTps)/period),
		allGoroutineCount/int64(collect.RunBenchTimeCost.Seconds()),
		maxDuration.Seconds(),
		minDuration.Seconds(),
		allGoroutineTimeCost.Seconds()/float64(allGoroutineCount),
		collect.ClientErrorCount.count,
		collect.ResponseLimitCount.count,
		collect.TimeoutErrorCount.count)
	if logger != nil {
		logger.Infof(record)
	}
}

func (collect *CountCollect) UpdateCount(index int, err error) {
	collect.ResponseCountList[index].Add()
	collect.TpsList[index].Add()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			collect.TimeoutErrorCount.Add()
		} else if strings.Contains(err.Error(), "Limit Exceeded") {
			collect.ResponseLimitCount.Add()
		} else {
			collect.ClientErrorCount.Add()
		}
	}
}

func (recorder *CountRecorder) Count() int64 {
	recorder.lock.RLock()
	defer recorder.lock.RUnlock()
	return recorder.count
}

func (recorder *CountRecorder) Add() {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	recorder.count += 1
}

func (recorder *CountRecorder) CountAndReset() int64 {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	temp := recorder.count
	recorder.count = 0
	return temp
}

type EncryptWorker struct {
	client            *sdk.TransferClient
	plainText         string
	keyId             string
	encryptionContext map[string]interface{}

	ca string
}

func NewEncryptWorker(client *sdk.TransferClient, plainText string, keyId string, encryptionContext map[string]interface{}, ca string) *EncryptWorker {
	return &EncryptWorker{
		client:            client,
		plainText:         plainText,
		keyId:             keyId,
		encryptionContext: encryptionContext,
		ca:                ca,
	}
}

func (worker *EncryptWorker) DoAction() (string, error) {
	res, err := worker.encrypt()
	if err != nil {
		return "", err
	}
	return tea.StringValue(res.Body.RequestId), nil
}

func (worker *EncryptWorker) encrypt() (*kms20160120.EncryptResponse, error) {
	request := &kms20160120.EncryptRequest{
		KeyId:             tea.String(worker.keyId),
		Plaintext:         tea.String(worker.plainText),
		EncryptionContext: worker.encryptionContext,
	}
	runtime := &teautil.RuntimeOptions{}
	if worker.ca != "" {
		runtime.Ca = tea.String(worker.ca)
	} else {
		runtime.IgnoreSSL = tea.Bool(true)
	}
	return worker.client.EncryptWithOptions(request, runtime)
}

type DecryptWorker struct {
	client            *sdk.TransferClient
	cipher            string
	encryptionContext map[string]interface{}

	ca string
}

func NewDecryptWorker(client *sdk.TransferClient, cipher string, encryptionContext map[string]interface{}, ca string) *DecryptWorker {
	return &DecryptWorker{
		client:            client,
		cipher:            cipher,
		encryptionContext: encryptionContext,
		ca:                ca,
	}
}

func (worker *DecryptWorker) DoAction() (string, error) {
	res, err := worker.decrypt()
	if err != nil {
		return "", err
	}
	return tea.StringValue(res.Body.RequestId), nil
}

func (worker *DecryptWorker) decrypt() (*kms20160120.DecryptResponse, error) {
	request := &kms20160120.DecryptRequest{
		CiphertextBlob:    tea.String(worker.cipher),
		EncryptionContext: worker.encryptionContext,
	}
	runtime := &teautil.RuntimeOptions{}
	if worker.ca != "" {
		runtime.Ca = tea.String(worker.ca)
	} else {
		runtime.IgnoreSSL = tea.Bool(true)
	}
	return worker.client.DecryptWithOptions(request, runtime)
}

type SignWorker struct {
	client    *sdk.TransferClient
	keyId     string
	digest    string
	algorithm string

	ca string
}

func NewSignWorker(client *sdk.TransferClient, keyId, digest, algorithm, ca string) *SignWorker {
	return &SignWorker{
		client:    client,
		keyId:     keyId,
		digest:    digest,
		algorithm: algorithm,
		ca:        ca,
	}
}

func (worker *SignWorker) DoAction() (string, error) {
	res, err := worker.sign()
	if err != nil {
		return "", err
	}
	return tea.StringValue(res.Body.RequestId), nil
}

func (worker *SignWorker) sign() (*kms20160120.AsymmetricSignResponse, error) {
	request := &kms20160120.AsymmetricSignRequest{
		KeyId:     tea.String(worker.keyId),
		Digest:    tea.String(worker.digest),
		Algorithm: tea.String(worker.algorithm),
	}
	runtime := &teautil.RuntimeOptions{}
	if worker.ca != "" {
		runtime.Ca = tea.String(worker.ca)
	} else {
		runtime.IgnoreSSL = tea.Bool(true)
	}
	return worker.client.AsymmetricSignWithOptions(request, runtime)
}

type VerifyWorker struct {
	client    *sdk.TransferClient
	keyId     string
	value     string
	digest    string
	algorithm string

	ca string
}

func NewVerifyWorker(client *sdk.TransferClient, keyId, value, digest, algorithm, ca string) *VerifyWorker {
	return &VerifyWorker{
		client:    client,
		keyId:     keyId,
		value:     value,
		digest:    digest,
		algorithm: algorithm,
		ca:        ca,
	}
}

func (worker *VerifyWorker) DoAction() (string, error) {
	res, err := worker.verify()
	if err != nil {
		return "", err
	}
	return tea.StringValue(res.Body.RequestId), nil
}

func (worker *VerifyWorker) verify() (*kms20160120.AsymmetricVerifyResponse, error) {
	request := &kms20160120.AsymmetricVerifyRequest{
		KeyId:     tea.String(worker.keyId),
		Value:     tea.String(worker.value),
		Digest:    tea.String(worker.digest),
		Algorithm: tea.String(worker.algorithm),
	}
	runtime := &teautil.RuntimeOptions{}
	if worker.ca != "" {
		runtime.Ca = tea.String(worker.ca)
	} else {
		runtime.IgnoreSSL = tea.Bool(true)
	}
	return worker.client.AsymmetricVerifyWithOptions(request, runtime)
}

type GetSecretValueWorker struct {
	client     *sdk.TransferClient
	secretName string

	ca     string
	logger seelog.LoggerInterface
}

func NewGetSecretValueWorker(client *sdk.TransferClient, secretName, ca string) *GetSecretValueWorker {
	return &GetSecretValueWorker{
		client:     client,
		secretName: secretName,
		ca:         ca,
	}
}

func (worker *GetSecretValueWorker) DoAction() (string, error) {
	res, err := worker.getSecretValue()
	if err != nil {
		return "", err
	}
	return tea.StringValue(res.Body.RequestId), nil
}

func (worker *GetSecretValueWorker) getSecretValue() (*kms20160120.GetSecretValueResponse, error) {
	request := &kms20160120.GetSecretValueRequest{
		SecretName: tea.String(worker.secretName),
	}
	runtime := &teautil.RuntimeOptions{}
	if worker.ca != "" {
		runtime.Ca = tea.String(worker.ca)
	} else {
		runtime.IgnoreSSL = tea.Bool(true)
	}
	return worker.client.GetSecretValueWithOptions(request, runtime)
}

func getLogs(config *Config) (seelog.LoggerInterface, seelog.LoggerInterface, error) {
	var reportLog, debugLog seelog.LoggerInterface
	var resultLogConfig, debugLogConfig string

	if config.LogPath != "" {
		logPath, err := config.GetLogDirPath()
		if err != nil {
			return nil, nil, err
		}
		resultLogConfig = strings.Replace(seelogConfig, "LOG_DIR", logPath+"/statistics.log", 1)
		debugLogConfig = strings.Replace(seelogConfig, "LOG_DIR", logPath+"/debug.log", 1)
	} else {
		resultLogConfig = seelogConsoleConfig
		debugLogConfig = seelogConsoleConfig
	}

	reportLog, err := seelog.LoggerFromConfigAsBytes([]byte(resultLogConfig))
	if err != nil {
		return nil, nil, err
	}

	if config.EnableDebugLog {
		debugLog, err = seelog.LoggerFromConfigAsBytes([]byte(debugLogConfig))
		if err != nil {
			return nil, nil, err
		}
	}
	return reportLog, debugLog, nil
}

func main() {
	exit := func(section string, err error) {
		fmt.Println(fmt.Sprintf("Section[%s]\tError[%v]", section, err))
		os.Exit(1)
	}

	flagCase := flag.String("case", "", "case name, required")
	flagEndpoint := flag.String("endpoint", "", "kms instance address, required")
	flagClientKeyPath := flag.String("client_key_path", "", "client key file path, required")
	flagP12Password := flag.String("client_key_password", "", "client key password, required")
	flagConcurrenceNums := flag.Int("concurrence_nums", 32, "concurrence nums")
	flagDuration := flag.Int64("duration", 600, "benchmark duration time")
	flagPeriod := flag.Int("period", 1, "benchmark result output period")
	flagLogPath := flag.String("log_path", "", "log path")
	flagKeyId := flag.String("key_id", "", "kms cmk id")
	flagDataSize := flag.Int("data_size", 32, "data size")
	flagSecretName := flag.String("secret_name", "", "secret name")
	flagCaPath := flag.String("ca_path", "", "ca cert file path")
	flag.Parse()

	config := &Config{
		Case:            *flagCase,
		Endpoint:        *flagEndpoint,
		KeyId:           *flagKeyId,
		ClientKeyPath:   *flagClientKeyPath,
		P12Password:     *flagP12Password,
		DataSize:        *flagDataSize,
		ConcurrenceNums: *flagConcurrenceNums,
		Duration:        *flagDuration,
		Period:          *flagPeriod,
		LogPath:         *flagLogPath,
		SecretName:      *flagSecretName,
		CaFilePath:      *flagCaPath,
	}
	err := config.ApplyFlag()
	if err != nil {
		exit("ApplyFlag", err)
	}

	fmt.Println(fmt.Sprintf("Start [%s] Case", *flagCase))

	reportLog, debugLog, err := getLogs(config)
	if err != nil {
		exit("getLogs", err)
	}

	defer func() {
		if reportLog != nil {
			reportLog.Flush()
		}
		if debugLog != nil {
			debugLog.Flush()
		}
	}()

	benchmark, err := NewBenchmark(config, reportLog, debugLog)
	if err != nil {
		exit("NewBenchmark", err)
	}
	benchmark.Run()

	fmt.Println(fmt.Sprintf("[%s] Case complete!", config.Case))
}
