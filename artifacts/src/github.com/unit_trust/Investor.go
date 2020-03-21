package main

import (
	"bytes"
	"strconv"
	// "encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Investors read all funds available
func (t *UnitTrustChaincode) ReadAllFunds(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// args check : should be 0
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0, got "+ strconv.Itoa(len(args)))
	}
	
	iterator, err := stub.GetStateByPartialCompositeKey(FUND,[]string{})
	if err != nil {
		return shim.Error("Couldn't Create Result iterator"+err.Error())
	}
	var buffer bytes.Buffer
	var firstRow = true

	buffer.WriteString(`{ "funds" : [`)
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