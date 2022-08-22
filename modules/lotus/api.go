package lotus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"gitee.com/szxjyt/filbox-backend/models"
	"gitee.com/szxjyt/filbox-backend/types"

	"github.com/sirupsen/logrus"
)

func WalletNew(link, token string) (string, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.WalletNew",
		"id":      1,
		"params":  []string{"secp256k1"},
	}
	bytesData, _ := json.Marshal(data)
	client := &http.Client{}
	req, _ := http.NewRequest("post", link+"/rpc/v0", bytes.NewReader(bytesData))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonResp map[string]interface{}
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return "", err
	}
	result := jsonResp["result"].(string)
	return result, err

}

func WalletBalance(link, token, wallet string) (string, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.WalletBalance",
		"id":      1,
		"params":  []string{wallet},
	}
	bytesData, _ := json.Marshal(data)
	client := &http.Client{}
	req, _ := http.NewRequest("post", link+"/rpc/v0", bytes.NewReader(bytesData))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonResp map[string]interface{}
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return "", err
	}
	result := jsonResp["result"].(string)
	return result, err
}

func ClientQueryAsk(link, token, miner string) (interface{}, error) {
	stateMinerInfo := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.StateMinerInfo",
		"id":      1,
		"params":  []interface{}{miner, nil},
	}
	peerid := ""
	result, _, err := doRequest(link, token, stateMinerInfo, 10)
	if err != nil {
		return nil, err
	}

	if minerinfo, ok := result.(map[string]interface{}); !ok {
		return nil, errors.New("stateMinerInfo err")
	} else {
		if vl, ok := minerinfo["PeerId"]; ok {
			peerid = vl.(string)
		}
	}

	if peerid == "" {
		return nil, errors.New("stateMinerInfo peerid err")
	}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.ClientQueryAsk",
		"id":      1,
		"params":  []string{peerid, miner},
	}
	res, _, err := doRequest(link, token, data, 10)
	return res, err
}

func GetDeal(link, token, dealcid string, deal *types.DealOnline) error {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.ClientGetDealInfo",
		"id":      1,
		"params":  []interface{}{map[string]string{"/": dealcid}},
	}
	_, bt, err := doRequest(link, token, data, 20)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, deal)
}

// TODO
func Makedeal(link, token string, file *models.File, deal *models.DealInfo) (string, error) {
	price, err := strconv.ParseFloat(deal.Price, 64)
	if err != nil {
		return "", err
	}
	epochPrice := fmt.Sprintf("%.0f", price*dealSize(float64(file.Size)))
	logrus.Infof("epoch price %s", epochPrice)
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Filecoin.ClientStartDeal",
		"id":      1,
		"params": []interface{}{map[string]interface{}{
			"Data": map[string]interface{}{
				"TransferType": "graphsync",
				"Root": map[string]interface{}{
					"/": deal.Filecid,
				},
			},
			"Wallet":             deal.Wallet,
			"Miner":              deal.Miner,
			"EpochPrice":         epochPrice,
			"MinBlocksDuration":  day2Epoch(deal.Duration),
			"ProviderCollateral": "0",
			"DealStartEpoch":     -1,
			"FastRetrieval":      true,
			"VerifiedDeal":       false,
		}},
	}
	_, bt, err := doRequest(link, token, data, 60*5)
	if err != nil {
		return "", err
	}
	// TODO
	value := struct {
		NAMING_FAILED string `json:"/"`
	}{}
	return value.NAMING_FAILED, json.Unmarshal(bt, &value)
}

func doRequest(link, token string, data map[string]interface{}, timeout int64) (interface{}, []byte, error) {
	bytesData, _ := json.Marshal(data)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("post", link+"/rpc/v0", bytes.NewReader(bytesData))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	logrus.Infof("%s", body)

	var jsonResp map[string]interface{}
	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return nil, nil, err
	}
	if res, ok := jsonResp["result"]; ok {
		bt, err := json.Marshal(res)
		if err != nil {
			logrus.Error("_______________%s", err.Error())
			return res, body, nil

		}
		return res, bt, nil
	}
	return jsonResp, body, errors.New(string(body))
}

func day2Epoch(d int) int {
	return d * 2880
}

// 订单字节数转换为Filecoin需要的字节
func dealSize(x float64) float64 {
	return 127 * (math.Pow(2, math.Ceil(math.Log2(math.Ceil(x/127))))) / math.Pow(10, 9)
}
