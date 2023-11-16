package lagrange

const (
	methodAPIToken = "/api_token"
	methodWallet   = "/jwt_info"
)

func (client *Client) APIToken() (key string, err error) {
	var result TokenDetail
	result.Token.Token = &key
	err = client.get(methodAPIToken, nil, &result)
	return
}

type TokenDetail struct {
	Token struct {
		CreatedAt string  `json:"created_at"`
		Token     *string `json:"token"`
		UpdatedAt string  `json:"updated_at"`
	} `json:"token"`
}

func (client *Client) WalletAddr() (wallet string, err error) {
	var result WalletAddrResult
	result.WalletAddress = &wallet
	err = client.get(methodWallet, nil, &result)
	return
}

type WalletAddrResult struct {
	WalletAddress *string `json:"wallet_address"`
}