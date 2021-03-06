package circuitbreaker

import (
	"time"
	"fmt"
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
)

//Some generic errors

func ErrCircuitBreakerKeyDoesNotExist(key string) error{
	return fmt.Errorf("Cannot be found a circuit Breaker with key: %s", key)
}

var ErrCircuitBreakerOpened = fmt.Errorf("Circuit breaker is OPEN")



// >>>>>>>>>>>>>>>>>>>
// Public Constructors
// >>>>>>>>>>>>>>>>>>>
func New(key string, closeAfter int, halfOpenAfter time.Duration) (ICircuitBreaker) { // todo - read this to understand why not to put a pointer to interface: http://stackoverflow.com/questions/13511203/why-cant-i-assign-a-struct-to-an-interface
	return theCircuitBreakerFactory.getOrCreateCircuitBreakerWithKey(key, closeAfter, halfOpenAfter)
}
func Get(key string) (ICircuitBreaker, error) {
	i, ok := theCircuitBreakerFactory.cbs[key]
	if ok{
		return i, nil
	}else{
		return nil, ErrCircuitBreakerKeyDoesNotExist(key)
	}
}

// This is a singleton - look somewhere else if you dont like it, you hipster
var	theCircuitBreakerFactory  = &circuitBreakerFactory{make(map[string] *CircuitBreaker)}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>
// Factory of Circuit Breakers
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>
type circuitBreakerFactory struct{
	cbs map[string] *CircuitBreaker
}
func (cbf *circuitBreakerFactory)getOrCreateCircuitBreakerWithKey(key string, closeAfter int, halfOpenAfter time.Duration) (ICircuitBreaker){
	//	if cbf.cbs == nil {
	//		cbf.cbs = make(map[string] *CircuitBreaker)
	//	}
	i, ok := cbf.cbs[key]
	if ok{
		return i
	}else{
		newCircuitBreaker := &CircuitBreaker{
			name:key,
			CircuitBreakerState: CLOSED,
			closeAfter: closeAfter,
			halfOpenAfter: halfOpenAfter,
			logger: stdlog.GetFromFlags(),
			failureCount: closeAfter}

		cbf.cbs[key] = newCircuitBreaker
		return newCircuitBreaker
	}
}

// >>>>>>>>>>>>>>>>>>>
// Circuit Breaker
// >>>>>>>>>>>>>>>>>>>

type ICircuitBreaker interface {
	Execute(command Command)(interface{}, error)
}

//states of the circuit breaker
type CircuitBreakerState uint8

const (
	OPEN CircuitBreakerState = iota
	HALF_OPEN
	CLOSED
)

type CircuitBreaker struct {
	name string
	CircuitBreakerState
	closeAfter int
	halfOpenAfter time.Duration
	logger log.Logger

	failureCount int // <-- if 0 trip the breaker
	// todo - analytics
	// todo - monitoring
}

func (c *CircuitBreaker) transitionTo(cbs CircuitBreakerState) {
	switch cbs{
	case OPEN:
		c.CircuitBreakerState = OPEN
	case HALF_OPEN:
		c.CircuitBreakerState = HALF_OPEN
	case CLOSED:
		c.CircuitBreakerState = CLOSED
	}
}
func (c *CircuitBreaker) resetFailureCount() {
	c.failureCount = c.closeAfter
}
func (c *CircuitBreaker) setResetTimer(){
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(c.halfOpenAfter)
		timeout <- true
	}()
	<-timeout
	c.attemptReset()
}
func (c *CircuitBreaker) tripBreaker() {
	c.transitionTo(OPEN)
	go c.setResetTimer()
	c.logger.Info("Circuit breaker has been tripped")

}
func (c *CircuitBreaker) attemptReset() {
	c.logger.Info("Attempting to reset the circuit")
	c.transitionTo(HALF_OPEN)
}

func (c *CircuitBreaker) resetBreaker() {
	c.logger.Info("Circuit breaker has been reset")
	c.transitionTo(CLOSED)
	c.resetFailureCount()
}

func (c CircuitBreaker) AllowRequest() bool{
	return c.CircuitBreakerState == CLOSED || c.CircuitBreakerState == HALF_OPEN
}

func (c *CircuitBreaker) MarkSuccess() {
	c.logger.Info("Service invocation has succeeded")
	switch c.CircuitBreakerState{
	case CLOSED:
		c.resetFailureCount()
	case HALF_OPEN:
		c.resetBreaker() //close it and reset the failure count
	}
}

func (c *CircuitBreaker) MarkFailure(){
	c.logger.Info("Service invocation has failed")
	switch c.CircuitBreakerState{
	case CLOSED:
		c.failureCount -=1
		if c.failureCount == 0 {
			c.tripBreaker() // state to open and set a timer to set it to half-open
		}
	case HALF_OPEN:
		c.tripBreaker()

	}
}

func (c *CircuitBreaker)Execute(command Command)(interface{}, error){
	if c.AllowRequest(){
		result,err := command.Run()
		if err != nil {
			c.MarkFailure()
		}else{
			c.MarkSuccess()
		}
		return result, err
	}else{
		c.logger.Info("Service invocation not allow due circuit breaker is OPEN")
		return -1, ErrCircuitBreakerOpened
	}

}

// Wrapper for your command that needs the protection of a Circuit Breaker
// For example:
// ------------
// 1. Like a struct:
// > type AddCommand struct {}
// > func (_ AddCommand) Run()(interface{}, error){ ... }
// > cb := circuitbreaker.New()
// > cb.Execute(AddCommand)
// 2. Like a func:
// > func AddCommandFunc()(interface{}, error){ ... }
// > cb.Execute(CommandFunc(AddCommand))



type Command interface {
	Run()(interface {}, error)
}
type CommandFunc func()(interface {}, error)
func (c CommandFunc) Run() (interface {}, error){
	return c()
}


