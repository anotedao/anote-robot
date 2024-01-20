package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"gopkg.in/telebot.v3"
)

func EncryptMessage(message string) string {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func DecryptMessage(message string) string {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if len(cipherText) < aes.BlockSize {
		log.Println(err)
		logTelegram(err.Error())
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func dataTransaction(key string, valueStr *string, valueInt *int64, valueBool *bool) error {
	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKey)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	a, _ := proto.NewAddressFromPublicKey(55, sender)
	log.Println(a)

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	tr := proto.NewUnsignedDataWithProofs(2, sender, Fee, uint64(ts))

	if valueStr == nil && valueInt == nil && valueBool == nil {
		tr.Entries = append(tr.Entries,
			&proto.DeleteDataEntry{
				Key: key,
			},
		)
	}

	if valueStr != nil {
		tr.Entries = append(tr.Entries,
			&proto.StringDataEntry{
				Key:   key,
				Value: *valueStr,
			},
		)
	}

	if valueInt != nil {
		tr.Entries = append(tr.Entries,
			&proto.IntegerDataEntry{
				Key:   key,
				Value: *valueInt,
			},
		)
	}

	if valueBool != nil {
		tr.Entries = append(tr.Entries,
			&proto.BooleanDataEntry{
				Key:   key,
				Value: *valueBool,
			},
		)
	}

	err = tr.Sign(55, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = cl.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func dataTransaction2(key string, valueStr *string, valueInt *int64, valueBool *bool) error {
	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyToday)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKeyToday)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	a, _ := proto.NewAddressFromPublicKey(55, sender)
	log.Println(a)

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	tr := proto.NewUnsignedDataWithProofs(2, sender, Fee, uint64(ts))

	if valueStr == nil && valueInt == nil && valueBool == nil {
		tr.Entries = append(tr.Entries,
			&proto.DeleteDataEntry{
				Key: key,
			},
		)
	}

	if valueStr != nil {
		tr.Entries = append(tr.Entries,
			&proto.StringDataEntry{
				Key:   key,
				Value: *valueStr,
			},
		)
	}

	if valueInt != nil {
		tr.Entries = append(tr.Entries,
			&proto.IntegerDataEntry{
				Key:   key,
				Value: *valueInt,
			},
		)
	}

	if valueBool != nil {
		tr.Entries = append(tr.Entries,
			&proto.BooleanDataEntry{
				Key:   key,
				Value: *valueBool,
			},
		)
	}

	err = tr.Sign(55, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = cl.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func dataTransactionAlpha(key string, valueStr *string, valueInt *int64, valueBool *bool) error {
	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyAlpha)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKeyAlpha)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	a, _ := proto.NewAddressFromPublicKey(55, sender)
	log.Println(a)

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	tr := proto.NewUnsignedDataWithProofs(2, sender, Fee, uint64(ts))

	if valueStr == nil && valueInt == nil && valueBool == nil {
		tr.Entries = append(tr.Entries,
			&proto.DeleteDataEntry{
				Key: key,
			},
		)
	}

	if valueStr != nil {
		tr.Entries = append(tr.Entries,
			&proto.StringDataEntry{
				Key:   key,
				Value: *valueStr,
			},
		)
	}

	if valueInt != nil {
		tr.Entries = append(tr.Entries,
			&proto.IntegerDataEntry{
				Key:   key,
				Value: *valueInt,
			},
		)
	}

	if valueBool != nil {
		tr.Entries = append(tr.Entries,
			&proto.BooleanDataEntry{
				Key:   key,
				Value: *valueBool,
			},
		)
	}

	err = tr.Sign(55, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = cl.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func getData(key string, address *string) (interface{}, error) {
	var a proto.WavesAddress

	wc, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if address == nil {
		pk, err := crypto.NewPublicKeyFromBase58(conf.PublicKey)
		if err != nil {
			return nil, err
		}

		a, err = proto.NewAddressFromPublicKey(55, pk)
		if err != nil {
			return nil, err
		}
	} else {
		a, err = proto.NewAddressFromString(*address)
		if err != nil {
			return nil, err
		}
	}

	ad, _, err := wc.Addresses.AddressesDataKey(context.Background(), a, key)
	if err != nil {
		return nil, err
	}

	if ad.GetValueType().String() == "string" {
		return ad.ToProtobuf().GetStringValue(), nil
	}

	if ad.GetValueType().String() == "boolean" {
		return ad.ToProtobuf().GetBoolValue(), nil
	}

	if ad.GetValueType().String() == "integer" {
		return ad.ToProtobuf().GetIntValue(), nil
	}

	return "", nil
}

func getData2(key string, address *string) (interface{}, error) {
	var a proto.WavesAddress

	wc, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if address == nil {
		pk, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyToday)
		if err != nil {
			return nil, err
		}

		a, err = proto.NewAddressFromPublicKey(55, pk)
		if err != nil {
			return nil, err
		}
	} else {
		a, err = proto.NewAddressFromString(*address)
		if err != nil {
			return nil, err
		}
	}

	ad, _, err := wc.Addresses.AddressesDataKey(context.Background(), a, key)
	if err != nil {
		return nil, err
	}

	if ad.GetValueType().String() == "string" {
		return ad.ToProtobuf().GetStringValue(), nil
	}

	if ad.GetValueType().String() == "boolean" {
		return ad.ToProtobuf().GetBoolValue(), nil
	}

	if ad.GetValueType().String() == "integer" {
		return ad.ToProtobuf().GetIntValue(), nil
	}

	return "", nil
}

func getDataAlpha(key string, address *string) (interface{}, error) {
	var a proto.WavesAddress

	wc, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if address == nil {
		pk, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyAlpha)
		if err != nil {
			return nil, err
		}

		a, err = proto.NewAddressFromPublicKey(55, pk)
		if err != nil {
			return nil, err
		}
	} else {
		a, err = proto.NewAddressFromString(*address)
		if err != nil {
			return nil, err
		}
	}

	ad, _, err := wc.Addresses.AddressesDataKey(context.Background(), a, key)
	if err != nil {
		return nil, err
	}

	if ad.GetValueType().String() == "string" {
		return ad.ToProtobuf().GetStringValue(), nil
	}

	if ad.GetValueType().String() == "boolean" {
		return ad.ToProtobuf().GetBoolValue(), nil
	}

	if ad.GetValueType().String() == "integer" {
		return ad.ToProtobuf().GetIntValue(), nil
	}

	return "", nil
}

func getHeight() uint64 {
	height := uint64(0)

	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bh, _, err := cl.Blocks.Height(ctx)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	} else {
		height = bh.Height
	}

	return height
}

func getCallerInfo() (info string) {

	// pc, file, lineNo, ok := runtime.Caller(2)
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	// funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("%s:%d: ", fileName, lineNo)
}

// func getStats() (*StatsResponse, error) {
// 	client := http.Client{
// 		Timeout: 30 * time.Second,
// 	}

// 	resp, err := client.Get("http://localhost:5001/stats")
// 	if err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 		resp.Body.Close()
// 		client.CloseIdleConnections()
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 		return nil, err
// 	}

// 	var result StatsResponse
// 	if err := json.Unmarshal(body, &result); err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 		return nil, err
// 	}

// 	resp, err = client.Get("http://localhost:5005/distribution")
// 	if err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err = io.ReadAll(resp.Body)

// 	var ds DistributionResponse
// 	if err := json.Unmarshal(body, &ds); err != nil {
// 		log.Println(err)
// 		logTelegram(err.Error())
// 		return nil, err
// 	}

// 	result.Holders = len(ds)

// 	return &result, nil
// }

func getStats() *Stats {
	var miners []*Miner
	sr := &Stats{}
	db.Find(&miners)
	height := getHeight()
	pc := 0

	for _, m := range miners {
		if height-uint64(m.MiningHeight) <= 1440 {
			sr.ActiveMiners++
			if m.ReferralID != 0 {
				sr.ActiveReferred++
			}
		}

		if height-uint64(m.MiningHeight) <= 1440 {
			sr.PayoutMiners++
			pc += int(m.PingCount)

			if hasAintHealth(m, true) {
				sr.ActiveUnits += 10
			} else {
				sr.ActiveUnits++
			}
		}
	}

	sr.InactiveMiners = len(miners) - sr.PayoutMiners
	sr.PingCount = pc

	return sr
}

func getRefCount(m *Miner) uint64 {
	var miners []*Miner

	height := getHeight()

	db.Where("referral_id = ? AND mining_height > ?", m.ID, height-2880).Find(&miners)
	count := len(miners)

	miners = nil

	return uint64(count)
}

type Stats struct {
	ActiveMiners   int `json:"active_miners"`
	ActiveReferred int `json:"active_referred"`
	PayoutMiners   int `json:"payout_miners"`
	InactiveMiners int `json:"inactive_miners"`
	PingCount      int `json:"ping_count"`
	ActiveUnits    int `json:"active_units"`
}

func hasAintHealth(m *Miner, second bool) bool {
	sma := StakeMobileAddress

	d, err := getData("%s__"+m.Address, &sma)
	if err != nil || d == nil {
		return false
	}

	aint := parseItem(d.(string), 0)
	if aint != nil {
		if second && aint.(int) >= (10*MULTI8) {
			return true
		} else if !second && aint.(int) >= MULTI8 {
			return true
		}
	}

	return false
}

type StatsResponse struct {
	ActiveMiners   int `json:"active_miners"`
	ActiveReferred int `json:"active_referred"`
	PayoutMiners   int `json:"payout_miners"`
	InactiveMiners int `json:"inactive_miners"`
	Holders        int `json:"holders"`
	ActiveUnits    int `json:"active_units"`
}

func parseItem(value string, index int) interface{} {
	values := strings.Split(value, Sep)
	var val interface{}
	types := strings.Split(values[0], "%")

	if index < len(values)-1 {
		val = values[index+1]
	}

	if val != nil && types[index+1] == "d" {
		intval, err := strconv.Atoi(val.(string))
		if err != nil {
			log.Println(err.Error())
			logTelegram(err.Error())
		}
		val = intval
	}

	return val
}

func updateItem(value string, newval interface{}, index int) string {
	values := strings.Split(value, Sep)
	types := strings.Split(values[0], "%")

	if index < len(values)-1 {
		switch newval.(type) {
		case int:
			values[index+1] = strconv.Itoa(newval.(int))
		default:
			values[index+1] = newval.(string)
		}
	} else if index < len(types)-1 {
		switch newval.(type) {
		case int:
			values = append(values, strconv.Itoa(newval.(int)))
		default:
			values = append(values, newval.(string))
		}
	}

	return strings.Join(values, Sep)
}

func getHoldersCount() uint64 {
	count := 0
	return uint64(count)
}

type DistributionResponse []struct {
	Address      string  `json:"address"`
	Balance      int64   `json:"balance"`
	BalanceFloat float64 `json:"balance_float"`
}

func saveTelegram(addr string, tids string) int {
	// resp, err := http.Get("http://localhost:5001/save-telegram/" + addr + "/" + tid)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// 	return 5
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)

	// var result SaveTelegramResponse
	// if err := json.Unmarshal(body, &result); err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// 	return 6
	// }

	// return result.Error

	tid, err := strconv.Atoi(tids)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return 1
	}

	m := &Miner{}
	db.First(m, &Miner{TelegramId: int64(tid)})

	if m.ID == 0 {
		db.FirstOrCreate(m, &Miner{TelegramId: int64(tid), Address: tids})
	}

	if strings.HasPrefix(addr, "3A") {
		m.Address = addr
	} else {
		refid, err := strconv.Atoi(addr)
		if err == nil && m.ReferralID == 0 && m.ID != uint(refid) {
			// result = db.FirstOrCreate(m, &Miner{TelegramId: int64(tid), Address: tids, ReferralID: uint(refid)})
			m.ReferralID = uint(refid)
		}
	}

	err = db.Save(m).Error
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return 3
	}

	return 0
}

