package main

import (
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

  
func (t *UnitTrustChaincode) CreateFund(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	// args check : should be type, value, validFrom, validTo
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, got "+ strconv.Itoa(len(args)))
	}

	account, err := getUser(stub)
	if err != nil {
		return shim.Error("Failed to get user: "+err.Error())
	}
	
	// create fund object
	fund := Fund{}
	fundId, err := stub.CreateCompositeKey(FUND,[]string{stub.GetTxID()}) 
	if err !=nil {
		return shim.Error("Failed to generate Fund ID: "+ err.Error())
	}
	fund.FundId = fundId
	fund.Type = args[0]
	fund.Value = args[1]
	fund.ValidFrom = args[2]
	fund.ValidTo = args[3]	
	fund.Owner = account.AccountId

	fundAsBytes, err := json.Marshal(fund) // convert fund object to bytes

	// push state to ledger 
	err = stub.PutState(fundId, fundAsBytes) // Write the state back to the ledger
	if err !=nil {
		return shim.Error("Failed to put state fund: "+ err.Error()+" fund Id "+ fundId)
	}

	// update txn history 
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get timestamp: "+ err.Error())
	}
	compositeKey, err := stub.CreateCompositeKey(TRANSACTION_HISTORY,[]string{stub.GetTxID()})
	if err != nil {
		return shim.Error("Failed to generate TRANSACTION_HISTORY ID: "+ err.Error())
	}
	txnHistory := TransactionHistory{}
	txnHistory.Status = "Fund: "+ fund.Type +" created by "+ account.AccountId
	txnHistory.Timestamp = timestamp.String()
	txnHistory.TxnId = compositeKey

	txnAsBytes, err := json.Marshal(txnHistory)
	err = stub.PutState(compositeKey, txnAsBytes) // Write the state back to the ledger
	if err !=nil {
		return shim.Error("Failed to put state txn history: "+ err.Error()+ " txn id "+ compositeKey)
	}

	return shim.Success(nil)
}

func (t *UnitTrustChaincode) ReadFund(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	
	// args check : should be fundId
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1, got "+ strconv.Itoa(len(args)))
	}

	// fund := Fund{}
	fundId, err := stub.CreateCompositeKey(FUND,[]string{args[0]}) 
	if err != nil {
		return shim.Error("Failed to get fund : "+err.Error())
	}
	fundAsBytes, err := stub.GetState(fundId)

	// err = json.Unmarshal(fundAsBytes, &fund)

	return shim.Success(fundAsBytes)
}

func (t *UnitTrustChaincode) CreateAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// args check : should be name, type
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, got "+ strconv.Itoa(len(args)))
	}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return shim.Error("Error Parsing Certificate : "+ err.Error())
	}
	compositeKey, err := stub.CreateCompositeKey(ACCOUNT,[]string{cert.Subject.CommonName})

	account := Account{}
	account.Name = args[0]
	account.Type = args[1]
	account.AccountId = cert.Subject.CommonName

	accountAsBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error("Failed to marshal account : "+err.Error())
	}
	
	
	err = stub.PutState(compositeKey, accountAsBytes)
	if err != nil {
		return shim.Error("Failed to put state account : "+ err.Error())
	}
	return shim.Success(nil)
}