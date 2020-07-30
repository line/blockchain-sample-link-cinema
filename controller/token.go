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
	"link/cinema/api"
	"link/cinema/config"
	"link/cinema/service"
)

type MovieDiscountBalance struct {
	UserInfo  *service.UserInfo        `json:"userInfo"`
	TokenInfo *service.FungibleBalance `json:"tokenInfo"`
	Txs       []*service.Transaction   `json:"transactions"`
}

//@Summary Get a movie-discount token balance
//@Description Retrieve a movie-discount token balance and summary by user
////@Tags token
//@Accept json
//@Produce json
//@Failure 200 {array} MovieDiscountBalance "Movie-Discount token and summary by user"
//@Failure 500 {string} string "Internal server error"
//@Router /token/balance/movie-discount [get]
func (ctr *Controller) GetMovieDiscountBalance(c *gin.Context) {
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	contractID := config.GetAPIConfig().ItemContractID
	tokenType := config.GetAPIConfig().FungibleTokenType

	userInfo, err := service.GetUserInfo(userProfile.UserID)

	if err != nil {
		c.String(500, err.Error())
		return
	}

	fungibleBalance, err := service.GetFungibleBalance(userProfile.UserID, contractID, tokenType)

	if err != nil {
		c.String(500, err.Error())
		return
	}

	txs, err := service.GetFungibleTransactionHistory(userProfile.UserID, contractID, tokenType)

	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, MovieDiscountBalance{
		UserInfo:  userInfo,
		TokenInfo: fungibleBalance,
		Txs:       txs,
	})
}

type MovieTicketBalance struct {
	Amount int                `json:"amount"`
	Tokens []MovieTicketToken `json:"tokens"`
}

type MovieTicketToken struct {
	Name         string                        `json:"name"`
	TokenID      string                        `json:"tokenId"`
	MovieInfo    service.MovieInfo             `json:"movieInfo"`
	TicketInfo   service.TicketInfo            `json:"ticketInfo"`
	PaymentInfo  service.PaymentInfo           `json:"paymentInfo"`
	Transactions *service.NonFungibleTxHistory `json:"transactions"`
}

//@Summary Get a movie-ticket token balance
//@Description Retrieve movie-ticket token balance and summary using its token index
//@Tags token
//@Accept json
//@Produce json
//@Success 200 {object} MovieTicketBalance "Movie-ticket token balance and summary with provided token index"
//@Failure 500 {string} string "Internal server error"
//@Router /token/balance/movie-ticket [get]
func (ctr *Controller) SearchTicketBalance(c *gin.Context) {
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	contractID := config.GetAPIConfig().ItemContractID
	tokenType := config.GetAPIConfig().NonFungibleTokenType

	nonFungibleInfos, err := service.GetNonFungibleInfo(userProfile.UserID, contractID, tokenType)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	tokens := make([]MovieTicketToken, 0)

	for _, nonFungibleInfo := range nonFungibleInfos {
		tokenIndex := nonFungibleInfo.TokenIndex

		meta := service.NonFungibleMetadata{}
		if err := json.Unmarshal([]byte(nonFungibleInfo.Meta), &meta); err != nil {
			c.String(500, err.Error())
			return
		}

		txs, err := service.GetNonFungibleTransactionHistory(userProfile.UserID, contractID, tokenType, tokenIndex)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		tokens = append(tokens, MovieTicketToken{
			Name:         nonFungibleInfo.Name,
			TokenID:      contractID + tokenType + tokenIndex,
			MovieInfo:    meta.MovieInfo,
			TicketInfo:   meta.TicketInfo,
			PaymentInfo:  meta.PaymentInfo,
			Transactions: txs,
		})

	}

	c.JSON(200, MovieTicketBalance{
		Amount: len(nonFungibleInfos),
		Tokens: tokens,
	})

}

type MovieTokenBalance struct {
	UserInfo  *service.UserInfo            `json:"userInfo"`
	TokenInfo *service.ServiceTokenBalance `json:"tokenInfo"`
	Txs       []*service.Transaction       `json:"transactions"`
}

//@Summary Get a movie token balance
//@Description Retrieve a movie token balance and summary by user
//@Tags token
//@Accept json
//@Produce json
//@Success 200 {object} MovieTokenBalance "Movie token balance and summary by user"
//@Failure 500 {string} string
//@Router /token/balance/movie [get]
func (ctr *Controller) GetMovieTokenBalance(c *gin.Context) {
	contractID := config.GetAPIConfig().ServiceContractID
	userProfile := api.UserProfile{
		UserID: config.GetAPIConfig().UserID,
	}

	userInfo, err := service.GetUserInfo(userProfile.UserID)

	if err != nil {
		c.String(500, err.Error())
	}

	serviceTokenBalance, err := service.GetServiceTokenBalance(userProfile.UserID, contractID)

	if err != nil {
		c.String(500, err.Error())
	}

	txs, err := service.GetServiceTokenTransactionHistory(userProfile.UserID, contractID)

	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, MovieTokenBalance{
		UserInfo:  userInfo,
		TokenInfo: serviceTokenBalance,
		Txs:       txs,
	})
}

type BaseCoinBalance struct {
	UserInfo *service.UserInfo        `json:"userInfo"`
	CoinInfo *service.BaseCoinBalance `json:"coinInfo"`
	Txs      []*service.Transaction   `json:"transactions"`
}

//@Summary Get a base coin balance
//@Description Retrieve a base coin balance and summary by user
//@Tags token
//@Accept json
//@Produce json
//@Success 200 {object} BaseCoinBalance "Base coin balance and summary by user"
//@Failure 500 {string} string
//@Router /token/balance/base-coin [get]
func (ctr *Controller) GetBaseCoinBalance(c *gin.Context) {
	userID := config.GetAPIConfig().UserID
	userProfile := api.UserProfile{
		UserID: userID,
	}

	userInfo, err := service.GetUserInfo(userProfile.UserID)
	if err != nil {
		c.String(500, err.Error())
	}

	baseCoinInfo, err := service.GetBaseCoinBalance(userID)
	if err != nil {
		c.String(500, err.Error())
	}

	txs, err := service.GetBaseCoinTransactionHistory(userID)
	if err != nil {
		c.String(500, err.Error())
	}

	c.JSON(200, BaseCoinBalance{
		UserInfo: userInfo,
		CoinInfo: baseCoinInfo,
		Txs:      txs,
	})
}
