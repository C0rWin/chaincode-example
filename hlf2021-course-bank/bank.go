package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BankCC struct {
}

type BankAccount struct {
	PersonID      string  `json:person_id`
	AccountID 	  string  `json:account_number`
	Balance       float64 `json:balance`
}

type Transfer struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Value float64 `json:"value"`
}

var functions = map[string]func(args []string, stub shim.ChaincodeStubInterface) peer.Response{
	"addAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error("invalid add person chaincode invocation")
		}

		var account BankAccount
		
		if err := json.Unmarshal([]byte(args[0]), &account); err != nil {
			return shim.Error(fmt.Sprintf("error unmarshalling account %s", err))
		}

		personID := account.PersonID
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

		accountState, err := stub.GetState(account.AccountID)

		if err != nil {
			return shim.Error(fmt.Sprintf("failed to create bank account %s", err))
		}

		if accountState != nil {
			return shim.Error(fmt.Sprintf("bank account %s already exists", account.AccountID))
		}

		if err := stub.PutState(account.AccountID, []byte(args[0])); err != nil {
			shim.Error(fmt.Sprintf("failed to save bank account %s : %s", account.AccountID, err))
		}

		return shim.Success(nil)
	},
	"deleteAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of parameters"))
		}

		err := stub.DelState(args[0])
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to delete account information, id %s, due to %s", args[0], err))
		}

		return shim.Success(nil)
	},
	"getAccount": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of arguments"))
		}
		
		accountState, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(fmt.Sprintf("failed to get account %s state due to %s", args[0], err))
		}

		if accountState == nil {
			return shim.Error(fmt.Sprintf("Account with number %s doesnt exist", args[0]))
		}

		return shim.Success(accountState)
	},
	"getBalance": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {

		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of arguments"))
		}

		state, err := stub.GetState(args[0])

		if err != nil {
			return shim.Error(err)
		}

		if state == nil {
			return shim.Error("bank account doesn't exists")
		}

		var account BankAccount
		if err = json.Unmarshal([]byte(state), &account); err != nil {
			return shim.Error(err)
		}

		return shim.Success([]byte(account.Balance))
	},
	"transfer": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		if len(args) != 1 {
			return shim.Error(fmt.Sprintf("wrong number of parameters"))
		}

		var transfer Transfer

		err := json.Unmarshal([]byte(args[0]), &transfer)
		if err != nil {
			return shim.Error(err)
		}

		from, err := stub.GetState(transfer.From)

		if err != nil {
			return shim.Error(err)
		}

		if from == nil {
			return shim.Error(fmt.Sprintf("from account %s doesn't exist", transfer.From))
		}

		var fromAccount BankAccount
		if err = json.Unmarshal([]byte(from), &fromAccount); err != nil {
			return shim.Error(err)
		}

		to, err := stub.GetState(transfer.To)

		if err != nil {
			return shim.Error(err)
		}

		if to == nil {
			return shim.Error(fmt.Sprintf("to account %s doesn't exist", transfer.To))
		}

		var toAccount BankAccount
		if err = json.Unmarshal([]byte(to), &toAccount); err != nil {
			return shim.Error(err)
		}

		issueMoneyAccountID := "0"
		if (fromAccount.Balance-transfer.Value) < 0 && fromAccount.AccountNumber != issueMoneyAccountID {
			return shim.Error(fmt.Sprintf("Insufficient balance on 'from' account %s", transfer.From))
		}

		fromAccount.Balance = fromAccount.Balance - transfer.Value
		toAccount.Balance = toAccount.Balance + transfer.Value

		fromString, err := json.Marshal(fromAccount)
		if err != nil {
			return shim.Error(err)
		}

		toString, err := json.Marshal(toAccount)
		if err != nil {
			return shim.Error(err)
		}

		err = stub.PutState(transfer.From, []byte(fromString))
		if err != nil {
			return shim.Error(err)
		}

		err = stub.PutState(transfer.To, []byte(toString))
		if err != nil {
			return shim.Error(err)
		}

		return shim.Success(nil)
	},
	"getHistory": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
		historyIterator, err := stub.GetHistoryForKey(args[0])

		if err != nil {
			return shim.Error(err)
		}

		result := make([]string, 0, 6)

		var historyEntry BankAccount

		prevBalance := 0.0

		for historyIterator.HasNext() {
			modification, err := historyIterator.Next()
			if err != nil {
				return shim.Error(err)
			}
			err = json.Unmarshal(modification.Value, &historyEntry)
			if err != nil {
				return shim.Error(err)
			}

			if math.Abs(historyEntry.Balance-prevBalance) > 0.0001 {
				result = append(result, fmt.Sprintf("%+f:%d", historyEntry.Balance-prevBalance, modification.Timestamp.GetSeconds()))
			}
			prevBalance = historyEntry.Balance
		}

		jsonString, _ := json.Marshal(result)
		return shim.Success([]byte(jsonString))
	},
}

func (b *BankCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Bank Management chaincode is initialized")
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
