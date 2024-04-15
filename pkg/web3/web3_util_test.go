package web3

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"goboot/internal/config"
	"goboot/internal/consts"
	"goboot/pkg/mathx"
	"math/big"
	"testing"
)

const TestRPCurl = "https://1rpc.io/base-goerli"
const TestChainId = 84531
const ChainId = 80001
const GasLimit = 300000
const PrivateKey = "2be938b70fb87de337f94c65effa87c357fa0c4c4184347db3d20f498e800876"
const Address = "0x98B0eDd62367Bae6e5285eE7F103fAD4e97B716D"
const ContractAddr = "0x71c82f1894Ce89aBe824695A0b64b32ed4bf5dC8"
const AbiPath = "../web3/abi/CCToken.abi"
const NFTPath = "../../config/abi/EYWorker.abi"
const USDTPath = "../../config/abi/USDTToken.abi"

func TestMint(t *testing.T) {
	var abi = config.LoadAbiFile(AbiPath)
	fmt.Println(abi)
	params := []any{common.HexToAddress(Address), FloatToWei(10)}
	reqParam := ContractSendParam{
		RpcUrl:       TestRPCurl,
		ChainId:      ChainId,
		GasLimit:     GasLimit,
		From:         Address,
		PrivateKey:   PrivateKey,
		ContractAddr: ContractAddr,
		Abi:          *abi,
		Method:       "mint",
		Params:       params,
	}
	txHash, err := ContractSendWithAbi(reqParam, context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txHash)
}

func TestToWei(t *testing.T) {
	sum := float64(1)
	for sum <= 10 {
		sum += 0.1
		amount := FloatToWei(sum)
		fmt.Println(amount)
	}
}

func TestGenerateWallet(t *testing.T) {
	w, _ := GenerateWallet()
	fmt.Printf("addr: %s, priv: %s, pub: %s", w.Address, w.PrivateKey, w.PublicKey)
}

func TestFromWeiStr(t *testing.T) {
	str := big.NewInt(1000000000000000)
	fmt.Println(FromWeiStr(str))
}

func TestGetGas(t *testing.T) {
	//gas := GetGas("https://eth.llamarpc.com", big.NewInt(21000), ToGWei(50), context.Background())
	amountFloat := 0.015291822347228
	bal := 0.014291822347228
	gasPrice, gas := GetGas("https://mainnet.infura.io/v3/cad08cc0470f425193c5a3ec2a7893b7", big.NewInt(21000), 5, context.Background())
	//gas := GetGas("http://54.179.196.171", big.NewInt(21000), ToGWei(5), context.Background())
	gasFloat, _ := FromWei(gas).Float64()
	maxTransfer := GetMaxLeft(bal, gasFloat)
	chainAmount := amountFloat
	if amountFloat > maxTransfer {
		chainAmount = mathx.Trunc(mathx.Sub(chainAmount, 0.001), 6)
	}
	//gasFin := ToWei(big.NewFloat(gasFloat))
	fmt.Printf("gas: %v \n", gas)
	fmt.Printf("gasFloat:%f \n", gasFloat)
	fmt.Printf("maxTrans: %f \n", maxTransfer)
	fmt.Printf("chainAmount: %f \n", chainAmount)
	fmt.Printf("gasPrice:%v \n", FromGWei(gasPrice))
	TestSuggestGas(t)
}

func TestSuggestGas(t *testing.T) {
	rpcUrl := "https://mainnet.infura.io/v3/cad08cc0470f425193c5a3ec2a7893b7"
	gasPrice, gas := GetGas(rpcUrl, big.NewInt(74037), consts.NormalGasIncr, context.Background())
	gasFloat, _ := FromWei(gas).Float64()
	fmt.Printf("gas:%v", mathx.Mul(gasFloat, 2000))
	fmt.Printf("gasPrice:%v", FromGWei(gasPrice))
}

func TestGetMaxLeft(t *testing.T) {
	//left := GetMaxLeft(0.01, 0.001193235283086)
	//s:=FromWei(27747798116)
	//fmt.Println(left)
	m := 0.009224 + 0.000582703760436
	fmt.Println(m)
}

