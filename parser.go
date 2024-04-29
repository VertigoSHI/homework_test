package main

import (
	"eth_test/client/eth"
	"github.com/go-errors/errors"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"time"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]eth.Transaction, error)
}

func NewParser(endpoint string) Parser {
	client := eth.NewETHClient(endpoint)
	return &parser{
		client:                client,
		lastBlock:             0,
		subscribedUserAddress: map[string]struct{}{},
		userTransactions:      map[string][]string{},
		transactionDetailMap:  map[string]eth.Transaction{},
	}
}

type parser struct {
	mutex                 sync.Mutex
	lastBlock             int
	client                eth.ETHClient
	subscribedUserAddress map[string]struct{}
	userTransactions      map[string][]string
	transactionDetailMap  map[string]eth.Transaction
}

func (p *parser) GetCurrentBlock() int {
	return p.lastBlock
}

func (p *parser) Subscribe(address string) bool {
	p.subscribedUserAddress[address] = struct{}{}
	return true
}

func (p *parser) GetTransactions(address string) ([]eth.Transaction, error) {
	latestBlockNum, err := p.client.GetBlockNum()
	if err != nil {
		return nil, errors.Errorf("[parser][GetTransactions]Get block num failed. Error %v", err)
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var recordedTrans []eth.Transaction
	if _, ok := p.subscribedUserAddress[address]; !ok {
		return []eth.Transaction{}, errors.Errorf("[parser][GetTransactions]User %v not subscribed", address)
	}
	transHash := p.userTransactions[address]
	for _, hash := range transHash {
		recordedTrans = append(recordedTrans, p.transactionDetailMap[hash])
	}
	if latestBlockNum == p.lastBlock {
		return recordedTrans, nil
	}

	waitGroup := errgroup.Group{}
	var routineNum atomic.Int32
	routineNum.Store(1000)
	for checkingBlockNum := p.lastBlock + 1; checkingBlockNum <= latestBlockNum; checkingBlockNum++ {
		checkingBlockNum := checkingBlockNum
		for routineNum.Load() <= 0 {
			time.Sleep(500 * time.Millisecond)
		}
		routineNum.Add(-1)
		waitGroup.Go(func() error {
			defer routineNum.Add(1)
			block, err := p.client.GetBlockDetailByNum(checkingBlockNum)
			if err != nil {

				return err
			}
			for _, transaction := range block.Transactions {
				p.userTransactions[transaction.From] = append(p.userTransactions[transaction.From], transaction.Hash)
				p.userTransactions[transaction.To] = append(p.userTransactions[transaction.To], transaction.Hash)
				p.transactionDetailMap[transaction.Hash] = transaction
				if transaction.From == address || transaction.To == address {
					recordedTrans = append(recordedTrans, transaction)
				}
			}
			return nil
		})
	}
	err = waitGroup.Wait()
	if err != nil {
		return []eth.Transaction{}, err
	}
	return recordedTrans, nil
}
