package kucoinapi

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	BIT_BASE_10 = 10
	BIT_SIZE_64 = 64
	//BIT_SIZE_32 = 32
)

type RequestType string

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

var NIL_REQBODY = []byte{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var log = logrus.New()

func SetLogger(logger *logrus.Logger) {
	log = logger
}

var httpTimeout = 100 * time.Second

func SetHttpTimeout(timeout time.Duration) {
	httpTimeout = timeout
}

func GetPointer[T any](v T) *T {
	return &v
}

func HmacSha256(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

type Kucoin struct{}

const (
	KUCOIN_API_SPOT_HTTP   = "api.kucoin.com"
	KUCOIN_API_FUTURE_HTTP = "api-futures.kucoin.com"

	IS_GZIP = true
)

type ApiType int

const (
	SPOT ApiType = iota
	MARGIN
	FUTURE
)

func (apiType *ApiType) String() string {
	switch *apiType {
	case SPOT:
		return "SPOT"
	case MARGIN:
		return "MARGIN"
	case FUTURE:
		return "FUTURE"
	}
	return ""
}

type Client struct {
	ApiKey     string
	ApiSecret  string
	Passphrase string
}
type RestClient struct {
	c *Client
}

type SpotRestClient RestClient
type MarginRestClient RestClient
type FutureRestClient RestClient

func (*Kucoin) NewSpotRestClient(apiKey, apiSecret, passphrase string) *SpotRestClient {
	client := &SpotRestClient{
		&Client{
			ApiKey:     apiKey,
			ApiSecret:  apiSecret,
			Passphrase: passphrase,
		},
	}
	return client
}

func (*Kucoin) NewMarginRestClient(apiKey, apiSecret, passphrase string) *MarginRestClient {
	client := &MarginRestClient{
		&Client{
			ApiKey:     apiKey,
			ApiSecret:  apiSecret,
			Passphrase: passphrase,
		},
	}
	return client
}

func (*Kucoin) NewFutureRestClient(apiKey, apiSecret, passphrase string) *FutureRestClient {
	client := &FutureRestClient{
		&Client{
			ApiKey:     apiKey,
			ApiSecret:  apiSecret,
			Passphrase: passphrase,
		},
	}
	return client
}

var serverTimeDelta int64 = 0

func setServerTimeDelta(delta int64) {
	serverTimeDelta = delta
}

func Request(url string, reqBody []byte, method string, isGzip bool) ([]byte, error) {
	return RequestWithHeader(url, method, map[string]string{}, isGzip)
}

func RequestWithHeader(url string, method string, headerMap map[string]string, isGzip bool) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headerMap {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: httpTimeout,
	}
	if isGzip {
		req.Header.Add("Content-Encoding", "gzip")
		req.Header.Add("Accept-Encoding", "gzip")
	}

	log.Debug(method, ": ", req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	data, err := io.ReadAll(body)
	return data, err
}

func kucoinCallApi[T any](client *Client, url url.URL, reqBody []byte, method string) (*T, error) {
	body, err := Request(url.String(), reqBody, method, IS_GZIP)
	if err != nil {
		return nil, err
	}
	res, err := handlerCommonRest[T](body)
	if err != nil {
		return nil, err
	}
	return &res.Data, res.handlerError()
}

func kucoinCallApiWithSecret[T any](client *Client, url, endpoint, method, body string) (*T, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), BIT_BASE_10)
	signStr := HmacSha256(client.ApiSecret, timestamp+method+endpoint+body)
	respBody, err := RequestWithHeader(url, method, map[string]string{
		"KC-API-SIGN":        base64.StdEncoding.EncodeToString([]byte(signStr)),
		"KC-API-TIMESTAMP":   timestamp,
		"KC-API-KEY":         client.ApiKey,
		"KC-API-PASSPHRASE":  client.Passphrase,
		"KC-API-KEY-VERSION": "2",
	}, IS_GZIP)
	if err != nil {
		return nil, err
	}
	res, err := handlerCommonRest[T](respBody)
	if err != nil {
		return nil, err
	}
	return &res.Data, res.handlerError()
}

func kucoinHandlerRequestApiWithPathQueryParam[T any](apiType ApiType, request *T, name string) url.URL {
	query := kucoinHandlerReq(request)
	u := url.URL{
		Scheme:   "https",
		Host:     kucoinGetRestHostByApiType(apiType),
		Path:     name,
		RawQuery: query,
	}
	return u
}

