package util

import (
	"encoding/json"

	"io/ioutil"

	"github.com/schwarzeni/govideodownloader/model"
)

func ParseJsonApi(url string, header map[string]string) ([]model.SingleVedioApi, error) {
	resp, err := MethodGet(url, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//var jsonBodyResult model.VedioApi
	var jsonBodyRawResult map[string]*json.RawMessage
	var jsonBodyResult []model.SingleVedioApi
	json.Unmarshal(jsonBody, &jsonBodyRawResult)
	json.Unmarshal(*jsonBodyRawResult["durl"], &jsonBodyResult)

	return jsonBodyResult, nil
}
