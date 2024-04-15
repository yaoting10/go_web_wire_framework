package web3

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"goboot/pkg/mathx"
	"log"
	"math/big"
	"regexp"
	"strings"
)

const defaultGasIncrGWei = float64(1)
const (
	TRecharge = 1
	TWithdraw = 2
	TSend     = 3
)

func GetEthParams(rpcUrl string, from string, gasIncrGWei float64, ctx context.Context) (*ethclient.Client, uint64, *big.Int, error) {
	client, err := ethclient.Dial(rpcUrl)
	nonce, err := client.NonceAt(ctx, common.HexToAddress(from), big.NewInt(-1))
	gasIncr := ToGWei(gasIncrGWei)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, 0, nil, err
	}
	gasPrice = gasPrice.Add(gasPrice, gasIncr)
	return client, nonce, gasPrice, nil
}

func GetGas(rpcUrl string, gasLimit *big.Int, gasIncr float64, ctx context.Context) (*big.Int, *big.Int) {
	client, err := ethclient.Dial(rpcUrl)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return big.NewInt(0), big.NewInt(0)
	}
	gasAmount := big.NewInt(0).Add(gasPrice, ToGWei(gasIncr))
	gas := big.NewInt(0).Mul(gasAmount, gasLimit)
	return gasAmount, gas
}

func GetGasPrice(rpcUrl string, gasIncr *big.Int, ctx context.Context) *big.Int {
	client, err := ethclient.Dial(rpcUrl)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return big.NewInt(0)
	}
	gasPrice = gasPrice.Add(gasPrice, gasIncr)
	return gasPrice
}

func GetNFTGasLimit(amount float64) int64 {
	gasLimit := int64(143882)

	if amount > 100000 {
		gasLimit = 181654
	}
	return gasLimit
}

func GetOtherGasLimit(typeVal int, amount float64) int64 {
	gasLimit := int64(21000)
	if typeVal == TRecharge {
		if amount <= 10 {
			gasLimit = 67814
		} else if amount < 100 {
			gasLimit = 74037
		} else if amount < 5000 {
			gasLimit = 85062
		} else {
			gasLimit = 109357
		}
	} else if typeVal == TWithdraw {
		if amount <= 10 {
			gasLimit = 67814
		} else if amount < 100 {
			gasLimit = 74037
		} else if amount < 5000 {
			gasLimit = 85062
		} else {
			gasLimit = 109357
		}
	}
	return gasLimit
}

func GetMaxLeft(amount float64, gasFloat float64) float64 {
	//gasFloat := float64(gas.Int64())
	//return mathx.Div(mathx.Sub(mathx.Mul(math.Pow10(18), amount), gasFloat), math.Pow10(18), 18)
	//0.008806764716914
	return mathx.Sub(amount, gasFloat)
}

func GetBalance(rpcUrl string, address string, blockNumber *big.Int, ctx context.Context) (*big.Float, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	account := common.HexToAddress(address)
	balance := big.NewInt(0)
	if blockNumber == nil || blockNumber.Int64() <= 0 {
		balance, err = client.BalanceAt(ctx, account, nil)
	} else {
		balance, err = client.BalanceAt(ctx, account, blockNumber)
	}
	if err != nil {
		return nil, err
	}
	amount := FromWei(balance)
	return amount, nil
}

func Transfer(p TransferParams, ctx context.Context) (string, error) {
	gasIncr := getGasIncr(p.GasIncrGWei)
	client, nonce, gasPrice, err := GetEthParams(p.RpcUrl, p.From, gasIncr, ctx)
	if err != nil {
		return "", err
	}
	if p.GasPrice != nil && p.GasPrice.Int64() > 0 {
		gasPrice = p.GasPrice
	}
	if p.Nonce != 0 {
		nonce = p.Nonce
	}
	//fmt.Printf("final gasPrice :%v", gasPrice)
	v, _ := p.Value.Float64()
	value := FloatToWei(v) // in wei (1 eth)
	toAddress := p.To
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, p.GasLimit, gasPrice, data)
	privateKey, err := crypto.HexToECDSA(p.PrivateKey)
	chainID, err := client.NetworkID(ctx)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}
	// fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

func ContractCallWithAbi(p ContractCallParam, ctx context.Context) (interface{}, error) {
	client, err := ethclient.Dial(p.RpcUrl)
	if err != nil {
		return "", err
	}
	m := p.Abi.Methods[p.Method]
	if m.Name == "" {
		return nil, errors.New("abi or method errors")
	}
	data := m.ID
	if len(p.Params) > 0 {
		inputData, err := m.Inputs.Pack(p.Params...)
		if err != nil {
			return nil, errors.New("params error")
		}
		data = append(data, inputData...)
	}

	contractAddress := common.HexToAddress(p.ContractAddr)
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	out, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		panic(err)
	}

	response, err := p.Abi.Unpack(p.Method, out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response : %v\n", response)
	if len(response) != 1 {
		return response, nil
	}
	return response[0], nil
}

