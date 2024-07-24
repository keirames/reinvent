package sync

import "fmt"

type SyncTarget interface {
	CanHandle(t string) bool
	Execute() error
}

type AccountTarget struct {
}

func NewAccountTarget() *AccountTarget {
	return new(AccountTarget)
}

func (at *AccountTarget) Execute() error {
	fmt.Println("account sync execute")
	return nil
}

func (at *AccountTarget) CanHandle(t string) bool {
	return t == "Account"
}

type ValueTarget struct {
}

func NewValueTarget() *ValueTarget {
	return new(ValueTarget)
}

func (vt *ValueTarget) Execute() error {
	fmt.Println("value sync execute")
	return nil
}

func (vt *ValueTarget) CanHandle(t string) bool {
	return t == "Value"
}

func GetExecutorByToken(token string) (SyncTarget, error) {
	arr := []SyncTarget{}

	arr = append(arr, NewAccountTarget())
	arr = append(arr, NewValueTarget())

	var handler SyncTarget
	for _, t := range arr {
		if t.CanHandle(token) {
			handler = t
			break
		}
	}
	if handler == nil {
		fmt.Println("no handler found!")
	}

	return handler, nil
}
