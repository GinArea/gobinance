package binapi

type AccountInfoRequest struct{}

// https://developers.binance.com/docs/wallet/account
type AccountInfo struct {
	VipLevel                       int  `json:"vipLevel"`
	IsMarginEnabled                bool `json:"isMarginEnabled"`
	IsFutureEnabled                bool `json:"isFutureEnabled"`
	IsOptionsEnabled               bool `json:"isOptionsEnabled"`
	IsPortfolioMarginRetailEnabled bool `json:"isPortfolioMarginRetailEnabled"`
}

func (o *Client) GetAccountInfo() Response[AccountInfo] {
	return AccountInfoRequest{}.Do(o)
}

func (o AccountInfoRequest) Do(c *Client) Response[AccountInfo] {
	c.WithBaseUrl(SpotOrSignedBaseUrl)
	return Get(c, "sapi/v1/account/info", o, identity[AccountInfo])
}

type APIRestrictionsRequest struct{}

// https://developers.binance.com/docs/wallet/account/api-key-permission
type APIRestrictions struct {
	IPRestrict                   bool  `json:"ipRestrict"`
	CreateTime                   int64 `json:"createTime"`
	EnableReading                bool  `json:"enableReading"`
	EnableWithdrawals            bool  `json:"enableWithdrawals"`
	EnableInternalTransfer       bool  `json:"enableInternalTransfer"`
	EnableMargin                 bool  `json:"enableMargin"`
	EnableFutures                bool  `json:"enableFutures"`
	PermitsUniversalTransfer     bool  `json:"permitsUniversalTransfer"`
	EnableVanillaOptions         bool  `json:"enableVanillaOptions"`
	EnableFixApiTrade            bool  `json:"enableFixApiTrade"`
	EnableFixReadOnly            bool  `json:"enableFixReadOnly"`
	EnableSpotAndMarginTrading   bool  `json:"enableSpotAndMarginTrading"`
	EnablePortfolioMarginTrading bool  `json:"enablePortfolioMarginTrading"`
}

func (o *Client) GetApiRestrictions() Response[APIRestrictions] {
	return APIRestrictionsRequest{}.Do(o)
}

func (o APIRestrictionsRequest) Do(c *Client) Response[APIRestrictions] {
	c.WithBaseUrl(SpotOrSignedBaseUrl)
	return Get(c, "sapi/v1/account/apiRestrictions", o, identity[APIRestrictions])
}
