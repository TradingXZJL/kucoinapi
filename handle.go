package kucoinapi

import "fmt"

type KucoinErrorRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type KucoinRestRes[T any] struct {
	KucoinErrorRes
	Data T `json:"data"`
}

func handlerCommonRest[T any](data []byte) (*KucoinRestRes[T], error) {
	res := &KucoinRestRes[T]{}
	log.Info(string(data))
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Error("rest返回值获取失败", err)
	}
	return res, err
}

func (err *KucoinErrorRes) handlerError() error {
	if err.Code != "0" && err.Code != "200000" {
		return fmt.Errorf("request error:[code:%v][message:%v]", err.Code, err.Msg)
	}

	return nil
}
