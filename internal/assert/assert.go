package assert

import (
	"cmp"
	"fmt"
	"log"
)

func Equal[T comparable](lhs T, rhs T) {
	if lhs != rhs {
		log.Panicf("%s %s", messageHead("Equal"), binaryOperatorMessageBody("==", lhs, rhs))
	}
}

func InEqual[T comparable](lhs T, rhs T) {
	if lhs == rhs {
		log.Panicf("%s %s", messageHead("InEqual"), binaryOperatorMessageBody("!=", lhs, rhs))
	}
}

func GreaterThan[T cmp.Ordered](lhs T, rhs T) {
	if lhs <= rhs {
		log.Panicf("%s %s", messageHead("GreaterThan"), binaryOperatorMessageBody(">", lhs, rhs))
	}
}

func GreaterThanOrEqual[T cmp.Ordered](lhs T, rhs T) {
	if lhs < rhs {
		log.Panicf("%s %s", messageHead("GreaterThanOrEqual"), binaryOperatorMessageBody(">=", lhs, rhs))
	}
}

func LessThan[T cmp.Ordered](lhs T, rhs T) {
	if lhs >= rhs {
		log.Panicf("%s %s", messageHead("LessThan"), binaryOperatorMessageBody("<", lhs, rhs))
	}
}

func LessThanOrEqual[T cmp.Ordered](lhs T, rhs T) {
	if lhs > rhs {
		log.Panicf("%s %s", messageHead("LessThanOrEqual"), binaryOperatorMessageBody("<=", lhs, rhs))
	}
}

func NilError(err error) {
	if err != nil {
		log.Panicf("%s %s", messageHead("NilError"), fmt.Sprintf("Expected nil error; Received error: %s", err.Error()))
	}
}

func messageHead(funcName string) string {
	return fmt.Sprintf("[assert.%s]", funcName)
}

func binaryOperatorMessageBody(operator string, lhs any, rhs any) string {
	return fmt.Sprintf("Expected lhs %s rhs; Received %v %s %v", operator, lhs, operator, rhs)
}
