
package main


import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("UnitTrustChaincode")


type UnitTrustChaincode struct {
}

func (t *UnitTrustChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### UnitTrustChaincode Init ###########")

	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	logger.Info("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)


}

// Transaction makes payment of X units from A to B
func (t *UnitTrustChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### UnitTrustChaincode Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	
	if function == "delete" {
		return t.delete(stub, args)
	}
	if function == "createFund" {
		return t.CreateFund(stub, args)
	}
	if function == "readFund" {
		return t.ReadFund(stub, args)
	}
	if function == "createAccount" {
		return t.CreateAccount(stub, args)
	}
	if function == "applyAccount" {
		return t.ApplyAccount(stub, args)
	}
	if function == "sellFund" {
		return t.SellFund(stub, args)
	}
	if function == "readAllFunds" {
		return t.ReadAllFunds(stub, args)
	}
	if function == "approveAccount" {
		return t.ApproveAccount(stub, args)
	}
	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", function)
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", function))
}


// Deletes an entity from state
func (t *UnitTrustChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}


func main() {
	err := shim.Start(new(UnitTrustChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
