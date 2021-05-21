package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BankCC struct {
}

type Account struct {
	Number string `json:"number"`
	BankName    string `json:"bank_name"`
	INN      string `json:"inn"`
	PersonId string `json:"person_id"`
	CardId string `json:"card_id"`
}

type Card struct {
	Number string `json:"number"`
	Type    string `json:"type"`
	Active string `json:"active"`
}

var functions = map[string]func(args []string, stub shim.ChaincodeStubInterface) peer.Response{
	"addAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		account := &Account{}
		if err := json.Unmarshal([]byte(args[0]), account); err != nil {
			return shim.Error(err.Error())
		}

		pp, err := stub.GetState(account.Number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp != nil {
			return shim.Error(fmt.Sprintf("Account with number %s already exists.", account.Number))
		}

		response := stub.InvokeChaincode(
			"person",
			[][]byte{[]byte("getPerson"), []byte(account.PersonId)},
			stub.GetChannelID(),
		)

		if response.Status == shim.ERROR {
			return shim.Error("PersonCC response error.")
		}

		if response.Payload == nil{
			return shim.Error("Person doesnt exist. Impossible to create bank account.")
		}

		if account.CardId != "" {
			pp, err := stub.GetState(account.CardId)
			if err != nil {
				return shim.Error(err.Error())
			}
			if pp == nil {
				return shim.Error(fmt.Sprintf("Card with number %s doesnt exist. It's impossible to link card.", account.CardId))
			}
		}

		if err := stub.PutState(account.Number, []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"getAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		number := args[0]
		pp, err := stub.GetState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Account with number %s doesnt exist.", number))
		}

		return shim.Success(pp)
	},
	"updateAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		account := &Account{}
		if err := json.Unmarshal([]byte(args[0]), account); err != nil {
			return shim.Error(err.Error())
		}

		pp, err := stub.GetState(account.Number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Account with number %s doesnt exist.", account.Number))
		}

		response := stub.InvokeChaincode(
			"person",
			[][]byte{[]byte("getPerson"), []byte(account.PersonId)},
			stub.GetChannelID(),
		)

		if response.Status == shim.ERROR {
			return shim.Error("PersonCC response error.")
		}

		if response.Payload == nil{
			return shim.Error("Person doesnt exist. Impossible to update bank account.")
		}

		if account.CardId != "" {
			pp, err := stub.GetState(account.CardId)
			if err != nil {
				return shim.Error(err.Error())
			}
			if pp == nil {
				return shim.Error(fmt.Sprintf("Card with number %s doesnt exist. It's impossible to link card.", account.CardId))
			}
		}

		if err := stub.PutState(account.Number, []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"deleteAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		number := args[0]
		pp, err := stub.GetState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Account with number %s doesnt exist.", number))
		}

		err = stub.DelState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"accountHistory": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		number := args[0]
		pp, err := stub.GetState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Account with number %s doesnt exist.", number))
		}

		history, err := stub.GetHistoryForKey(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		records := []string{}
		i := 1
		for history.HasNext(){
			record, _ := history.Next()
			var buf = fmt.Sprintf("%d record : %s %s", i, record.GetTxId(), record.GetValue())
			records = append(records, buf)
			i++
		}

		var payload string
		for _, v := range records {
			payload = payload + v
		}
		return shim.Success([]byte(payload))
	},

	"addCard": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		card := &Card{}
		if err := json.Unmarshal([]byte(args[0]), card); err != nil {
			return shim.Error(err.Error())
		}

		pp, err := stub.GetState(card.Number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp != nil {
			return shim.Error(fmt.Sprintf("Card with number %s already exists.", card.Number))
		}

		if err := stub.PutState(card.Number, []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"getCard": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		number := args[0]
		pp, err := stub.GetState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Card with number %s doesnt exist.", number))
		}

		return shim.Success(pp)
	},
	"updateCard": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		card := &Card{}
		if err := json.Unmarshal([]byte(args[0]), card, ); err != nil {
			return shim.Error(err.Error())
		}

		pp, err := stub.GetState(card.Number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Card with number %s doesnt exist.", card.Number))
		}

		if err := stub.PutState(card.Number, []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"cardHistory": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) < 2 {
			return shim.Error("Unsufficient amount of arguments.")
		}

		number := args[0]
		pp, err := stub.GetState(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp == nil {
			return shim.Error(fmt.Sprintf("Card with number %s doesnt exist.", number))
		}

		history, err := stub.GetHistoryForKey(number)
		if err != nil {
			return shim.Error(err.Error())
		}

		records := []string{}
		i := 1
		for history.HasNext(){
			record, _ := history.Next()
			var buf = fmt.Sprintf("%d record : %s %s", i, record.GetTxId(), record.GetValue())
			records = append(records, buf)
			i++
		}

		var payload string
		for _, v := range records {
			payload = payload + v
		}
		return shim.Success([]byte(payload))
	},
}

func (b *BankCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("PersonCC has been initialized")
	return shim.Success(nil)
}

func (b *BankCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functionName, args := stub.GetFunctionAndParameters()

	f, ok := functions[functionName]
	if !ok {
		return shim.Error("unknown function name for chaincode PersonCC")
	}

	return f(args, stub)
}

func main() {
	err := shim.Start(new(BankCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}