func ContractSendWithAbi(p ContractSendParam, ctx context.Context) (string, error) {
	gasIncr := getGasIncr(p.GasIncrGWei)
	client, nonce, gasPrice, err := GetEthParams(p.RpcUrl, p.From, gasIncr, ctx)
	if err != nil {
		return "", err
	}
	if p.GasPrice != nil && p.GasPrice.Int64() > 0 {
		gasPrice = p.GasPrice
	}
	if p.Nonce != 0 {
		nonce = p.Nonce
	}
	contractABI := p.Abi
	if contractABI.Methods[p.Method].Name != p.Method {
		return "", errors.New("abi errors")
	}
	// todo original data
	// callData, err := contractABI.Pack(p.Method, common.HexToAddress(to), toWei(amount))
	callData, err := contractABI.Pack(p.Method, p.Params...)
	if callData == nil {
		return "", errors.New("abi params errors")
	}
	tx := types.NewTransaction(nonce, common.HexToAddress(p.ContractAddr), big.NewInt(0), p.GasLimit, gasPrice, callData)
	privateKey, err := crypto.HexToECDSA(p.PrivateKey)
	SignedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(p.ChainId)), privateKey)
	err = client.SendTransaction(ctx, SignedTx)
	hash := SignedTx.Hash().String()
	if err != nil {
		return "", err
	}
	log.Printf("hash:" + hash)
	return hash, nil
}

func getGasIncr(gasIncrGWei float64) float64 {
	gasIncr := float64(0)
	if gasIncrGWei == 0 {
		gasIncr = defaultGasIncrGWei
	} else {
		gasIncr = gasIncrGWei
	}
	return gasIncr
}

func TryCallContract(p ContractSendParam, ctx context.Context) (any, error) {
	client, err := ethclient.Dial(p.RpcUrl)
	if err != nil {
		return "", err
	}
	m := p.Abi.Methods[p.Method]
	if m.Name == "" {
		return nil, errors.New("abi or method errors")
	}
	data := m.ID
	if len(p.Params) > 0 {
		inputData, err := m.Inputs.Pack(p.Params...)
		if err != nil {
			return nil, errors.New("params error")
		}
		data = append(data, inputData...)
	}

	contractAddress := common.HexToAddress(p.ContractAddr)
	msg := ethereum.CallMsg{
		From: common.HexToAddress(p.From),
		To:   &contractAddress,
		Data: data,
	}
	client.CallContract(ctx, msg, nil)
	out, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		panic(err)
	}

	response, err := p.Abi.Unpack(p.Method, out)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response : %v\n", response)
	if len(response) != 1 {
		return response, nil
	}
	return response[0], nil
}

func QueryTransactionByHash(rpcUrl string, hash string, ctx context.Context) (tx *types.Transaction, isPending bool, err error) {
	client, _ := ethclient.Dial(rpcUrl)
	return client.TransactionByHash(ctx, common.HexToHash(hash))
}

func FromWei(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func FromGWei(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.GWei))
}

func From6Wei(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(1e6))
}

func FromWeiStr(amount *big.Int) string {
	compact_amount := big.NewInt(0)
	reminder := big.NewInt(0)
	divisor := big.NewInt(1e18)
	compact_amount.QuoRem(amount, divisor, reminder)
	return fmt.Sprintf("%v.%018s", compact_amount.String(), reminder.String())
}

func FloatToWei(amount float64) *big.Int {
	s := fmt.Sprintf("%f", amount)
	a := strings.Split(s, ".")
	for i := 0; i < 18; i++ {
		if len(a) > 1 && len(a[1]) > i {
			a[0] += string(a[1][i])
		} else {
			a[0] += "0"
		}
	}
	b, ok := new(big.Int).SetString(a[0], 10)
	if !ok {
		panic("Could not set big.Int string for value " + s)
	}
	return b
}

func FloatTo6Wei(amount float64) *big.Int {
	s := fmt.Sprintf("%f", amount)
	a := strings.Split(s, ".")
	for i := 0; i < 6; i++ {
		if len(a) > 1 && len(a[1]) > i {
			a[0] += string(a[1][i])
		} else {
			a[0] += "0"
		}
	}
	b, ok := new(big.Int).SetString(a[0], 10)
	if !ok {
		panic("Could not set big.Int string for value " + s)
	}
	return b
}

