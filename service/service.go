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
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"link/cinema/api"
	"link/cinema/config"
	"regexp"
	"strconv"
	"strings"
)


var (
	errInvalidParam = errors.New("invalid URL params")
)

func checkUrlParam(params ...string) bool {
	for _, param := range params {
		matched, err := regexp.MatchString("^[a-zA-Z0-9-_=]*$", param)
		if err != nil {
			return true
		}
		if len(param) == 0 || !matched {
			return true
		}
	}
	return false
}

func GetUserInfo(userID string) (*UserInfo, error) {
	if checkUrlParam(userID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s", userID)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	user := &UserInfo{}

	if err := json.Unmarshal(apiResult, &user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetServiceTokenBalance(userID, contractID string) (*ServiceTokenBalance, error) {
	if checkUrlParam(userID, contractID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/service-tokens/%s", userID, contractID)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	serviceTokenBalance := &ServiceTokenBalance{}

	if err := json.Unmarshal(apiResult, serviceTokenBalance); err != nil {
		return nil, err
	}

	return serviceTokenBalance, nil
}

func GetFungibleBalance(userID, contractID, tokenType string) (*FungibleBalance, error) {
	if checkUrlParam(userID, contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/item-tokens/%s/fungibles/%s", userID, contractID, tokenType)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	fungibleBalance := &FungibleBalance{}

	if err := json.Unmarshal(apiResult, fungibleBalance); err != nil {
		return nil, err
	}

	return fungibleBalance, nil
}

func GetNonFungibleInfo(userID, contractID, tokenType string) ([]*NonFungibleInfo, error) {
	if checkUrlParam(userID, contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/item-tokens/%s/non-fungibles/%s", userID, contractID, tokenType)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	nonFungibleInfos := make([]*NonFungibleInfo, 0)

	if err := json.Unmarshal(apiResult, &nonFungibleInfos); err != nil {
		return nil, err
	}

	return nonFungibleInfos, nil

}

func GetBaseCoinBalance(userID string) (*BaseCoinBalance, error) {
	if checkUrlParam(userID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/base-coin", userID)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	baseCoinBalance := &BaseCoinBalance{}

	if err := json.Unmarshal(apiResult, baseCoinBalance); err != nil {
		return nil, err
	}

	return baseCoinBalance, nil
}
func GetTransaction(txHash string) (*Transaction, error) {
	if checkUrlParam(txHash) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/transactions/%s", txHash)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	tx := &Transaction{}
	if err := json.Unmarshal(apiResult, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func GetTransactionHistory(userID, before, after, limit, page, orderBy, msgType string) ([]*Transaction, error) {
	if checkUrlParam(userID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/transactions", userID)

	query := map[string]string{
	}

	if before != "" {
		query["before"] = before
	}

	if after != "" {
		query["after"] = after
	}

	if limit != "" {
		query["limit"] = limit
	}

	if page != "" {
		query["page"] = page
	}

	if orderBy != "" {
		query["orderBy"] = orderBy
	}

	if msgType != "" {
		query["msgType"] = msgType
	}

	apiResult, err := api.CallAPI(path, "GET", query, nil)

	if err != nil {
		return nil, err
	}

	txs := make([]*Transaction, 0)

	if err := json.Unmarshal(apiResult, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}

func GetBaseCoinTransactionHistory(userID string) ([]*Transaction, error) {
	result := make([]*Transaction, 0)
	var (
		txs []*Transaction
		err error
	)
	for page := 1; ; page++ {
		txs, err = GetTransactionHistory(userID, "", "", "", strconv.Itoa(page), "", "link/MsgSend")
		if err != nil {
			return result, err
		}

		if txs == nil || len(txs) == 0 {
			return result, nil
		}

		for _, tx := range txs {
			result = append(result, tx)
			//TODO make limit dynamic
			if len(result) == 5 {
				return result, nil
			}
		}
	}
}

func GetServiceTokenTransactionHistory(userID, contractID string) ([]*Transaction, error) {
	result := make([]*Transaction, 0)
	var (
		txs []*Transaction
		err error
	)
	for page := 1; ; page++ {
		txs, err = GetTransactionHistory(userID, "", "", "", strconv.Itoa(page), "", "token/MsgTransfer")
		if err != nil {
			return result, err
		}

		if txs == nil || len(txs) == 0 {
			return result, nil
		}

		for _, tx := range txs {
			for _, msg := range tx.Tx.Value.Message {
				val := TransferServiceTokenMsg{}
				marshaled, _ := json.Marshal(msg.Value)
				if err := json.Unmarshal(marshaled, &val); err == nil {
					if val.ContractID == contractID {
						result = append(result, tx)
						//TODO make limit dynamic
						if len(result) == 5 {
							return result, nil
						}
					}
				}
			}
		}
	}
}

func GetFungibleTransactionHistory(userID, contractID, tokenType string) ([]*Transaction, error) {
	result := make([]*Transaction, 0)
	var (
		txs []*Transaction
		err error
	)
	for page := 1; ; page++ {
		txs, err = GetTransactionHistory(userID, "", "", "", strconv.Itoa(page), "", "")
		if err != nil {
			return result, err
		}

		if txs == nil || len(txs) == 0 {
			return result, nil
		}

		for _, tx := range txs {
			for _, msg := range tx.Tx.Value.Message {
				if msg.Type == "collection/MsgBurnFT" || msg.Type == "collection/MsgBurnFTFrom" || msg.Type == "collection/MsgMintFT" || msg.Type == "collection/MsgTransferFT" {
					val := FungibleMsg{}
					marshaled, _ := json.Marshal(msg.Value)
					if err := json.Unmarshal(marshaled, &val); err == nil {
						if val.ContractID == contractID {
							hasToken := false
							for _, amt := range val.Amount {
								if strings.HasPrefix(amt.TokenID, tokenType) {
									hasToken = true
									break
								}
							}
							if hasToken {
								result = append(result, tx)
								//TODO make limit dynamic
								if len(result) == 5 {
									return result, nil
								}
							}
						}
					}
				}
			}
		}
	}
}

//TODO store txhash with tokenID as a key in localDB
func GetNonFungibleTransactionHistory(userID, contractID, tokenType, tokenIndex string) (*NonFungibleTxHistory, error) {
	result := &NonFungibleTxHistory{}
	tokenID := contractID + tokenType + tokenIndex
	var (
		txs []*Transaction
		err error
	)
	for page := 1; ; page++ {
		txs, err = GetTransactionHistory(userID, "", "", "", strconv.Itoa(page), "", "collection/MsgMintNFT")
		if err != nil {
			return result, err
		}

		if txs == nil || len(txs) == 0 {
			return result, nil
		}

		for _, tx := range txs {
			innerContractID := ""
			innerTokenID := ""
			for _, log := range tx.Logs {
				for _, event := range log.Events {
					if event.Type == "mint_nft" {
						for _, attr := range event.Attributes {
							if attr.Key == "contract_id" {
								innerContractID = attr.Value
							}
							if attr.Key == "token_id" {
								innerTokenID = attr.Value
							}
						}
					}
				}
			}
			if tokenID == innerContractID+innerTokenID {
				result.MintTransaction = tx
				break
			}
		}
		if result.MintTransaction != nil {
			break
		}
	}

	mintMsg := MintNonFungibleMsg{}
	for _, msg := range result.MintTransaction.Tx.Value.Message {
		marshaledMintMsg, err := json.Marshal(msg.Value)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(marshaledMintMsg, &mintMsg); err != nil {
			return nil, err
		}
	}

	meta := NonFungibleMetadata{}
	if err := json.Unmarshal([]byte(mintMsg.Meta), &meta); err != nil {
		return nil, err
	}

	result.PaymentTransaction, err = GetTransaction(meta.PaymentInfo.PaymentTransaction)
	if err != nil {
		return nil, err
	}

	result.PointTransaction, err = GetTransaction(meta.PaymentInfo.PointTransaction)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func TransferBaseCoin(userID, amount string) (*TransactionAccepted, error) {
	if checkUrlParam(config.GetAPIConfig().WalletAddress) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/wallets/%s/base-coin/transfer", config.GetAPIConfig().WalletAddress)

	params := map[string]interface{}{
		"walletSecret": config.GetAPIConfig().WalletSecret,
		"toUserId":     userID,
		"amount":       amount,
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)

	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil

}

func TransferServiceToken(userID, contractID, amount string) (*TransactionAccepted, error) {
	if checkUrlParam(config.GetAPIConfig().WalletAddress, contractID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/wallets/%s/service-tokens/%s/transfer", config.GetAPIConfig().WalletAddress, contractID)

	params := map[string]interface{}{
		"walletSecret": config.GetAPIConfig().WalletSecret,
		"toUserId":     userID,
		"amount":       amount,
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)

	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil

}

func MintFungible(userID, contractID, tokenType, amount string) (*TransactionAccepted, error) {
	if checkUrlParam(contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/fungibles/%s/mint", contractID, tokenType)

	params := map[string]interface{}{
		"toUserId":     userID,
		"ownerAddress": config.GetAPIConfig().WalletAddress,
		"ownerSecret":  config.GetAPIConfig().WalletSecret,
		"amount":       amount,
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)
	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil

}

func MintNonFungible(userID, contractID, tokenType string, meta NonFungibleMetadata) (*TransactionAccepted, error) {
	if checkUrlParam(contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/mint", contractID, tokenType)

	marshaledMeta, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"toUserId":     userID,
		"name":         "MovieTicket",
		"meta":         string(marshaledMeta),
		"ownerAddress": config.GetAPIConfig().WalletAddress,
		"ownerSecret":  config.GetAPIConfig().WalletSecret,
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)
	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil
}

func BurnFungible(userID, contractID, tokenType, amount string) (*TransactionAccepted, error) {
	if checkUrlParam(contractID, tokenType) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/item-tokens/%s/fungibles/%s/burn", contractID, tokenType)

	params := map[string]interface{}{
		"amount":       amount,
		"fromUserId":   userID,
		"ownerAddress": config.GetAPIConfig().WalletAddress,
		"ownerSecret":  config.GetAPIConfig().WalletSecret,
	}

	apiResult, err := api.CallAPI(path, "POST", nil, params)

	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil
}

func RequestBaseCoinTransfer(userID, amount string) (*TransferRequestResult, error) {
	if checkUrlParam(userID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/base-coin/request-transfer/", userID)

	query := map[string]string{
		"requestType": "redirectUri",
	}

	params := map[string]interface{}{
		"toAddress": config.GetAPIConfig().WalletAddress,
		"amount":    amount,
		//"landingUri": fmt.Sprintf("%s/swagger/index.html", config.GetAPIConfig().Endpoint),
	}

	apiResult, err := api.CallAPI(path, "POST", query, params)

	if err != nil {
		return nil, err
	}

	txReqResult := &TransferRequestResult{}
	if err := json.Unmarshal(apiResult, txReqResult); err != nil {
		return nil, err
	}

	return txReqResult, nil
}

func RequestServiceTransfer(userID, contractID, amount string) (*TransferRequestResult, error) {
	if checkUrlParam(userID, contractID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/service-tokens/%s/request-transfer", userID, contractID)

	query := map[string]string{
		"requestType": "redirectUri",
	}

	params := map[string]interface{}{
		"toAddress": config.GetAPIConfig().WalletAddress,
		"amount":    amount,
		//"landingUri": fmt.Sprintf("%s/swagger/index.html", config.GetAPIConfig().Endpoint),
	}

	apiResult, err := api.CallAPI(path, "POST", query, params)

	if err != nil {
		return nil, err
	}

	txReqResult := &TransferRequestResult{}
	if err := json.Unmarshal(apiResult, txReqResult); err != nil {
		return nil, err
	}

	return txReqResult, nil
}

func RequestProxy(userID, contractID string) (*TransferRequestResult, error) {
	if checkUrlParam(userID, contractID) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/item-tokens/%s/request-proxy", userID, contractID)

	query := map[string]string{
		"requestType": "redirectUri",
	}

	params := map[string]interface{}{
		"ownerAddress": config.GetAPIConfig().WalletAddress,
		//"landingUri":   fmt.Sprintf("%s/swagger/index.html", config.GetAPIConfig().Endpoint),
	}

	apiResult, err := api.CallAPI(path, "POST", query, params)

	if err != nil {
		return nil, err
	}

	txReqResult := &TransferRequestResult{}
	if err := json.Unmarshal(apiResult, txReqResult); err != nil {
		return nil, err
	}

	return txReqResult, nil
}

func GetProxyStatus(token string) ([]byte, error) {
	if checkUrlParam(token) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/user-requests/%s", token)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	return apiResult, nil
}

func GetProxySetting(userID, contractID string) (bool, error) {
	if checkUrlParam(userID, contractID) {
		return false, errInvalidParam
	}
	path := fmt.Sprintf("/v1/users/%s/item-tokens/%s/proxy", userID, contractID)

	apiResult, err := api.CallAPI(path, "GET", nil, nil)

	if err != nil {
		return false, err
	}

	result := make(map[string]bool)

	if err := json.Unmarshal(apiResult, &result); err != nil {
		return false, err
	}

	return result["isApproved"], nil
}

func CommitTransferRequest(token string) (*TransactionAccepted, error) {
	if checkUrlParam(token) {
		return nil, errInvalidParam
	}
	path := fmt.Sprintf("/v1/user-requests/%s/commit", token)

	apiResult, err := api.CallAPI(path, "POST", nil, nil)
	if err != nil {
		return nil, err
	}

	txAccepted := &TransactionAccepted{}

	if err := json.Unmarshal(apiResult, txAccepted); err != nil {
		return nil, err
	}

	return txAccepted, nil

}
