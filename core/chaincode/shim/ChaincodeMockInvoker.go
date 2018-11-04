package shim

import (
	"fmt"
	"flag"
	"encoding/json"
	"strings"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ChaincodeMockInvoker struct{

	stub *MockStub;

}

func (invoker *ChaincodeMockInvoker) Start(ci interface{}) {
	cc := ci.(Chaincode);
	invoker.stub = NewMockStub("ex02", cc);
}

func (invoker *ChaincodeMockInvoker) Invoke() error{

	//argStr := "[\"init\",\"A\",\"567\",\"B\",\"678\"]";
	//argStr := `["init","A","567","B","678"]`;


	argStrOrg := flag.String("invokeargs", "", "passed from web client")
	flag.Parse()

	var argStr string = strings.Replace(*argStrOrg,"'","\"",-1);

	var argStrArray []string;
	json.Unmarshal([]byte(argStr),&argStrArray);


	//args := [][]byte{[]byte("init"), []byte("A"), []byte("567"), []byte("B"), []byte("678")};

	args := convertByte(argStrArray);

	var res pb.Response;
	switch argStrArray[0]{
		case "init" :
			res = invoker.stub.MockInit("1", args)
		case "invoke":
			invoker.stub.MockTransactionStart("1");
			invoker.stub.PutState("A",[]byte("567"));
			invoker.stub.PutState("B",[]byte("678"));
			invoker.stub.MockTransactionEnd("1")
			res = invoker.stub.MockInvoke("1", args)
		default:
	}
	if res.Status != OK {
		fmt.Println(argStrArray[0]," failed", string(res.Message))
	}

	fmt.Println(invoker.stub.State)

	return nil;
}

func MockStart(cc Chaincode){

	invorker := new (ChaincodeMockInvoker);
	invorker.Start(cc);
	invorker.Invoke();


}

func convertByte(str []string)([][]byte){
	var bt [][]byte = make([][]byte,len(str));

	for  i := 0; i < len(str); i++ {
		bt[i] = []byte(str[i]);
	}
	//fmt.Println(str)
	return bt;
}