type SaveTelegramResponse struct {
	Success bool `json:"success"`
	Error   int  `json:"error"`
}

func telegramMine(code string, tid int64) string {
	adnum, err := getData2("%s__adnum", nil)
	miner := getMiner(tid)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	mNotCode := fmt.Sprintf("This code is not valid, it should be 3 numbers.\n\nYou can see the daily mining code <a href=\"https://t.me/AnoteAds/%d\">here</a>.", adnum.(int64))
	mWrongCode := fmt.Sprintf("This code is not correct.\n\nYou can see the daily mining code <a href=\"https://t.me/AnoteAds/%d\">here</a>.", adnum.(int64))
	mSuccess := "You successfully started your Anote mining cycle. 🚀\n\nReward is sent each day when you enter daily mining code."
	mAlreadyMining := "You miner is already mining. You will get notified when you need to repeat the mining cycle."

	if (miner.MiningHeight == 0 || miner.MinedTelegram > 0) && !strings.HasPrefix(miner.Address, "3A") {
		mSuccess += "\n\nTo collect your mining reward, open app.anotedao.com and click 'Connect Telegram' button on the bottom."
	}

	if !isFollower(tid) {
		mSuccess += "\n\nSubscribe to both @AnoteAds and @AnoteNews channels to receive 10% bigger reward!"
	}

	// bot.ChatMemberOf()

	if len(code) != 3 {
		return mNotCode
	}

	codeInt, err := strconv.Atoi(code)
	if err != nil || codeInt > 999 {
		return mNotCode
	}

	resp, err := http.Get("http://localhost:5001/telegram-mine/" + strconv.Itoa(int(tid)) + "/" + code)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()

	for err != nil {
		time.Sleep(time.Millisecond * 500)
		resp, err := http.Get("http://localhost:5001/telegram-mine/" + strconv.Itoa(int(tid)) + "/" + code)
		if err != nil {
			log.Println(err)
			logTelegram(err.Error())
		}
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	var result MineResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if !result.Success {
		if result.Error == 3 {
			return mWrongCode
		}
		if result.Error == 4 {
			return mAlreadyMining
		}
		return "test"
	}

	return mSuccess
}

type MineResponse struct {
	Success bool `json:"success"`
	Error   int  `json:"error"`
}

type MinerResponse struct {
	ID            uint   `json:"id"`
	Address       string `json:"address"`
	Referred      int64  `json:"referred"`
	Active        int64  `json:"active"`
	HasTelegram   bool   `json:"has_telegram"`
	MiningHeight  int64  `json:"mining_height"`
	Height        uint64 `json:"height"`
	Exists        bool   `json:"exists"`
	MinedMobile   uint64 `json:"mined_mobile"`
	MinedTelegram uint64 `json:"mined_telegram"`
	TelegramId    int64  `json:"telegram_id"`
}

func getMiner(tid int64) *MinerResponse {
	mr := &MinerResponse{}
	var miners []*Miner
	height := getHeight()

	// resp, err := http.Get("http://localhost:5001/tminer/" + strconv.Itoa(int(tid)))
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	// if err := json.Unmarshal(body, mr); err != nil {
	// 	log.Println(err)
	// 	logTelegram(err.Error())
	// }

	u := getMinerTel(tid)

	if u.ID != 0 {
		mr.ID = u.ID
		mr.Exists = true
		mr.Address = u.Address
		mr.MiningHeight = u.MiningHeight
		mr.Height = height
		mr.MinedMobile = u.MinedMobile
		mr.MinedTelegram = u.MinedTelegram
		mr.TelegramId = u.TelegramId

		if u.TelegramId != 0 {
			mr.HasTelegram = true
		}

		db.Where("referral_id = ?", u.ID).Find(&miners).Count(&mr.Referred)

		db.Where("referral_id = ? AND mining_height > ?", u.ID, height-2880).Find(&miners).Count(&mr.Active)
	}

	return mr
}

func withdraw(tid int64) *MineResponse {
	mr := &MineResponse{}

	resp, err := http.Get("http://localhost:5001/withdraw/" + strconv.Itoa(int(tid)))
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err := json.Unmarshal(body, mr); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	return mr
}

func getAlphaBalance(address string) uint64 {
	balance := uint64(0)

	ad := &AlphaDistribution{}

	resp, err := http.Get("https://static.anote.digital/alpha-distribution.json")
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	if err := json.Unmarshal(body, ad); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}

	for _, b := range *ad {
		if b.Address == address {
			balance = uint64(b.Balance)
		}
	}

	return balance
}

