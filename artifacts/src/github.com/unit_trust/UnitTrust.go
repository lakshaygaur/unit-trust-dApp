
package main


import (
	"fmt"
	// "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("UnitTrustChaincode")


type UnitTrustChaincode struct {
}

func (t *UnitTrustChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### UnitTrustChaincode Init ###########")

	// Initialize the chaincode

	return shim.Success(nil)
}

// List of invoke functions
func (t *UnitTrustChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### UnitTrustChaincode Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	
	if function == "deleteFund" {
		return t.DeleteFund(stub, args)
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
	logger.Errorf("Unknown action, check the first argument, must be one of 'deleteFund', 'createFund','readFund','createAccount', 'applyAccount','sellFund', 'readAllFunds' or 'approveAccount'. But got: %v", function)
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'deleteFund', 'createFund','readFund','createAccount', 'applyAccount','sellFund', 'readAllFunds' or 'approveAccount'. But got: %v", function))
}



func main() {
	err := shim.Start(new(UnitTrustChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
