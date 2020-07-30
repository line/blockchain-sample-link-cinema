/*
Copyright 2020 LINE Corporation

LINE Corporation licenses this file to you under the Apache License,
version 2.0 (the "License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at:

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations
under the License
*/
package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"io"
	"io/ioutil"
	"link/cinema/config"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func CallAPI(path, method string, query map[string]string, params map[string]interface{}) ([]byte, error) {
	client := http.Client{}
	var body io.Reader
	queryStr := ""
	if method == "POST" {
		jsonParams, _ := json.Marshal(params)
		body = bytes.NewReader(jsonParams)
	}
	if query != nil {
		prefix := "?"
		for k, v := range query {
			queryStr += fmt.Sprintf("%s%s=%s", prefix, k, v)
			prefix = "&"
		}
	}

	req, err := http.NewRequest(method, config.GetAPIConfig().LBDAPIEndpoint+path+queryStr, body)
	if err != nil {
		return nil, err
	}

	timestamp, err := GetServerTime()

	if err != nil {
		return nil, err
	}

	nonce := makeNonce(8)

	sig := getSignature(nonce, timestamp, method, path, queryStr, params)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("service-api-key", config.GetAPIConfig().APIKey)
	req.Header.Add("signature", sig)
	req.Header.Add("nonce", nonce)
	req.Header.Add("timestamp", timestamp)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	apiResult, err := ioutil.ReadAll(resp.Body)

	type response struct {
		ResponseTime  uint64      `json:"responseTime"`
		StatusCode    int         `json:"statusCode"`
		StatusMessage string      `json:"statusMessage"`
		ResponseData  interface{} `json:"responseData"`
	}

	unmarshalResult := response{}
	err = json.Unmarshal(apiResult, &unmarshalResult)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("invalid API response, path: %s, response: %s(%d %s)", path, resp.Status, unmarshalResult.StatusCode, unmarshalResult.StatusMessage)
	}

	if err != nil {
		return nil, err
	}

	if unmarshalResult.StatusCode >= 1000 && unmarshalResult.StatusCode <= 1999 {
		return json.Marshal(unmarshalResult.ResponseData)
	}
	return nil, errors.New(fmt.Sprintf("%d: %s", unmarshalResult.StatusCode, unmarshalResult.StatusMessage))
}

func makeNonce(length int) string {
	result := make([]byte, 0)
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvxwyz0123456789"
	for i := 0; i < length; i++ {
		n := rand.Intn(len(charset))
		result = append(result, charset[n])
	}

	return string(result)
}

func getSignature(nonce, timestamp, method, path string, query string, params map[string]interface{}) string {
	msg := nonce + timestamp + method + path + query
	prefix := "?"
	if len(query) > 0 {
		prefix = "&"
	}

	if params != nil {
		paramMap := make(map[string]string)
		parseParams(paramMap, "", params)

		sortable := make([]string, 0)

		for k, _ := range paramMap {
			sortable = append(sortable, k)
		}

		sort.Slice(sortable, func(i, j int) bool {
			return strings.Compare(sortable[i], sortable[j]) < 0
		})

		for _, k := range sortable {
			msg += fmt.Sprintf("%s%s=%s", prefix, k, paramMap[k])
			prefix = "&"
		}
	}

	hash := hmac.New(sha512.New, []byte(config.GetAPIConfig().APISecret))
	hash.Write([]byte(msg))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func parseParams(result map[string]string, key string, params interface{}) {
	if mp, ok := params.(map[string]interface{}); ok {
		for k, v := range mp {
			newKey := fmt.Sprintf("%s.%s", key, k)
			if len(key) == 0 {
				newKey = k
			}
			parseParams(result, newKey, v)
		}
		return
	}

	result[key] = fmt.Sprint(params)
}

func GetServerTime() (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", config.GetAPIConfig().LBDAPIEndpoint+"/v1/time", nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("service-api-key", config.GetAPIConfig().APIKey)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	type timeResult struct {
		ResponseTime uint64 `json:"responseTime"`
	}

	result := timeResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return strconv.FormatUint(result.ResponseTime, 10), nil

}

func GetUserProfileFromSession(session sessions.Session) (*UserProfile, error) {
	profile := &UserProfile{}

	accessTokenI := session.Get("accessToken")

	if accessTokenI == nil {
		return nil, errors.New("accessToken not found")
	}

	accessToken, ok := accessTokenI.(string)

	if !ok {
		return nil, errors.New("invalid accessToken")
	}

	tokenTypeI := session.Get("tokenType")

	if tokenTypeI == nil {
		return nil, errors.New("tokenType not found")
	}

	tokenType, ok := tokenTypeI.(string)

	if !ok {
		return nil, errors.New("invalid tokenType")
	}

	client := http.Client{}

	url := fmt.Sprintf("%s/v2/profile", config.GetAPIConfig().LINEAPIEndpoint)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, profile); err != nil {
		return nil, err
	}

	return profile, nil

}
