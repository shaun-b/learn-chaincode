/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	//write the specified `key` and `value` into the ledger: i.e. the first argument sent in the deployment
	err := stub.PutState("hello world",[]byte(args[0]))


	return nil, err
}

// Invoke is our entry point to invoke a chaincode function.
// Invoke functions are captured as transactions, which get grouped into blocks for writing to the ledger.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {	//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub,args)
	}

	fmt.Println("invoke did not find func: " + function)//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {

	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. write() expecting 2: name of variable and value to set")
	}

	var name, value string
	var err error

	name = args[0]
	value = args[1]
	//write the variable into the chaincode state (i.e. the ledger)
	err = stub.PutState(name, []byte(value))

	return nil, err
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error)  {

	fmt.Println("running read()")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. read() is expecting the key name to query")
	}

	var name, jsonResp string
	var err error

	name = args[0]
	valAsBytes, err := stub.GetState(name)

	if err != nil {
		jsonResp = "{\"error\":\"failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsBytes, nil
}