package customerrors

import "fmt"

type URLNotFound struct {
	info string
}

func (err URLNotFound) Error() string {
	return fmt.Sprintf("url not found: %v", err.info)
}

type WrongURL struct {
	info string
}

func (err WrongURL) Error() string {
	return fmt.Sprintf("wrong url: %v", err.info)
}