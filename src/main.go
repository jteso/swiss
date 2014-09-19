package main

import (
	"circuitbreaker"
	"errors"
	"time"
)

type DodgyCommand struct {}
func (_ DodgyCommand) Run()(interface{}, error) {
	return 0, errors.New("oh oh,... service is returning an error")
}

func healthyService() (interface {}, error){
	return 0, nil
}




func main() {
	cb := circuitbreaker.New("test", 3, time.Duration(10) * time.Second)
	numFailures := 10
	for {
		if numFailures > 0 {
			time.Sleep(3 * time.Second)
			cb.Execute(DodgyCommand{})
			numFailures--
		}else{
			time.Sleep(3 * time.Second)
			cb.Execute(circuitbreaker.CommandFunc(healthyService))
		}
	}


}
