package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type MainInfo struct {
	Name string `json:"name"`
	Phone string `json:"phone"`
	Id string `json:"id"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "createMainInfo" {
		//개인정보 생성
		return s.createMainInfo(APIstub, args)
	} else if function == "getAllMainInfo" {
		//원장의 모든 정보 가져오기
		return s.getAllMainInfo(APIstub)
	} else if function == "getMainInfoByIdentifier" {
		//식별자로 정보 가져오기
		return s.getMainInfoByIdentifier(APIstub, args)
	} else if function == "queryMainInfoByName" {
		//이름으로 정보가져오기
		return s.queryMainInfoByName(APIstub, args)
	} else if function == "queryMainInfoByPhone" {
		//연락처로 정보 가져오기
		return s.queryMainInfoByPhone(APIstub, args)
	} else if function == "queryMainInfoById" {
		//아이디로 정보 가져오기
		return s.queryMainInfoById(APIstub, args)
	} else if function == "queryMainInfoByQueryString" {
		//쿼리로 정보 가져오기
		return s.queryMainInfoByQueryString(APIstub, args)
	} else if  function == "getHistoryMainInfo" {
		//정보 이력 가져오기
		return s.getHistoryMainInfo(APIstub, args)
	} else if function == "updateMainInfo" {
		//정보 수정하기
		return s.updateMainInfo(APIstub, args)
	} else if function == "deleteMainInfo" {
		//정보 삭제하기
		return s.deleteMainInfo(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name. ")
}

// 개인정보 생성 함수
func (s *SmartContract) createMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//같은 식별자로 등록되어 있는것이 있는지 확인한다.
	resultsAsBytes, _ := APIstub.GetState(args[0])
	if resultsAsBytes != nil {
		return shim.Error("Already exists!!!")
	}
	
	var mainInfo = MainInfo{Name: args[1], Phone: args[2], Id: args[3]}
	mainInfoAsBytes, _ := json.Marshal(mainInfo)
	APIstub.PutState(args[0], mainInfoAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getAllMainInfo(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	//json으로 이쁘게 변환함
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(buffer.Bytes())
}

//식별자로 정보 가져오는 함수
func (s *SmartContract) getMainInfoByIdentifier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Excepting 1")
	}

	mainInfoAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(mainInfoAsBytes)
}

//이름으로 정보 가져오는 함수
func (s *SmartContract) queryMainInfoByName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"name\":\"%s\"}}", name)
	
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//연락처로 정보 가져오는 함수
func (s *SmartContract) queryMainInfoByPhone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, Excepting 1")
	}

	phone := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"phone\":\"%s\"}}", phone)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//아이디로 정보 가져오는 함수
func (s *SmartContract) queryMainInfoById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//쿼리로 정보 가져오기
func (s *SmartContract) queryMainInfoByQueryString(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := strings.ToLower(args[0])

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//식별자에 해당하는 정보의 이력 가져오기
func (s *SmartContract) getHistoryMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of argument. Expecting 1")
	}

	identifier := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(identifier)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//정보 수정을 위한 함수
func (s *SmartContract) updateMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	mainInfoAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if mainInfoAsBytes == nil {
		return shim.Error("Info does not exist")
	}

	mainInfo := MainInfo{}
	json.Unmarshal(mainInfoAsBytes, &mainInfo)

	if args[1] != "" {
		mainInfo.Name = args[1]
	}

	if args[2] != "" {
		mainInfo.Phone = args[2]
	}

	if args[3] != "" {
		mainInfo.Id = args[3]
	}

	mainInfoAsBytes, _ = json.Marshal(mainInfo)
	APIstub.PutState(args[0], mainInfoAsBytes)

	return shim.Success(nil)
}

//정보를 삭제하는 함수
func (s *SmartContract) deleteMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
	var jsonResp string
	var mainInfoJSON MainInfo
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	identifier := args[0]

	valAsBytes, err := APIstub.GetState(identifier)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + identifier + "\"}"
		return shim.Error(jsonResp)
	} else if valAsBytes == nil {
		jsonResp = "{\"Error\":\"identifier does not exist: " + identifier + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsBytes), &mainInfoJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + identifier + "\"}"
		return shim.Error(jsonResp)
	}

	err = APIstub.DelState(identifier)
	if err != nil {
		return shim.Error("Failed to delete state : " + err.Error())
	}

	return shim.Success(nil)
}

//iterator를 json으로 이쁘게 변환하기 위한 함수
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	return &buffer, nil
}

//couchDB에 쿼리 날리고 결과값 받아오는 함수
func getQueryResultForQueryString (APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	//쿼리 날림
	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	//결과값을 json으로 이쁘게 변환
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

//Iterator size 재는 함수인데 식별자를 키로 쓰면 쓸필요 없음 
func getIteratorSize(iterator shim.StateQueryIteratorInterface) int {
	result := 0
	for iterator.HasNext() {
		iterator.Next()
		result++
	}
	return result
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}