package eth

import (
	"encoding/json"
	"eth_test/util"
	"github.com/go-errors/errors"
	"log"
	"strconv"
)

type ETHClient interface {
	GetBlockDetailByNum(num int) (Block, error)
	GetBlockNum() (int, error)
	GetTransactionByHash(hash string) (Transaction, error)
}

func NewETHClient(endpoint string) ETHClient {
	return &ETHClientImpl{
		endpoint:       endpoint,
		userTransMap:   map[string]string{},
		transDetailMap: map[string]Transaction{},
	}
}

type ETHClientImpl struct {
	endpoint       string
	userTransMap   map[string]string
	transDetailMap map[string]Transaction
}

func (e *ETHClientImpl) GetBlockNum() (int, error) {
	req := util.NewJsonRequest(METHOD_GET_BLOCK_NUM)
	resp, err := util.RPCCall(e.endpoint, req)
	if err != nil {
		return 0, errors.Errorf("Rpc call failed, error %v", err)
	}
	hexString := string(resp.Result)[3 : len(resp.Result)-1]
	number, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(number), nil
}

func (e *ETHClientImpl) GetBlockDetailByNum(num int) (Block, error) {
	req := util.NewJsonRequest(METHOD_GET_BLOCK_BY_NUM, util.IntToHexString(num), true)
	resp, err := util.RPCCall(e.endpoint, req)
	if err != nil {
		return Block{}, errors.Errorf("Rpc call failed, error %v", err)
	}
	var bytes []byte
	err = resp.Result.UnmarshalJSON(bytes)
	if err != nil {
		return Block{}, errors.Errorf("Parse result failed, error %v", err)
	}
	block := Block{}
	err = json.Unmarshal(bytes, &block)
	if err != nil {
		return Block{}, err
	}
	return block, nil
}

func (e *ETHClientImpl) GetTransactionByHash(hash string) (Transaction, error) {

	req := util.NewJsonRequest(METHOD_GET_TRANSACTION_BY_HASH, hash, true)
	resp, err := util.RPCCall(e.endpoint, req)
	if err != nil {
		return Transaction{}, errors.Errorf("Rpc call failed, error %v", err)
	}
	transaction := Transaction{}
	err = json.Unmarshal(resp.Result, &transaction)
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}
