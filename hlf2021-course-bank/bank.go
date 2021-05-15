package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BankCC struct {
}

var functions = map[string]func(args []string, stub shim.ChaincodeStubInterface) peer.Response{
	"addAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		// TODO: add validity checks

		personID := args[0]
		response := stub.InvokeChaincode(
			"personCC",
			[][]byte{[]byte("getPerson"), []byte(personID)},
			stub.GetChannelID(),
		)
		if response.Status == shim.ERROR {
			return shim.Error("person doesn't exist cannot create bank account")
		}

		if response.Payload == nil {
			return shim.Error("person doesn't exist cannot create bank account")
		}

		// TODO: add new account

		return shim.Success(nil)
	},
	"deleteAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
	"getAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
	"accountHistory": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
}

func (b *BankCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	panic("implement me")
}

func (b *BankCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	panic("implement me")
}
