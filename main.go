package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	//"net/http"
	"os"
	"sync"

	//"strconv"

	bsctoken "BSCTokenList/build"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Tokens struct {
	Tokens []Token `json:"tokens"`
}

type Token struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

func TokenBal(tokenaddress string, wg *sync.WaitGroup, tokenAddress common.Address, tokens Tokens, i int, client *ethclient.Client) {

	instance, err := bsctoken.NewBsctoken(tokenAddress, client)
	defer wg.Done()
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x6877654e79119a7f9A8182BB2389E797CA421D24")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	if len(bal.Bits()) > 0 {
		fmt.Println("\n============================================")
		fmt.Println("Token Name: " + tokens.Tokens[i].Name)
		fmt.Println("Token Symbol: " + tokens.Tokens[i].Symbol)
		fmt.Println("Balance: ", value)
	}
}

func main() {

	// resp, err := http.Get("https://wispy-bird-88a7.uniswap.workers.dev/?url=http://tokens.1inch.eth.link")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //read the respond body on the line
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //write the complete body
	// err = ioutil.WriteFile("tokenlist.json", body, 0644)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	jsonFile, err := os.Open("bsctokenlist.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var tokens Tokens

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &tokens)

	wg := &sync.WaitGroup{}

	client, err := ethclient.Dial("https://speedy-nodes-nyc.moralis.io/2abeb0fd5d1beb3b221a0f05/bsc/mainnet")
	if err != nil {
		log.Fatal(err)
	}
	// we iterate through every Tokens within our Token array and
	// print out the token Name, Symbol, and their Address
	// as just an example
	for i := 0; i < len(tokens.Tokens); i++ {
		tokenaddress := tokens.Tokens[i].Address
		tokenAddress := common.HexToAddress(tokenaddress)
		wg.Add(1)
		go TokenBal(tokenaddress, wg, tokenAddress, tokens, i, client)
	}
	wg.Wait()

}
