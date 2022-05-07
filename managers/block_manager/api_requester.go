package block_manager

import (
	"fmt"
	"github.com/AlexandrGayun/go_test_task/models/block"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
)

const BASE_URL = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%#x&boolean=true&apikey="
const WEI_TO_ETH_COEF = 1.e-18

var apiUrl = BASE_URL + os.Getenv("ETHERSCAN_API_KEY")

var httpClient = &http.Client{Timeout: 5 * time.Second}

func getBlockData(blockNumber uint64) (*block.Block, error) {
	body, err := performRequest(blockNumber)
	if err != nil {
		return nil, err
	}
	nongroupedData, err := parseJSONResponse(body)
	if err != nil {
		return nil, err
	}
	transactionsCount, totalAmount, err := groupTransactions(nongroupedData)
	if err != nil {
		return nil, err
	}
	block := &block.Block{Number: blockNumber, TransactionsCount: transactionsCount, TotalAmount: totalAmount}
	return block, nil
}

func performRequest(blockNumber uint64) ([]byte, error) {
	apiPath := fmt.Sprintf(apiUrl, blockNumber)
	resp, err := httpClient.Get(apiPath)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func parseJSONResponse(body []byte) ([]string, error) {
	p := fastjson.Parser{}
	val, err := p.ParseBytes(body)
	if err != nil {
		return nil, err
	}
	jsonArr := val.GetArray("result", "transactions")
	res := make([]string, len(jsonArr))
	for k, v := range jsonArr {
		str := string(v.GetStringBytes("value"))
		res[k] = str
	}
	return res, err
}

func groupTransactions(arr []string) (int, float64, error) {
	var (
		transactionsCount      int
		acc, transactionAmount = new(big.Float), new(big.Float)
		err                    error
	)
	for _, v := range arr {
		// could use high perf(?) community approved solution
		// https://pkg.go.dev/github.com/ethereum/go-ethereum/common/hexutil#DecodeBig
		if _, ok := transactionAmount.SetString(v); !ok {
			return 0, 0, fmt.Errorf("failed to convert transaction amount string val %s", v)
		}
		acc.Add(acc, transactionAmount)
		transactionsCount++
	}
	totalAmountWei, _ := acc.Float64()
	totalAmountEth := totalAmountWei * WEI_TO_ETH_COEF
	return transactionsCount, totalAmountEth, err
}
