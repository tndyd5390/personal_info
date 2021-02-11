package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	_"strconv"
	_"time"
	_"reflect"
	"crypto/sha256"
	"encoding/hex"


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
		return s.createMainInfo(APIstub, args)
	} else if function == "deleteMainInfo"{
		return s.deleteMainInfo(APIstub, args)
	} else if function == "modificateMainInfo"{
		return s.modificateMainInfo(APIstub, args)
	} else if function == "getAllMainInfo" {
		return s.getAllMainInfo(APIstub)
	} else if function == "getMainInfoByIdentifier" {
		return s.getMainInfoByIdentifier(APIstub, args)
	} else if function == "queryMainInfoByName" {
		return s.queryMainInfoByName(APIstub, args)
	} else if function == "queryMainInfoByPhone" {
		return s.queryMainInfoByPhone(APIstub, args)
	} else if function == "queryMainInfoById" {
		return s.queryMainInfoById(APIstub, args)
	} else if function == "queryMainInfoByQueryString" {
		return s.queryMainInfoByQueryString(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name. ")
}

// 개인정보 생성 함수
func (s *SmartContract) createMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	name := args[0]
	phone := args[1]
	id := args[2]
	identifier := makeIdentifier(&name, &phone, &id)

	//여기서 identifier 중복 체크
	valAsByte, _ := APIstub.GetState(identifier)
	if valAsByte != nil {
		return shim.Error("{\"Error\":\"Already exist!!!\"}")
	}

	var mainInfo = MainInfo{Name: args[0], Phone: args[1], Id: args[2]}
	mainInfoAsBytes, _ := json.Marshal(mainInfo)
	APIstub.PutState(identifier, mainInfoAsBytes)

	return shim.Success(nil)
}
func makeIdentifier(name *string, phone *string, id *string) string {
        sha := sha256.New()
        sha.Write([]byte(*name))
        sha.Write([]byte(*phone))
        sha.Write([]byte(*id))
        identifier := sha.Sum(nil)
        identifier2 := hex.EncodeToString(identifier)
        return identifier2
}

func (s *SmartContract) deleteMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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

func (s *SmartContract) modificateMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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


func (s *SmartContract) getAllMainInfo(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

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
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
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

	fmt.Printf("- getQueryResultForQueryString queryResult: \n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

/*func (s *SmartContract) getHistoryMainInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

        if len(args) < 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
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
                // Add a comma before array members, suppress it for the first array member
                if bArrayMemberAlreadyWritten == true {
                        buffer.WriteString(",")
                }
                buffer.WriteString("{\"TxId\":")
                buffer.WriteString("\"")
                buffer.WriteString(response.TxId)
                buffer.WriteString("\"")

                buffer.WriteString(", \"Value\":")
                // if it was a delete operation on given key, then we need to set the
                //corresponding value null. Else, we will write the response.Value
                //as-is (as the Value itself a JSON marble)
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

        fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

        return shim.Success(buffer.Bytes())
}*/



func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
