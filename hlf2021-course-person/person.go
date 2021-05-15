package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type PersonCC struct {
}

type Person struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PassportID string `json:"passport_id"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}

var functions = map[string]func(args []string, stub shim.ChaincodeStubInterface) peer.Response{
	"addPerson": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error("invalid add person chaincode invocation")
		}

		person := &Person{}
		if err := json.Unmarshal([]byte(args[0]), person); err != nil {
			return shim.Error(err.Error())
		}

		pp, err := stub.GetState(person.PassportID)
		if err != nil {
			return shim.Error(err.Error())
		}

		if pp != nil {
			return shim.Error(fmt.Sprintf("persion with id %s already exists", person.PassportID))
		}

		if err := stub.PutState(person.PassportID, []byte(args[0])); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	},
	"getPerson": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
	"deletePerson": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
	"personHistory": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	},
}

func (p *PersonCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("PersonCC has been initialized")
	return shim.Success(nil)
}

func (p *PersonCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	functionName, args := stub.GetFunctionAndParameters()

	f, ok := functions[functionName]
	if !ok {
		return shim.Error("unknown function name for chaincode PersonCC")
	}

	return f(args, stub)
}

func main() {
	err := shim.Start(new(PersonCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
