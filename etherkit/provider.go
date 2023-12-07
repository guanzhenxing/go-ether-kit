package etherkit

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

// EthProvider 需要通过链上查询的，但是不需要账户的
type EthProvider interface {
	GetEthClient() *ethclient.Client
	GetRpcClient() *rpc.Client
	Close()
	GetNetworkID() (*big.Int, error)
	GetChainID() (*big.Int, error)
	GetBlockByHash(hash common.Hash) (*types.Block, error)
	GetBlockByNumber(number *big.Int) (*types.Block, error)
	GetBlockNumber() (uint64, error)
	GetSuggestGasPrice() (*big.Int, error)
	GetTransactionByHash(hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error)
	GetContractBytecode(address common.Address) (string, error)
	IsContractAddress(address common.Address) (bool, error)
	EstimateGas(from, to common.Address, nonce uint64, gasPrice, value *big.Int, data []byte) (uint64, error)
	GetFromAddress(tx *types.Transaction) (common.Address, error)
}

func NewDefaultEthProvider(rawUrl string) (*DefaultEthProvider, error) {

	rpcClient, err := rpc.Dial(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to rpc.Dial(): %w", err)
	}

	return &DefaultEthProvider{
		rc: rpcClient,
		ec: ethclient.NewClient(rpcClient),
	}, nil
}

func NewDefaultEthProviderWithChainId(rawUrl string, chainId int64) (*DefaultEthProvider, error) {

	p, err := NewDefaultEthProvider(rawUrl)
	if err != nil {
		return nil, err
	}
	p.chainId = big.NewInt(chainId)

	return p, nil
}

type DefaultEthProvider struct {
	rc      *rpc.Client
	ec      *ethclient.Client
	chainId *big.Int
}

// GetEthClient 获得ethClient客户端
func (p *DefaultEthProvider) GetEthClient() *ethclient.Client {
	return p.ec
}

// GetRpcClient 获得rpcClient客户端
func (p *DefaultEthProvider) GetRpcClient() *rpc.Client {
	return p.rc
}

// Close 关闭ethClient客户端和rpcClient客户端
func (p *DefaultEthProvider) Close() {
	p.ec.Close()
	p.rc.Close()
}

// GetNetworkID 获得NetworkId
func (p *DefaultEthProvider) GetNetworkID() (*big.Int, error) {
	return p.ec.NetworkID(context.Background())
}

// GetChainID 获得ChainId
func (p *DefaultEthProvider) GetChainID() (*big.Int, error) {

	if p.chainId == nil {
		chainId, err := p.ec.ChainID(context.Background())
		if err != nil {
			return nil, err
		}
		p.chainId = chainId
	}

	return p.chainId, nil
}

// GetBlockByHash 根据区块Hash获得区块信息
func (p *DefaultEthProvider) GetBlockByHash(blkHash common.Hash) (*types.Block, error) {
	return p.ec.BlockByHash(context.Background(), blkHash)
}

// GetBlockByNumber 根据区块号获得区块信息
func (p *DefaultEthProvider) GetBlockByNumber(number *big.Int) (*types.Block, error) {
	return p.ec.BlockByNumber(context.Background(), number)
}

// GetBlockNumber 获得最新区块
func (p *DefaultEthProvider) GetBlockNumber() (uint64, error) {
	return p.ec.BlockNumber(context.Background())
}

// GetSuggestGasPrice 获得建议的Gas
func (p *DefaultEthProvider) GetSuggestGasPrice() (*big.Int, error) {
	return p.ec.SuggestGasPrice(context.Background())
}

// GetTransactionByHash 根据txHash获得交易信息
func (p *DefaultEthProvider) GetTransactionByHash(txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return p.ec.TransactionByHash(context.Background(), txHash)
}

// GetTransactionReceipt 根据txHash获得交易Receipt
func (p *DefaultEthProvider) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return p.ec.TransactionReceipt(context.Background(), txHash)
}

// GetContractBytecode 根据合约地址获得bytecode
func (p *DefaultEthProvider) GetContractBytecode(address common.Address) (string, error) {
	bytecode, err := p.ec.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytecode), nil
}

// IsContractAddress 是否是合约地址。
func (p *DefaultEthProvider) IsContractAddress(address common.Address) (bool, error) {
	//获取一个代币智能合约的字节码并检查其长度以验证它是一个智能合约
	if bytecode, err := p.GetContractBytecode(address); err == nil {
		return len(bytecode) > 0, nil
	} else {
		return false, err
	}
}

// EstimateGas 预估手续费
func (p *DefaultEthProvider) EstimateGas(from, to common.Address, nonce uint64, gasPrice, value *big.Int, data []byte) (uint64, error) {
	return p.ec.EstimateGas(context.Background(), ethereum.CallMsg{
		From:       from,
		To:         &to,
		GasPrice:   gasPrice,
		Value:      value,
		Data:       data,
		Gas:        0,
		GasFeeCap:  nil,
		GasTipCap:  nil,
		AccessList: nil,
	})
}

// GetFromAddress 获得交易的fromAddress
func (p *DefaultEthProvider) GetFromAddress(tx *types.Transaction) (common.Address, error) {
	return types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
}