func TestFloatTo6Wei(t *testing.T) {
	w := FloatToWei(1)
	wu := FloatTo6Wei(10)
	fmt.Println(w)
	fmt.Println(wu)
}
func TestFrom6Wei(t *testing.T) {
	s := FromWei(big.NewInt(1000000000000000000))
	s2 := From6Wei(big.NewInt(100000000000000000))
	fmt.Println(s)
	fmt.Println(s2)
}
func TestTo6Wei(t *testing.T) {
	s := To6Wei(big.NewFloat(12300))
	fmt.Println(s)
}

func TestTransferNft(t *testing.T) {
	var abi = config.LoadAbiFile(NFTPath)
	fmt.Println(abi)
	sysAddr := "0x98B0eDd62367Bae6e5285eE7F103fAD4e97B716D"
	//sysPri := "2be938b70fb87de337f94c65effa87c357fa0c4c4184347db3d20f498e800876"
	toAddr := "0x27e44C6bBa4Cb79f2f14C3e9f1D13d8968d28B98"
	toPriv := "d2ddea126047258a06059a997bc5a6d3c385f4f0f370a257ec37d58bb85e5255"
	nftContract := "0x69b757A836C9F8254755086eBfC6BC620D64004b"

	params := []any{common.HexToAddress(toAddr), common.HexToAddress(sysAddr), big.NewInt(int64(74804330))}
	reqParam := ContractSendParam{
		RpcUrl:       TestRPCurl,
		ChainId:      TestChainId,
		GasLimit:     300000,
		From:         toAddr,
		PrivateKey:   toPriv,
		ContractAddr: nftContract,
		Abi:          *abi,
		Method:       "transferFrom",
		Params:       params,
	}
	_, err := TryCallContract(reqParam, context.Background())
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return
	}
	txHash, err := ContractSendWithAbi(reqParam, context.Background())
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	fmt.Println(txHash)
}

func TestAddSpeed(t *testing.T) {
	from := "0xA8d6A4168A0Cc31E12bF9371A2fED8E851D4951D"
	to := "0x2c33d63F5def798874306C0A115C221B95Fa64bB"
	priv := "f38114bf58b3c6ce9661695761a70bb7ceafbfb525f8304d11d806f2c7f2e903"
	rpcUrl := "https://mainnet.infura.io/v3/cad08cc0470f425193c5a3ec2a7893b7"
	gasPrice := ToGWei(30)
	param := TransferParams{
		RpcUrl:     rpcUrl,
		ChainId:    1,
		GasLimit:   21000,
		From:       from,
		PrivateKey: priv,
		To:         common.HexToAddress(to),
		Nonce:      654,
		GasPrice:   gasPrice,
		Value:      big.NewFloat(0.01397),
	}

	txHash, err := Transfer(param, context.Background())
	fmt.Println(txHash)
	fmt.Println(err)
}

func TestTransferERC20(t *testing.T) {
	from := "0x27e44C6bBa4Cb79f2f14C3e9f1D13d8968d28B98"
	to := "0x0cadBCFB0a309F7c52C9D6A01A5957cb982dBe2e"
	priv := "d2ddea126047258a06059a997bc5a6d3c385f4f0f370a257ec37d58bb85e5255"
	rpcUrl := "https://mainnet.infura.io/v3/cad08cc0470f425193c5a3ec2a7893b7"
	contractAddr := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	var abi = config.LoadAbiFile(USDTPath)
	param := []any{common.HexToAddress(to), To6Wei(big.NewFloat(3))}
	sendParam := ContractSendParam{
		RpcUrl:   rpcUrl,
		ChainId:  1,
		GasLimit: 67814,
		//GasIncrGWei:  5,
		GasPrice:     ToGWei(35),
		Nonce:        47,
		From:         from,
		PrivateKey:   priv,
		ContractAddr: contractAddr,
		Abi:          *abi,
		Method:       "transfer",
		Params:       param,
	}
	txHash, err := ContractSendWithAbi(sendParam, context.Background())
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("hash:", txHash)
	}

}

func TestTransaction(t *testing.T) {
	hash := "0x22b06a1af6e41e3ecef92004bc953b1def60ed04ea4405d710140c2447005e8f"
	rpc := "https://mainnet.infura.io/v3/cad08cc0470f425193c5a3ec2a7893b7"
	tx, isPadding, err := QueryTransactionByHash(rpc, hash, context.Background())
	if !isPadding && err == nil {
		println("ok")
	}
	fmt.Println(tx.To())
	//fmt.Println(padding)
	//fmt.Println(err)
}