type AlphaDistribution []struct {
	Address      string  `json:"address"`
	Balance      int64   `json:"balance"`
	BalanceFloat float64 `json:"balance_float"`
}

func sendAsset(amount uint64, assetId string, recipient string) error {
	var networkByte byte
	var nodeURL string

	networkByte = 55
	nodeURL = AnoteNodeURL

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(conf.PublicKeyAlpha)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(conf.PrivateKeyAlpha)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	asset, err := proto.NewOptionalAssetFromString(assetId)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	assetW, err := proto.NewOptionalAssetFromString("")
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	rec, err := proto.NewAddressFromString(recipient)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	tr := proto.NewUnsignedTransferWithSig(sender, *asset, *assetW, uint64(ts), amount, Fee, proto.Recipient{Address: &rec}, nil)

	err = tr.Sign(networkByte, sk)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return err
	}

	return nil
}

func isFollower(uid int64) bool {
	u, err := bot.ChatByID(uid)
	if err != nil {
		log.Println(err)
		return false
	}

	ct, err := bot.ChatByID(TelAnoteToday)
	if err != nil {
		log.Println(err)
		return false
	}

	cn, err := bot.ChatByID(TelAnoteNews)
	if err != nil {
		log.Println(err)
		return false
	}

	cm1, err := bot.ChatMemberOf(ct, u)
	if err != nil {
		log.Println(err)
		return false
	}

	cm2, err := bot.ChatMemberOf(cn, u)
	if err != nil {
		log.Println(err)
		return false
	}

	if (cm1.Role == "member" ||
		cm1.Role == "administrator") &&
		(cm2.Role == "member" ||
			cm2.Role == "administrator") {
		return true
	}

	return false
}

func getAmountNode() float64 {
	var am float64

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: AnoteNodeURL, Client: &http.Client{}, ApiKey: " "})
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return 0
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pc, _, err := cl.Peers.Connected(ctx)
	if err != nil {
		log.Println(err)
		logTelegram(err.Error())
		return 0
	}

	am = (1440 * 0.01) / float64(len(pc))

	return am
}

func newMessage(m *telebot.Message) {
	m3 = m2
	m2 = m1
	m1 = m
}
