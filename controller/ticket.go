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
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"link/cinema/api"
	"link/cinema/config"
	"link/cinema/service"
	"math/big"
	"strconv"
	"time"
)


var (
	movieTokenNotUsed = "0"
)

func checkPrice(info service.PriceInfo) bool{
	ticketPrice := service.DefaultTicket.Price

	if info.SubTotal != ticketPrice {
		return false
	}

	if info.UsedFungible < 0 || info.UsedFungible > 1 {
		return false
	}

	if info.UsedServiceToken < 0 || info.UsedServiceToken > 1000 || info.UsedServiceToken % 1000 != 0{
		return false
	}

	discount := info.UsedServiceToken / 1000 + info.UsedFungible * 5

	if -discount != info.Discount {
		return false
	}

	if ticketPrice - discount != info.GrandTotal {
		return false
	}

	return true
}

//@Summary Get a purchase info
//@Description Retrieve a purchase info about given movie ticket
//@Tags ticket
//@Accept json
//@Produce json
//@Success 200 {object} service.PurchaseInfo "Ticket info"
//@Failure 500 {string} string "Internal server error"
//@Router /ticket [get]
func (ctr *Controller) GetPurchaseInfo(c *gin.Context) {
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	serviceContractID := config.GetAPIConfig().ServiceContractID
	itemContractID := config.GetAPIConfig().ItemContractID
	tokenType := config.GetAPIConfig().FungibleTokenType

	resp := service.PurchaseInfo{
		MovieInfo:  service.DefaultMovie,
		TicketInfo: service.DefaultTicket,
		PriceInfo:  service.PriceInfo{},
	}

	fungibleBalance, err := service.GetFungibleBalance(userProfile.UserID, itemContractID, tokenType)

	discount := 0

	if err != nil {
		c.String(500, err.Error())
		return
	}

	fungibleAmt, err := strconv.Atoi(fungibleBalance.Amount)

	if err != nil {
		c.String(500, err.Error())
		return
	}

	//TODO make fungible ratio dynamic
	fungibleRatio := 5

	if fungibleAmt > 1 {
		fungibleAmt = 1
	}
	discount -= fungibleAmt * fungibleRatio

	serviceTokenBalance, err := service.GetServiceTokenBalance(userProfile.UserID, serviceContractID)

	if err != nil {
		c.String(500, err.Error())
		return
	}
	//TODO make service ratio dynamic
	serviceRatio := big.NewInt(1000)

	serviceAmt, ok := new(big.Int).SetString(serviceTokenBalance.Amount, 10)
	if !ok {
		c.String(500, "Invalid movie token amount")
	}

	for i := 0; i < serviceTokenBalance.Decimals; i++ {
		serviceAmt.Div(serviceAmt, big.NewInt(10))
	}

	serviceAmt.Div(serviceAmt, serviceRatio)
	if serviceAmt.Cmp(big.NewInt(1)) == 1 {
		serviceAmt = big.NewInt(1)
	}

	discount -= int(serviceAmt.Int64())

	resp.PriceInfo = service.PriceInfo{
		UsedFungible:     fungibleAmt,
		UsedServiceToken: int(serviceAmt.Mul(serviceAmt, serviceRatio).Int64()),
		SubTotal:         resp.TicketInfo.Price,
		Discount:         discount,
		GrandTotal:       resp.TicketInfo.Price + discount,
	}

	c.JSON(200, resp)
}

//@Summary Request user to purchase
//@Description Request user to transfer token at LBW
//@Tags ticket
//@Accept json
//@Produce json
//@Param purchase_info body service.PurchaseInfo true "Purchase info"
//@Success 200 {object} service.TransferRequestResult "Session token and redirect url to transfer token"
//@Failure 500 {string} string "Internal server error"
//@Router /ticket/purchase [post]
func (ctr *Controller) RequestTicketPurchasing(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	purchaseInfo := &service.PurchaseInfo{}

	if err := json.Unmarshal(reqBody, purchaseInfo); err != nil {
		c.String(500, err.Error())
		return
	}

	if !checkPrice(purchaseInfo.PriceInfo) {
		c.String(500, "Invalid price info")
		return
	}

	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	if purchaseInfo.PriceInfo.UsedFungible > 0 {
		isApproved, err := service.GetProxySetting(userProfile.UserID, config.GetAPIConfig().ItemContractID)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		if !isApproved {
			c.String(500, "Cannot transfer movie-discount token without proxy setting")
			return
		}
	}

	amt := big.NewInt(int64(purchaseInfo.PriceInfo.GrandTotal))
	amt.Mul(amt, big.NewInt(1000000))

	reqResult, err := service.RequestBaseCoinTransfer(userProfile.UserID, amt.String())

	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, reqResult)
	//c.Redirect(http.StatusMovedPermanently, resp.RedirectURI)
}

