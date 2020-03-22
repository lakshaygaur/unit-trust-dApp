package main

import (
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"bytes"
)

// Create Fund and make it available in the network
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
	compositeKey, err := stub.CreateCompositeKey(TRANSACTION_HISTORY,[]string{stub.GetTxID(),stub.GetTxID()})
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

// Read a fund's info with TxnHistory
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

	iterator, err := stub.GetStateByPartialCompositeKey(TRANSACTION_HISTORY,[]string{args[0]})
	if err != nil {
		return shim.Error("Couldn't Create Result iterator"+err.Error())
	}
	var buffer bytes.Buffer
	var firstRow = true
	
	buffer.WriteString(`{ "fund" : `)
	buffer.Write(fundAsBytes)
	buffer.WriteString(`, "txnHistory":[`)

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if firstRow != true {
			buffer.WriteString(",")
		}
		assetAddress := string(queryResponse.Key)
		assetAsBytes, err :=stub.GetState(assetAddress)
		if err != nil {
			return shim.Error("Failed to get state of asset : address "+ assetAddress+" err: "+err.Error())
		}
		buffer.Write(assetAsBytes)
		firstRow = false
	}
	buffer.WriteString(`]}`)

	return shim.Success(buffer.Bytes())
}

// Create HQ account
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
	account.Status = true
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

// Create HQ account
func (t *UnitTrustChaincode) ApproveAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// args check : should be name
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1, got "+ strconv.Itoa(len(args)))
	}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return shim.Error("Error Parsing Certificate : "+ err.Error())
	}
	compositeKey, err := stub.CreateCompositeKey(ACCOUNT,[]string{AGENT,cert.Subject.CommonName})

	accountAsBytes, err := stub.GetState(compositeKey)
	if err != nil {
		return shim.Error("Failed to get state account : "+ err.Error())
	}
	account := Account{}
	err = json.Unmarshal( accountAsBytes, &account)
	if err != nil {
		return shim.Error("Failed to unmarshal account : "+ err.Error())
	}

	account.Status = true // account approved
	
	accountAsBytes, err = json.Marshal(account)
	if err != nil {
		return shim.Error("Failed to marshal account : "+err.Error())
	}
	
	// put state to ledger
	err = stub.PutState(compositeKey, accountAsBytes)
	if err != nil {
		return shim.Error("Failed to put state account : "+ err.Error())
	}
	return shim.Success(nil)
}


// Deletes a fund from state
func (t *UnitTrustChaincode) DeleteFund(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// args check : should be fundId
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1, got "+ strconv.Itoa(len(args)))
	}

	// Delete the key from the state in ledger
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}