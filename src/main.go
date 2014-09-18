package main

import (
	"fmt"
	_ "net/rpc"
	_ "log"
	"strconv"
	"reflect"
	"errors"
)
type C struct {
	name string
}
type Args struct {
	A,B int
	c *C
}

func main() {

	var args = new(Args)

	fmt.Println("args.A:" + strconv.Itoa(args.A))
	fmt.Println("args.B:" + strconv.Itoa(args.B))
	fmt.Println("type of args.A:", reflect.TypeOf(args.A))
	fmt.Println("type of args.B:", reflect.TypeOf(args.B))
	fmt.Println("type of args.C:", reflect.TypeOf(args.c))

	args2 := Args{}
	fmt.Println("args.A:" + strconv.Itoa(args2.A))
	fmt.Println("args.B:" + strconv.Itoa(args2.B))
	fmt.Println("type of args.A:", reflect.TypeOf(args.A))
	fmt.Println("type of args.B:", reflect.TypeOf(args.B))
	fmt.Println("type of args.C:", reflect.TypeOf(args.c))



//	client,err := rpc.DialHTTP("tcp", "localhost" + ":1234")
//	if err != nil {
//		log.Fatal("Error found while getting a client")
//	}
//
//	wrapper := NewClientWrapper(client)
//
//	var reply int
//	wrapper.execute("addOperation", Args{7,8}, &reply)
//	fmt.Println("The reply is:" + reply)
}