//@Summary Request user to purchase extra token
//@Description Request user to transfer movie-token used for discounting ticket price
//@Tags ticket
//@Accept json
//@Produce json
//@Param purchase_info body service.PurchaseInfo true "Purchase info"
//@Success 200 {object} service.TransferRequestResult "Session token and redirect url to transfer a token"
//@Failure 500 {string} string "Internal server error"
//@Router /ticket/purchase/extra [post]
func (ctr *Controller) RequestExtraPurchase(c *gin.Context) {
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
	}

	purchaseInfo := &service.PurchaseInfo{}

	if err := json.Unmarshal(reqBody, purchaseInfo); err != nil {
		c.String(500, err.Error())
		return
	}

	if !checkPrice(purchaseInfo.PriceInfo) {
		c.String(500, "Invalid price info")
		return
	}

	if purchaseInfo.PriceInfo.UsedServiceToken > 0 {

		amount := big.NewInt(int64(purchaseInfo.PriceInfo.UsedServiceToken))
		amount.Mul(amount, big.NewInt(1000000))
		txReqResult, err := service.RequestServiceTransfer(userProfile.UserID, config.GetAPIConfig().ServiceContractID, amount.String())

		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.JSON(200, txReqResult)
		return
	}

	c.JSON(200, service.TransferRequestResult{
		RequestSessionToken: movieTokenNotUsed,
		RedirectURI:         "",
	})
}

//@Summary Commit a purchasing movie-ticket token
//@Description Commit transactions to purchase movie-ticket token and mint a movie-ticket token to user wallet
//@Tags ticket
//@Accept json
//@Produce json
//@Param purchase_info body service.PurchaseInfo true "Purchase info"
//@Param baseCoinTransferToken path string true "Base coin transfer session Token"
//@Param movieTokenTransferToken path string true "Base coin transfer session Token"
//@Success 200 {array} string "Transaction hashes has executed"
//@Failure 500 {string} string "Internal server error"
//@Router /ticket/purchase/commit/{baseCoinTransferToken}/{:movieTokenTransferToken} [post]
func (ctr *Controller) CommitPurchasingTicket(c *gin.Context) {

	resp := make([]string, 0)
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	serviceContractID := config.GetAPIConfig().ServiceContractID
	itemContractID := config.GetAPIConfig().ItemContractID
	fungibleTokenType := config.GetAPIConfig().FungibleTokenType
	nonFungibleTokenType := config.GetAPIConfig().NonFungibleTokenType

	baseSessionToken := c.Param("baseCoinTransferToken")
	serviceSessionToken := c.Param("movieTokenTransferToken")

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	purchaseInfo := service.PurchaseInfo{}
	if err := json.Unmarshal(reqBody, &purchaseInfo); err != nil {
		c.String(500, err.Error())
		return
	}

	if !checkPrice(purchaseInfo.PriceInfo) {
		c.String(500, "Invalid price info")
		return
	}

	if fungibleAmt := purchaseInfo.PriceInfo.UsedFungible; fungibleAmt > 0 {
		tx, err := service.BurnFungible(userProfile.UserID, itemContractID, fungibleTokenType, strconv.Itoa(fungibleAmt))
		if err != nil {
			c.String(500, err.Error())
			return
		}
		resp = append(resp, tx.TxHash)
	}

	serviceAmt := new(big.Int).Mul(big.NewInt(int64(purchaseInfo.PriceInfo.GrandTotal)), big.NewInt(1000000))
	serviceAmt.Div(serviceAmt, big.NewInt(10))
	serviceAmt.Mul(serviceAmt, big.NewInt(1000))
	serviceTx, err := service.TransferServiceToken(userProfile.UserID, serviceContractID, serviceAmt.String())
	if err != nil {
		c.String(500, err.Error())
		return
	}

	if serviceSessionToken != movieTokenNotUsed {
		tx, err := service.CommitTransferRequest(serviceSessionToken)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		resp = append(resp, tx.TxHash)
	}

	baseTx, err := service.CommitTransferRequest(baseSessionToken)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	resp = append(resp, baseTx.TxHash)

	meta := service.NonFungibleMetadata{
		MovieInfo:  purchaseInfo.MovieInfo,
		TicketInfo: purchaseInfo.TicketInfo,
		PaymentInfo: service.PaymentInfo{
			PaymentDate:        time.Now(),
			PaymentTransaction: baseTx.TxHash,
			PointTransaction:   serviceTx.TxHash,
		},
	}

	tx, err := service.MintNonFungible(userProfile.UserID, itemContractID, nonFungibleTokenType, meta)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	resp = append(resp, tx.TxHash)

	c.JSON(200, resp)

}