// FIXME loss precision
func ToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func RemoveTo6Wei(amount *big.Int) *big.Int {
	v := amount.Int64()
	f := float64(v) * 0.000000000001
	return big.NewInt(int64(f))
}

func To6Wei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(1e6))
	fracStr := strings.Split(fmt.Sprintf("%.6f", eth), ".")[1]
	fracStr += strings.Repeat("0", 6-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func ToGWei(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	expF := new(big.Float)
	expF.SetInt(exp)
	bigval.Mul(bigval, expF)
	result := new(big.Int)
	bigval.Int(result) // store converted number in result
	return result
}

func ParseBigFloat(value string) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	if err != nil {
		panic(err)
	}
	return f
}

func MnemonicFun() (string, error) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// GenerateWithMnemonic generate with mnemonic need password to seed
func GenerateWithMnemonic(mnemonic string, password string) (Wallet, error) {
	// mnemonic, _ := MnemonicFun()
	seed := bip39.NewSeed(mnemonic, password) // 这里可以选择传入指定密码或者空字符串，不同密码生成的助记词不同
	wallet, err := hdwallet.NewFromSeed(seed)
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") // 最后一位是同一个助记词的地址id，从0开始，相同助记词可以生产无限个地址
	account, err := wallet.Derive(path, false)

	address := account.Address.Hex()
	privateKey, _ := wallet.PrivateKeyHex(account)
	publicKey, _ := wallet.PublicKeyHex(account)
	return Wallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Mnemonic:   mnemonic,
	}, err
}

func GenerateWallet() (Wallet, error) {
	privateKey, _ := crypto.GenerateKey()
	return GenerateWithPrivateKey(privateKey)
}

func GenerateWithPrivate(private string) (Wallet, error) {
	priKey, _ := crypto.HexToECDSA(private)
	return GenerateWithPrivateKey(priKey)
}

func GenerateWithPrivateKey(priKey *ecdsa.PrivateKey) (Wallet, error) {
	var err error
	pubKey := priKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("private key error")
	}
	privateKeyBytes := crypto.FromECDSA(priKey)
	privateKey := hexutil.Encode(privateKeyBytes)[2:]
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKey := hexutil.Encode(publicKeyBytes)[4:]
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// TODO 无法逆向生成 mnemonic
	// entropy, _ := hex.DecodeString(privateKey)
	// mnemonic, _ := bip39.NewMnemonic(entropy)
	return Wallet{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, err
}

func IsValidAddress(address string) bool {
	// return ToValidateAddress(address) == address
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func ToValidateAddress(address string) string {
	addrLowerStr := strings.ToLower(address)
	if strings.HasPrefix(addrLowerStr, "0x") {
		addrLowerStr = addrLowerStr[2:]
		address = address[2:]
	}
	var binaryStr string
	addrBytes := []byte(addrLowerStr)
	hash256 := crypto.Keccak256Hash([]byte(addrLowerStr)) // 注意，这里是直接对字符串转换成byte切片然后哈希

	for i, e := range addrLowerStr {
		// 如果是数字则跳过
		if e >= '0' && e <= '9' {
			continue
		} else {
			binaryStr = fmt.Sprintf("%08b", hash256[i/2]) // 注意，这里一定要填充0
			if binaryStr[4*(i%2)] == '1' {
				addrBytes[i] -= 32
			}
		}
	}
	return "0x" + string(addrBytes)
}

func ContractContains(contracts []string, txTo string) bool {
	for _, v := range contracts {
		if strings.ToLower(v) == strings.ToLower(txTo) {
			return true
		}
	}
	return false
}

type Wallet struct {
	Address    string
	PrivateKey string
	PublicKey  string
	Mnemonic   string
}

type TransferParams struct {
	RpcUrl      string
	ChainId     int64
	GasLimit    uint64
	From        string
	PrivateKey  string
	GasIncrGWei float64
	Nonce       uint64
	GasPrice    *big.Int
	To          common.Address
	Value       *big.Float
}

type ContractSendParam struct {
	RpcUrl       string
	ChainId      int64
	GasLimit     uint64
	GasIncrGWei  float64
	GasPrice     *big.Int
	Nonce        uint64
	From         string
	PrivateKey   string
	ContractAddr string
	Abi          abi.ABI
	Method       string
	Params       []any
}

type ContractCallParam struct {
	RpcUrl       string
	ContractAddr string
	Abi          abi.ABI
	Method       string
	Params       []any
}

type ZeroValueCallMsg struct {
	From common.Address `json:"from,omitempty"`
	To   common.Address `json:"to"`
	Data CallMsgData    `json:"data"`
}

type CallMsgData []byte
