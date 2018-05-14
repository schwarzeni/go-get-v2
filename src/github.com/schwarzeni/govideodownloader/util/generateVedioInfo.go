package util

import "github.com/schwarzeni/govideodownloader/model"

func GenerateVideoInfo(jsonApiUrl string, header map[string]string) (model.Video, error) {
	urlsets, err := ParseJsonApi(jsonApiUrl, header)
	if err != nil {
		return model.Video{}, err
	}
	return model.Video{
		Cookie:       header["Cookie"],
		ApiUrl:       jsonApiUrl,
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
		VedioSection: urlsets}, nil
}

func GenerateVideoInfoForWeb(jsonApiUrl string, webpageUrl string, header map[string]string) (model.Video, error) {
	urlsets, err := ParseJsonApi(jsonApiUrl, header)
	if err != nil {
		return model.Video{}, err
	}
	return model.Video{
		Cookie:       header["Cookie"],
		ApiUrl:       jsonApiUrl,
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
		VedioSection: urlsets,
		WebPageUrl:   webpageUrl}, nil
}
