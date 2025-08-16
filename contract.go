package etherkit

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

//############ Contract ############

// GetABI 从abi字符串中获得ABI对象
func GetABI(abiStr string) (abi.ABI, error) {
	abiContract, err := abi.JSON(strings.NewReader(abiStr))
	return abiContract, err
}

// GetContractMethodId 获得合约的methodId
// 参数method，如：transfer(address,uint256)
func GetContractMethodId(method string) string {
	methodId := hexutil.Encode(crypto.Keccak256([]byte(method))[:4])
	return methodId
}

// GetEventTopic 获得事件的topic。event字符串如：transfer(address,uint256)
func GetEventTopic(event string) string {
	return crypto.Keccak256Hash([]byte(event)).String()
}

// BuildContractInputData 构建合约的input data
func BuildContractInputData(contract abi.ABI, name string, args ...interface{}) ([]byte, error) {
	return contract.Pack(name, args...)
}
