package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type MjmallProduct struct {
	Name string `json:"name"`
	Qty  string    `json:"qty"`
	Owner string `json:"owner"`
}

type ProductKey struct {
	Key string
	Idx int
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "registProducts" {
		return t.registProducts(stub, args)
	}else if function == "getProductList" {
		return t.getProductList(stub, args)
	}else if function == "transferOwner" {
		return t.transferOwner(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) registProducts(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	keyAsBytes,err := stub.GetState("latestKey")
	if err != nil {
		return shim.Error(err.Error())
	}
	productKey := ProductKey{}
	json.Unmarshal(keyAsBytes, &productKey)
	var tempIdx string;
	tempIdx = strconv.Itoa(productKey.Idx)
	fmt.Println("Last ProductKey is " + productKey.Key + " : " + tempIdx)

	var product = MjmallProduct{Name: args[0] , Qty: args[1], Owner:args[2]}

	productAsBytes, _ := json.Marshal(product)
	var keyIdx int;
	keyIdx = productKey.Idx
	keyIdx++
	newKeyIdx := strconv.Itoa(keyIdx)
	var newProductKey = productKey.Key + newKeyIdx 
	stub.PutState(newProductKey, productAsBytes)

	productKey.Idx = keyIdx
	newProductKeyAsBytes,_ := json.Marshal(productKey)
	stub.PutState("latestkey", newProductKeyAsBytes)

	return shim.Success(nil)
}

func (t *SimpleChaincode) getProductList(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	keyAsBytes,_ := stub.GetState("latestKey")
	productKey := ProductKey{}
	json.Unmarshal(keyAsBytes, &productKey)

	idxStr := strconv.Itoa(productKey.Idx)

	var startKey = "PD0"
	var endKey = productKey.Key + idxStr

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) transferOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	productAsBytes, _ := stub.GetState(args[0])
	product := MjmallProduct{}
	json.Unmarshal(productAsBytes, &product)

	product.Owner = args[1]

	newProductAsBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(args[0], newProductAsBytes)

	return shim.Success(nil)
}
