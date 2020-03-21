package main

import (
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

// Allow agent to request to join the network
func (t *UnitTrustChaincode) ApplyAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// args check : should be name, type
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, got "+ strconv.Itoa(len(args)))
	}
	if args[1] != AGENT {
		return shim.Error("Account type should be agent only")
	}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return shim.Error("Error Parsing Certificate : "+ err.Error())
	}
	compositeKey, err := stub.CreateCompositeKey(ACCOUNT,[]string{AGENT,cert.Subject.CommonName})
	if err!=nil {
		return shim.Error("Failed to generate agent account key : "+ err.Error())
	}
	account := Account{}
	account.Name = args[0]
	account.Type = args[1]
	account.AccountId = cert.Subject.CommonName
	account.Status = false // un-approved
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

// Agents sell funds to investors
func (t *UnitTrustChaincode) SellFund(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// args check : should be fundId, sellingTo
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2, got "+ strconv.Itoa(len(args)))
	}
	
	_,_ = getAgent(stub) // check status of agent's account
	
	fund := Fund{}
	fundId, err := stub.CreateCompositeKey(FUND,[]string{args[0]})
	if err != nil {
		return shim.Error("Failed to create fundId: "+ err.Error())
	}
	fundAsBytes, err:= stub.GetState(fundId)
	if err != nil {
		return shim.Error("Failed to get fund data: id "+ fundId+" err : "+ err.Error())
	}
	err = json.Unmarshal(fundAsBytes, &fund)
	if err != nil {
		return shim.Error("Failed to unmarshal fund: "+ err.Error())
	}
	
	fund.Owner = args[1] // change owner
	
	// push state to ledger 
	fundAsBytes, err = json.Marshal(fund)
	if err != nil {
		return shim.Error("Failed to marshal fund data: id "+ fundId+" err : "+ err.Error())
	}	
	err = stub.PutState(fundId, fundAsBytes)
	if err !=nil {
		return shim.Error("Failed to put state fund: "+ err.Error()+" fund Id "+ fundId)
	}

	// update txn history 
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Failed to get timestamp: "+ err.Error())
	}
	compositeKey, err := stub.CreateCompositeKey(TRANSACTION_HISTORY,[]string{args[0],stub.GetTxID()})
	if err != nil {
		return shim.Error("Failed to generate TRANSACTION_HISTORY ID: "+ err.Error())
	}
	txnHistory := TransactionHistory{}
	txnHistory.Status = "Fund: "+ fund.Type +" sold to "+ args[1]
	txnHistory.Timestamp = timestamp.String()
	txnHistory.TxnId = compositeKey

	txnAsBytes, err := json.Marshal(txnHistory)
	err = stub.PutState(compositeKey, txnAsBytes) // Write the state back to the ledger
	if err !=nil {
		return shim.Error("Failed to put state txn history: "+ err.Error()+ " txn id "+ compositeKey)
	}

	return shim.Success(nil)
}