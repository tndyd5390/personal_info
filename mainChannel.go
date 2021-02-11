package main

import {
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
}

type SmartContract struct {
}

type BasicInfo struct {
	Identifier string `json:"identifier"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Id string `json:"id"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response{
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	
	function, args := APIstub.GetFunctionAndParameters()

	if function == "createBasicInfo"{
		return s.CreateBasicInfo(APIstub, args)
	} else if function == "queryAllBasicInfo" {
		return s.QueryAllBasicInfo(APIstub)
	} else if function == "queryBasicInfoByKeyValue" {
		return s.QueryBasicInfoByKeyValue(APIstub, args)
	}
}

func (s *SmartContract) CreateBasicInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments.Expecting 5")
	}
	
	var basicInfo = BasicInfo{Identifier: args[1], Name: args[2], Phone: args[3], Id: args[4]}

	basicInfoAsBytes, _ := json.Marshal(basicInfo)
	APIstub.PutState(args[0], basicInfoAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) QueryAllBasicInfo(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) QueryBasicInfoByKeyValue(APIstub ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	key := strings.ToLower(args[0])
	key := strings.ToLower(args[1])
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}", key, value)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.bufferbuffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}