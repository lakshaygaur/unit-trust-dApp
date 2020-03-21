
package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"encoding/json"
)

// getUser returns Account of user invoking the transaction.
func getUser(stub shim.ChaincodeStubInterface) (*Account , error) {
	account := &Account{}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return nil, errors.New("Error Parsing Certificate : "+ err.Error())
	}
	compositeKey,err := stub.CreateCompositeKey(ACCOUNT,[]string{cert.Subject.CommonName})
	if err !=nil {
		return nil,errors.New("Error creating composite Key : "+err.Error())
	}
	accountAsBytes,err := stub.GetState(compositeKey)
	if err !=nil {
		return nil,errors.New("Error fetching user : "+ err.Error())
	}
	err = json.Unmarshal(accountAsBytes,account)
	if err != nil {
		return nil, errors.New("Error decoding account : "+ err.Error())
	}
	return account,nil
}


func getAgent(stub shim.ChaincodeStubInterface) (*Account , error) {
	account := &Account{}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return nil, errors.New("Error Parsing Certificate : "+ err.Error())
	}
	compositeKey,err := stub.CreateCompositeKey(AGENT,[]string{cert.Subject.CommonName})
	if err !=nil {
		return nil,errors.New("Error creating composite Key : "+err.Error())
	}
	accountAsBytes,err := stub.GetState(compositeKey)
	if err !=nil {
		return nil,errors.New("Error fetching user : "+ err.Error())
	}
	err = json.Unmarshal(accountAsBytes,account)
	if err != nil {
		return nil, errors.New("Error decoding account : "+ err.Error())
	}
	// account status check
	if account.Status == false {
		return nil, errors.New("Agent Account is not approved yet")
	}
	return account,nil
}