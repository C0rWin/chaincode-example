package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type CreditCardsCC struct {
}

type BankAccount struct {
	PersonID      string  `json:person_id`
	AccountID 	  string  `json:account_number`
	Balance       float64 `json:balance`
}

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PassportID string `json:"passport_id"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

type CreditCard struct {
	CardNumber    string `json:"card_number"`
	ExpireDate    string `json:"expire_date"`
	PersonID      uint64 `json:"person_id"`
	AccountNumber string `json:"account_number"`
}



var functions = map[string]func(args []string, stub shim.ChaincodeStubInterface) peer.Response{
	"addCard": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error("wrong number of input parameters")
		}

		var card CreditCard
		err := json.Unmarshal([]byte(args[0]), &card)
		if err != nil {
			return shim.Error("failed to deserialize", err)
		}

		// Existence of person
		personID := card.PersonID
		response := stub.InvokeChaincode("personCC", [][]byte{[]byte("getPerson"), []byte(personID)}, stub.GetChannelID())
		if response.Status == shim.ERROR {
			return shim.Error("failed to create credit card")
		}

		// Existance of bank account
		response = stub.InvokeChaincode("bankCC", [][]byte{[]byte("getAccount"), []byte(card.AccountNumber)}, stub.GetChannelID())
		if response.Status == shim.ERROR {
			return shim.Error("failed to create credit cart and link it to the bank account")
		}

		var account BankAccount
		err = json.Unmarshal(response.GetPayload(), &account)
		if err != nil {
			return shim.Error("failed to deserialize payload bankCC")
		}

		// Client is owner of the bank account
		if account.PersonID != card.PersonID {
			return shim.Error("person is not an owner of the bank account"))
		}

		// Existance of the card number and duplicate
		state, err := stub.GetState(card.CardNumber)
		if err != nil {
			return shim.Error(err)
		}

		if state != nil {
			return shim.Error("card already exist")
		}

		jsonString, _ := json.Marshal(card)
		err = stub.PutState(card.CardNumber, jsonString)
		if err != nil {
			return shim.Error(err)
		}

		return shim.Success([]byte(card.CardNumber))
	},
	"deleteAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of parameters"))
		}

		err := stub.DelState(args[0])
		if err != nil {
			return shim.Error(err)
		}

		return shim.Success(nil)
	},
	"getCard": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of arguments"))
		}
		
		cardState, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err)
		}

		if cardState == nil {
			return shim.Error("Card doesnt exist")
		}

		return shim.Success(cardState)
	},
}

func (b *CreditCardsCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Bank Management chaincode is initialized")
	return shim.Success(nil)
}

func (b *CreditCardsCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functionName, args := stub.GetFunctionAndParameters()

	f, ok := functions[functionName]
	if !ok {
		return shim.Error("unknown function name for chaincode CreditCardCC")
	}

	return f(args, stub)
}

func main() {
	err := shim.Start(new(CreditCardsCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
