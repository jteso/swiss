package main

import(
	"math"
	"fmt"
	"net/rpc"
	"log"
	"circuitbreaker"
	"circuitbreaker/integration"
	"errors"
)

// ---- API level-----
type errorHandler interface {
	invalidArgument(int) (int, error)
}

func sqrt(i int, eh errorHandler) (float64, error){
	if i < 0 {
		i, err := eh.invalidArgument(i)
		if err != nil {
			return i, err
		}
	}
	return math.Sqrt(float64(i)),nil
}


// Client

type myHandler struct {}

func (_ myHandler) invalidArgument(i int) (int, error){
	return i, nil
}

func main(){
	root, err := sqrt(3, myHandler{})
	if err != nil {
		panic("ohoh")
	}
	fmt.Println(root)
}


// ****************** Test circuit braker

type Args struct {
	a,b int
}
// >> Client uses a struct
type AddCommand struct {}

func (_ AddCommand) Run()(int, error){
	client,err := rpc.DialHTTP("tcp", "localhost" + ":1234")
	if err != nil {
		log.Fatal("Error here")
	}
	var reply int
	err := client.Call("addOperation", Args{7,8}, &reply)

	if err != nil {
		return -1, err
	}
	return reply, nil
}

// >> Client uses a func
func AddCommandFunc()(int, error){
	client, err := rpc.DialHTTP("tcp", "localhost" + ":1234")
	if err != nil {
		log.Fatal("Error here")
	}
	var reply int
	err := client.Call("addOperation", Args{7,8}, &reply)

	if err != nil {
		return -1, err
	}
	return reply, nil
}


// Test

func main() {
	errors.New
	//Default
	cb := circuitbreaker.New("add", 2, 12000)
	cb.Execute(CommandFunc(AddCommandFunc))
	//BYO parameters
	cb2, err := circuitbreaker.Get("add")
	if err != nil {
		cb2.Execute(AddCommand{})
	}

}



