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
	"github.com/gin-gonic/gin"
	"link/cinema/api"
	"link/cinema/config"
	"link/cinema/service"
)

//@Summary Get a transaction
//@Description Retrieve a Transaction using its hash
//@Tags test
//@Accept json
//@Produce json
//@Param txhash query string true "Transaction hash used for searching"
//@Success 200 {object} service.Transaction "Transaction with the provided hash"
//@Router /test/transaction [get]
func (ctr *Controller) GetTransaction(c *gin.Context) {
	txHash := c.Query("txhash")
	tx, _ := service.GetTransaction(txHash)

	c.JSON(200, tx)
}

//@Summary Init asset for test user
//@Description Transfer tokens to user
//@Tags test
//@Accept json
//@Produce json
//@Success 200 {array} string "transaction hashes has executed"
//@Failure 500 {string} string "Internal server error"
//@Router /test/init [get]
func (ctr *Controller) InitUser(c *gin.Context) {
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	txs := make([]string, 0)

	cfg := config.GetAPIConfig()

	tx, err := service.TransferBaseCoin(userProfile.UserID, "100000000")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	txs = append(txs, tx.TxHash)

	tx, err = service.TransferServiceToken(userProfile.UserID, cfg.ServiceContractID, "10000000000")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	txs = append(txs, tx.TxHash)

	tx, err = service.MintFungible(userProfile.UserID, cfg.ItemContractID, cfg.FungibleTokenType, "10")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	txs = append(txs, tx.TxHash)

	c.JSON(200, txs)
}

//@Summary Show a config
//@Description Show a config
//@Tags test
//@Accept json
//@Produce json
//@Success 200 {object} config.APIConfig "Server Configuration"
//@Router /test/config [get]
func (ctr *Controller) ShowConfig(c *gin.Context) {
	c.JSON(200, config.GetAPIConfig())
}
