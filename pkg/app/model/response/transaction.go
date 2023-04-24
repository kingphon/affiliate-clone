package responsemodel

import (
	"git.selly.red/Selly-Server/affiliate/external/constant"
	optionsmodel "git.selly.red/Selly-Server/affiliate/external/model/options"
	"git.selly.red/Selly-Server/affiliate/external/utils/ptime"
	"git.selly.red/Selly-Server/affiliate/internal/config"
	appconstant "git.selly.red/Selly-Server/affiliate/pkg/app/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ResponseTransactionBrief ...
type ResponseTransactionBrief struct {
	ID              primitive.ObjectID        `json:"_id"`
	Code            string                    `json:"code"`
	TransactionTime *ptime.TimeResponse       `json:"transactionTime"`
	Source          string                    `json:"source"`
	Commission      float64                   `json:"commission"`
	Status          ResponseTransactionStatus `json:"status"`
	Campaign        ResponseCampaignShortInfo `json:"campaign"`
	RejectedReason  string                    `json:"rejectedReason,omitempty"`
}

// ResponseTransactionStatus ...
type ResponseTransactionStatus struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

// ResponseTransactionDetail ...
type ResponseTransactionDetail struct {
	ID                 primitive.ObjectID                   `json:"_id"`
	Code               string                               `json:"code"`
	TransactionTime    *ptime.TimeResponse                  `json:"transactionTime"`
	Source             string                               `json:"source"`
	Commission         float64                              `json:"commission"`
	Status             ResponseTransactionStatus            `json:"status"`
	Campaign           ResponseCampaignShortInfo            `json:"campaign"`
	RejectedReason     string                               `json:"rejectedReason,omitempty"`
	EstimateCashbackAt *ptime.TimeResponse                  `json:"estimateCashbackAt"`
	Device             ResponseTransactionDevice            `json:"device"`
	DeviceText         string                               `json:"deviceText"`
	ProcessIcons       []appconstant.TransactionProcessIcon `json:"processIcons"`
	LastHistory        ResponseTransactionHistory           `json:"lastHistory"`
	SupportChannel     optionsmodel.ActionType              `json:"supportChannel"`
}

const (
	URLSellyChatDev  = "https://chat.unibag.xyz/?source=selly_app"
	URLSellyChatProd = "https://chat.selly.vn/?source=selly_app"
)

// getURLChat ...
func getURLChat() string {
	if config.IsEnvProduction() {
		return URLSellyChatProd
	}

	return URLSellyChatDev
}

// GetSupportChannel ...
func (s *ResponseTransactionDetail) GetSupportChannel() optionsmodel.ActionType {
	return optionsmodel.ActionType{
		Type:  constant.ActionTypeOpenAppBrowser,
		Value: getURLChat(),
	}
}

// ResponseTransactionHistories ...
type ResponseTransactionHistories struct {
	Data []ResponseTransactionHistory `json:"data"`
}

// ResponseTransactionHistory ...
type ResponseTransactionHistory struct {
	ID        primitive.ObjectID  `json:"_id"`
	Status    string              `json:"status"`
	Desc      string              `json:"desc"`
	CreatedAt *ptime.TimeResponse `json:"createdAt"`
}

// ResponseTransactionDevice ...
type ResponseTransactionDevice struct {
	Model          string `json:"model"`
	UserAgent      string `json:"userAgent"`
	Manufacturer   string `json:"manufacturer"`
	OSName         string `json:"osName"`
	OSVersion      string `json:"osVersion"`
	BrowserVersion string `json:"browserVersion"`
	BrowserName    string `json:"browserName"`
	DeviceType     string `json:"deviceType"`
	DeviceID       string `json:"deviceId"`
}

// ResponseTransactionFilter ...
type ResponseTransactionFilter struct {
	Data []KeyValue `json:"data"`
}
