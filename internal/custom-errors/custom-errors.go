package customerrors

import "fmt"

type URLNotFound struct {
	Info string
}

func (err URLNotFound) Error() string {
	return fmt.Sprintf("url not found: %v", err.Info)
}

type WrongURL struct {
	Info string
}

func (err WrongURL) Error() string {
	return fmt.Sprintf("wrong url: %v", err.Info)
}