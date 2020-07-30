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
package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"link/cinema/api"
	"link/cinema/config"
	"link/cinema/service"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//@Summary Request user to set proxy
//@Description Request user to set proxy to delegate managing item tokens by service
//@Tags user
//@Accept json
//@Produce json
//@Success 200 {object} service.TransferRequestResult "Session token and redirect url to set proxy"
//@Failure 500 {string} string "Internal server error"
//@Router /user/proxy [get]
func (ctr *Controller) RequestProxy(c *gin.Context){
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	proxyReqResult, err := service.RequestProxy(userProfile.UserID, config.GetAPIConfig().ItemContractID)

	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, proxyReqResult)
}


//@Summary Commit a request of setting proxy
//@Description Commit a request of setting proxy
//@Tags user
//@Accept json
//@Produce json
//@Param proxyToken path string true "Proxy session token"
//@Success 200 {string} string "Transaction hash has executed"
//@Failure 500 {string} string "Internal server error"
//@Router /user/proxy/commit/{proxyToken} [get]
func (ctr *Controller) CommitRequestProxy(c *gin.Context) {
	token := c.Param("proxyToken")

	apiResult, err := service.GetProxyStatus(token)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	proxyStatus := make(map[string]string)

	if err := json.Unmarshal(apiResult, &proxyStatus); err != nil {
		c.String(500, err.Error())
		return
	}

	if proxyStatus["status"] != "Authorized" {
		c.String(200, "Failed to request proxy")
		return
	}

	tx, err := service.CommitTransferRequest(token)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, tx.TxHash)
}


//@Summary Login to LINE
//@Description retrieve URL to login through LINE
//@Tags user
//@Accept json
//@Produce json
//@Success 200 {string} string "URL to redirect login page"
//@Failure 500 {string} string "Internal server error"
//@Router /user/login [get]
func (ctr *Controller) LINELogin(c *gin.Context) {
	url := fmt.Sprintf("%s/oauth2/v2.1/authorize", config.GetAPIConfig().LINEAccessEndpoint)

	query := map[string]string{
		"response_type": "code",
		"client_id":     config.GetAPIConfig().ChannelID,
		"redirect_uri":  fmt.Sprintf("%s/api/v0/user/login/callback", config.GetAPIConfig().Endpoint),
		"state":         strconv.FormatInt(rand.Int63(), 16),
		"scope":         "profile%20openid",
	}
	prefix := "?"
	for k, v := range query {
		url += fmt.Sprintf("%s%s=%s", prefix, k, v)
		prefix = "&"
	}

	c.JSON(200, url)
	//c.Redirect(http.StatusMovedPermanently, url)
}

func (ctr *Controller) LINELoginCallback(c *gin.Context) {
	client := http.Client{}

	code := c.Query("code")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", fmt.Sprintf("%s/api/v0/user/login/callback", config.GetAPIConfig().Endpoint))
	data.Add("client_id", config.GetAPIConfig().ChannelID)
	data.Add("client_secret", config.GetAPIConfig().ChannelSecret)

	apiURL := fmt.Sprintf("%s/oauth2/v2.1/token", config.GetAPIConfig().LINEAPIEndpoint)

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)

	apiResult, _ := ioutil.ReadAll(resp.Body)

	type tokenResult struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	token := tokenResult{}

	json.Unmarshal(apiResult, &token)

	session := sessions.Default(c)
	session.Set("accessToken", token.AccessToken)
	session.Set("tokenType", token.TokenType)
	session.Save()

	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
