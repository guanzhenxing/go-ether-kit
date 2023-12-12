package etherkit

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

// Wallet 钱包信息
type Wallet interface {
	GetEthSigner() Signer
	GetEthProvider() Provider
	GetClient() *ethclient.Client
	GetAddress() common.Address
	CloseWallet()
	GetNonce() (uint64, error)
	GetBalance() (*big.Int, error)
	NewTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error)
	SendTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error)
	NewTxWithHexInput(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error)
	SendTxWithHexInput(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error)
	BuildTxOpts(value, nonce, gasPrice *big.Int) (*bind.TransactOpts, error)
	SignTx(tx *types.Transaction) (*types.Transaction, error)
	SendSignedTx(signedTx *types.Transaction) (common.Hash, error)
	Signature(data []byte) ([]byte, error)
	CallContract(contractAddress common.Address, contractAbi abi.ABI, functionName string, params ...interface{}) ([]interface{}, error)
}

type EtherWallet struct {
	es Signer
	ep Provider
}

// NewEtherWallet 新建一个Wallet
func NewEtherWallet(hexPk string, rawUrl string) (*EtherWallet, error) {
	es, err := NewEtherSignerFromHexPrivateKey(hexPk)
	if err != nil {
		return nil, err
	}

	ep, err := NewEtherProvider(rawUrl)
	if err != nil {
		return nil, err
	}

	return NewEtherWallet1(es, ep)
}

// NewEtherWallet1 新建Wallet
func NewEtherWallet1(es Signer, ep Provider) (*EtherWallet, error) {
	return &EtherWallet{
		es: es,
		ep: ep,
	}, nil
}

// GetEthClient 获得ethClient客户端
func (w *EtherWallet) getEthClient() *ethclient.Client {
	return w.ep.GetEthClient()
}

// GetRpcClient 获得rpcClient客户端
func (w *EtherWallet) getRpcClient() *rpc.Client {
	return w.ep.GetRpcClient()
}

// GetEthSigner 获得EthSinger
func (w *EtherWallet) GetEthSigner() Signer {
	return w.es
}

// GetEthProvider 获得EthProvider
func (w *EtherWallet) GetEthProvider() Provider {
	return w.ep
}

func (w *EtherWallet) GetClient() *ethclient.Client {
	return w.getEthClient()
}

// GetAddress 获得地址
func (w *EtherWallet) GetAddress() common.Address {
	return w.es.GetAddress()
}

// CloseWallet 关闭Wallet
func (w *EtherWallet) CloseWallet() {
	w.ep.Close()
}

// GetNonce 获得nonce
func (w *EtherWallet) GetNonce() (uint64, error) {
	return w.getEthClient().PendingNonceAt(context.Background(), w.GetAddress())
}

// GetBalance 获得本位币的约
func (w *EtherWallet) GetBalance() (*big.Int, error) {
	return w.getEthClient().BalanceAt(context.Background(), w.GetAddress(), nil)
}

// NewTx 构建一笔交易。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *EtherWallet) NewTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (*types.Transaction, error) {

	if nonce == 0 {
		var err error
		nonce, err = w.GetNonce()
		if err != nil {
			return nil, err
		}
	}

	if gasPrice == nil || gasPrice.Sign() == 0 {
		var err error
		gasPrice, err = w.GetEthProvider().GetSuggestGasPrice()
		if err != nil {
			return nil, err
		}
	}

	if gasLimit == 0 {
		var err error
		gasLimit, err = w.ep.EstimateGas(w.GetAddress(), to, nonce, gasPrice, value, data)
		if err != nil {
			return nil, err
		}
	}

	return NewTx(to, nonce, gasLimit, gasPrice, value, data)
}

// SendTx 发送交易。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *EtherWallet) SendTx(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, data []byte) (common.Hash, error) {

	tx, err := w.NewTx(to, nonce, gasLimit, gasPrice, value, data)
	if err != nil {
		return [32]byte{}, err
	}

	signedTx, err := w.SignTx(tx)
	if err != nil {
		return [32]byte{}, err
	}

	return w.SendSignedTx(signedTx)
}

// NewTxWithHexInput 构建一笔交易，使用0x开头的input。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *EtherWallet) NewTxWithHexInput(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (*types.Transaction, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return nil, err
	}
	return w.NewTx(to, nonce, gasLimit, gasPrice, value, data)
}

// SendTxWithHexInput 发送一笔交易，使用0x开头的input。nonce传0表示字段计算；gasLimit传0表示字段计算；gasPrice穿nil或者big.NewInt(0)表示gasPrice自动计算。
func (w *EtherWallet) SendTxWithHexInput(to common.Address, nonce, gasLimit uint64, gasPrice, value *big.Int, input string) (common.Hash, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return [32]byte{}, err
	}
	return w.SendTx(to, nonce, gasLimit, gasPrice, value, data)
}

// BuildTxOpts 构建交易的选项
func (w *EtherWallet) BuildTxOpts(value, nonce, gasPrice *big.Int) (*bind.TransactOpts, error) {

	chainId, err := w.ep.GetChainID()
	if err != nil {
		return nil, err
	}

	txOpts, _ := bind.NewKeyedTransactorWithChainID(w.GetEthSigner().GetPrivateKey(), chainId)

	txOpts.Value = value

	if gasPrice != nil && gasPrice.Sign() == 1 {
		txOpts.GasPrice = gasPrice
	} else {
		_gasPrice, err := w.GetEthProvider().GetSuggestGasPrice()
		if err != nil {
			return nil, err
		}
		txOpts.GasPrice = _gasPrice
	}

	// 如果nonce不为nil，就用传入的值（这里默认 nonce >0 ）
	if nonce != nil && nonce.Sign() > 0 {
		txOpts.Nonce = nonce
	} else {
		_nonce, err := w.GetNonce()
		if err != nil {
			return nil, err
		}
		txOpts.Nonce = big.NewInt(int64(_nonce))
	}

	return txOpts, nil
}

// SignTx 对交易进行签名
func (w *EtherWallet) SignTx(tx *types.Transaction) (*types.Transaction, error) {

	chainId, err := w.ep.GetChainID()
	if err != nil {
		return nil, err
	}

	// 使用伦敦签名
	signer := types.NewLondonSigner(chainId)
	signedTx, err := types.SignTx(tx, signer, w.es.GetPrivateKey())
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

// SendSignedTx 发送签名后的Tx
func (w *EtherWallet) SendSignedTx(signedTx *types.Transaction) (common.Hash, error) {
	err := w.getEthClient().SendTransaction(context.Background(), signedTx)
	if err != nil {
		return [32]byte{}, err
	}
	return signedTx.Hash(), nil
}

// Signature 生成一个签名
func (w *EtherWallet) Signature(data []byte) ([]byte, error) {
	key := w.es.GetPrivateKey()
	hash := crypto.Keccak256Hash(data)

	return crypto.Sign(hash.Bytes(), key)
}

// CallContract 调用合约的方法，无需创建交易
func (w *EtherWallet) CallContract(contractAddress common.Address, contractAbi abi.ABI, functionName string, params ...interface{}) ([]interface{}, error) {

	inputData, err := BuildContractInputData(contractAbi, functionName, params...)
	if err != nil {
		return nil, err
	}

	res, err := w.getEthClient().CallContract(context.TODO(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: inputData,
	}, nil)
	if err != nil {
		return nil, err
	}

	response, err := contractAbi.Unpack(functionName, res)
	if err != nil {
		return nil, err
	}
	return response, nil
}
