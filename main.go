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
package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"link/cinema/config"
	"link/cinema/controller"
	"link/cinema/docs"
	"os"
	"strings"
)

// @title Link Cinema API
// @version 0.1
// @description This is sample dapp to provide trials of LBD service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v0
func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	if configPath := os.Getenv(config.Path); configPath != "" {
		config.LoadAPIConfig(configPath)
	}
	host := config.GetAPIConfig().Endpoint
	if strings.HasPrefix(host, "http://") {
		host = host[7:]
	}
	if strings.HasPrefix(host, "https://") {
		host = host[8:]
	}
	docs.SwaggerInfo.Host = host

	ctr := controller.NewController()

	v0 := r.Group("/api/v0")
	{
		user := v0.Group("/user")
		{
			//user.GET("/login", ctr.LINELogin)
			//user.GET("/login/callback", ctr.LINELoginCallback)
			user.GET("/proxy", ctr.RequestProxy)
			user.GET("/proxy/commit/:proxyToken", ctr.CommitRequestProxy)
		}

		ticket := v0.Group("/ticket")
		{
			ticket.GET("/", ctr.GetPurchaseInfo)
			ticket.POST("/purchase", ctr.RequestTicketPurchasing)
			ticket.POST("/purchase/extra", ctr.RequestExtraPurchase)
			ticket.POST("/purchase/commit/:baseCoinTransferToken/:movieTokenTransferToken", ctr.CommitPurchasingTicket)
		}

		token := v0.Group("/token")
		{
			token.GET("/balance/base-coin", ctr.GetBaseCoinBalance)
			token.GET("/balance/movie-discount", ctr.GetMovieDiscountBalance)
			token.GET("/balance/movie-ticket", ctr.SearchTicketBalance)
			token.GET("/balance/movie", ctr.GetMovieTokenBalance)
		}
		test := v0.Group("/test")
		{
			test.GET("/init", ctr.InitUser)

			test.GET("/transaction", ctr.GetTransaction)
			test.GET("/config", ctr.ShowConfig)
		}
	}

	url := ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", config.GetAPIConfig().Endpoint)) // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}