func kucoinHandlerRequestApiWithSecret[T any](apiType ApiType, request *T, name, secret string) string {
	query := kucoinHandlerReq(request)

	u := url.URL{
		Scheme:   "https",
		Host:     kucoinGetRestHostByApiType(apiType),
		Path:     name,
		RawQuery: query,
	}

	return u.String()
}

func kucoinHandlerRequestApi[T any](apiType ApiType, request *T, name string) string {
	query := kucoinHandlerReq(request)

	u := url.URL{
		Scheme:   "https",
		Host:     kucoinGetRestHostByApiType(apiType),
		Path:     name,
		RawQuery: query,
	}

	return u.String()
}

func kucoinGetRestHostByApiType(apiType ApiType) string {
	switch apiType {
	case SPOT:
		return KUCOIN_API_SPOT_HTTP
	case MARGIN:
		return KUCOIN_API_SPOT_HTTP
	case FUTURE:
		return KUCOIN_API_FUTURE_HTTP
	default:
		return ""
	}
}

// edit needed
func kucoinHandlerReq[T any](req *T) string {
	var paramBuffer bytes.Buffer
	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)
	if v.IsNil() {
		return ""
	}
	t = t.Elem()
	v = v.Elem()
	count := v.NumField()
	for i := 0; i < count; i++ {
		paramName := t.Field(i).Tag.Get("json")
		paramName = strings.ReplaceAll(paramName, ",omitempty", "")
		switch v.Field(i).Elem().Kind() {
		case reflect.String:
			paramBuffer.WriteString(paramName + "=" + v.Field(i).Elem().String() + "&")
		case reflect.Int, reflect.Int64:
			paramBuffer.WriteString(paramName + "=" + strconv.FormatInt(v.Field(i).Elem().Int(), BIT_BASE_10) + "&")
		case reflect.Float32, reflect.Float64:
			paramBuffer.WriteString(paramName + "=" + decimal.NewFromFloat(v.Field(i).Elem().Float()).String() + "&")
		case reflect.Bool:
			paramBuffer.WriteString(paramName + "=" + strconv.FormatBool(v.Field(i).Elem().Bool()) + "&")
		case reflect.Struct:
			sv := reflect.ValueOf(v.Field(i).Interface())
			ToStringMethod := sv.MethodByName("String")
			params := make([]reflect.Value, 0)
			result := ToStringMethod.Call(params)
			paramBuffer.WriteString(paramName + "=" + result[0].String() + "&")
		case reflect.Slice:
			s := v.Field(i).Interface()
			d, _ := json.Marshal(s)
			paramBuffer.WriteString(paramName + "=" + url.QueryEscape(string(d)) + "&")
		case reflect.Invalid:
		default:
			log.Errorf("req type error %s:%s", paramName, v.Field(i).Elem().Kind())
		}
	}
	return strings.TrimRight(paramBuffer.String(), "&")
}

func interfaceStringToFloat64(inter interface{}) float64 {
	return stringToFloat64(inter.(string))
}

func interfaceStringToInt64(inter interface{}) int64 {
	return int64(inter.(float64))
}

func stringToFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, BIT_SIZE_64)
	return f
}

type MySyncMap[K any, V any] struct {
	smap sync.Map
}

func NewMySyncMap[K any, V any]() MySyncMap[K, V] {
	return MySyncMap[K, V]{
		smap: sync.Map{},
	}
}
func (m *MySyncMap[K, V]) Load(k K) (V, bool) {
	v, ok := m.smap.Load(k)

	if ok {
		return v.(V), true
	}
	var resv V
	return resv, false
}
func (m *MySyncMap[K, V]) Store(k K, v V) {
	m.smap.Store(k, v)
}

func (m *MySyncMap[K, V]) Delete(k K) {
	m.smap.Delete(k)
}
func (m *MySyncMap[K, V]) Range(f func(k K, v V) bool) {
	m.smap.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (m *MySyncMap[K, V]) Length() int {
	length := 0
	m.Range(func(k K, v V) bool {
		length += 1
		return true
	})
	return length
}

func (m *MySyncMap[K, V]) MapValues(f func(k K, v V) V) *MySyncMap[K, V] {
	var res = NewMySyncMap[K, V]()
	m.Range(func(k K, v V) bool {
		res.Store(k, f(k, v))
		return true
	})
	return &res
}
