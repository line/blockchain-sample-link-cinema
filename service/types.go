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
	"time"
)

var (
	DefaultMovie = MovieInfo{
		Title:       "The LINK Movie",
		Score:       4.9,
		Country:     "South Korea",
		RunningTime: 132,
		Genre:       "Drama",
		Year:        2020,
	}
	loc, _        = time.LoadLocation("Local")
	DefaultTicket = TicketInfo{
		Date:    time.Date(2020, 2, 1, 3, 30, 0, 0, loc),
		Theater: "Seoul, World Tower, Theater 1",
		Sit:     "M14",
		Price:   20,
	}
)

type MovieInfo struct {
	Title       string  `json:"title"`
	Score       float32 `json:"score"`
	Country     string  `json:"country"`
	RunningTime uint64  `json:"runningTime"`
	Genre       string  `json:"genre"`
	Year        int     `json:"year"`
}

type TicketInfo struct {
	Date    time.Time `json:"date"`
	Theater string    `json:"theater"`
	Sit     string    `json:"sit"`
	Price   int       `json:"price"`
}

type PaymentInfo struct {
	PaymentDate        time.Time `json:"paymentDate"`
	PaymentTransaction string    `json:"paymentTransaction"`
	PointTransaction   string    `json:"pointTransaction"`
}

type PriceInfo struct {
	UsedFungible     int `json:"usedFungible"`
	UsedServiceToken int `json:"usedServiceToken"`
	SubTotal         int `json:"subTotal"`
	Discount         int `json:"discount"`
	GrandTotal       int `json:"grandTotal"`
}

type PurchaseInfo struct {
	MovieInfo  MovieInfo  `json:"movieInfo"`
	TicketInfo TicketInfo `json:"ticketInfo"`
	PriceInfo  PriceInfo  `json:"priceInfo"`
}

type UserInfo struct {
	userID        string `json:"userId"`
	WalletAddress string `json:"walletAddress"`
}

type ServiceTokenBalance struct {
	ContractID string `json:"contractId"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	ImgURI     string `json:"imgUri"`
	Amount     string `json:"amount"`
	Decimals   int    `json:"decimals"`
}

type FungibleBalance struct {
	Name      string `json:"name"`
	TokenType string `json:"tokenType"`
	Meta      string `json:"meta"`
	Amount    string `json:"amount"`
}

type NonFungibleInfo struct {
	Name       string `json:"name"`
	TokenIndex string `json:"tokenIndex"`
	Meta       string `json:"meta"`
}

type BaseCoinBalance struct {
	Symbol   string `json:"symbol"`
	Amount   string `json:"amount"`
	Decimals int    `json:"decimals"`
}

type Transaction struct {
	Height    uint64 `json:"height"`
	TxHash    string `json:"txhash"`
	Index     int    `json:"index"`
	Code      int    `json:"code"`
	RawLog    string `json:"raw_log"`
	Logs      []Log  `json:"logs"`
	GasWanted uint64 `json:"gasWanted"`
	GasUsed   uint64 `json:"gasUsed"`
	Tx        Tx     `json:"tx"`
	Timestamp string `json:"timestamp"`
}

type Log struct {
	MsgIndex int     `json:"msg_index"`
	Success  bool    `json:"success"`
	Log      string  `json:"log"`
	Events   []Event `json:"events"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Tx struct {
	Type  string  `json:"type"`
	Value TxValue `json:"value"`
}

type TxValue struct {
	Message    []Message   `json:"msg"`
	Fee        Fee         `json:"fee"`
	Signatures []Signature `json:"signatures"`
	Memo       string      `json:"memo"`
}

type Message struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TransferBaseCoinMsg struct {
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      Amount `json:"amount"`
}

type TransferServiceTokenMsg struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     uint64 `json:"amount"`
	ContractID string `json:"contractId"`
}

type FungibleMsg struct {
	Amount     []FungibleAmount `json:"amount"`
	ContractID string           `json:"contractId"`
}

type FungibleAmount struct {
	Amount  uint64 `json:"amount"`
	TokenID string `json:"tokenId"`
}

type MintNonFungibleMsg struct {
	From       string `json:"from"`
	To         string `json:"to"`
	ContractID string `json:"contractId"`
	Meta       string `json:"meta"`
	Name       string `json:"name"`
	TokenType  string `json:"tokenType"`
}

type Signature struct {
	PubKey    PubKey `json:"pubKey"`
	Signature string `json:"signature"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Fee struct {
	Amount []Amount `json:"amount"`
	Gas    int      `json:"gas"`
}

type Amount struct {
	Amount int    `json:"amount"`
	Denom  string `json:"denom"`
}

type NonFungibleTxHistory struct {
	PaymentTransaction *Transaction `json:"paymentTransaction"`
	MintTransaction    *Transaction `json:"mintTransaction"`
	PointTransaction   *Transaction `json:"pointTransaction"`
}

type NonFungibleMetadata struct {
	MovieInfo   MovieInfo   `json:"movieInfo"`
	TicketInfo  TicketInfo  `json:"ticketInfo"`
	PaymentInfo PaymentInfo `json:"paymentInfo"`
}

type TransactionAccepted struct {
	TxHash string `json:"txHash"`
}

type TransferRequestResult struct {
	RequestSessionToken string `json:"requestSessionToken"`
	RedirectURI         string `json:"redirectUri"`
}
