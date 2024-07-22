package sync

import "fmt"

type SyncTarget interface {
	getToken() string
	execute()
}

type AccountTarget struct {
}

func NewAccountTarget() *AccountTarget {
	return new(AccountTarget)
}

func (at *AccountTarget) execute() {
	fmt.Println("account sync execute")
}

func (at *AccountTarget) getToken() string {
	return "Account"
}

type ValueTarget struct {
}

func NewValueTarget() *ValueTarget {
	return new(ValueTarget)
}

func (vt *ValueTarget) execute() {
	fmt.Println("value sync execute")
}

func (vt *ValueTarget) getToken() string {
	return "Value"
}

func GetExecutorByToken(token string) (SyncTarget, error) {
	m := make(map[string]SyncTarget)

	m["Account"] = NewAccountTarget()
	m["Value"] = NewValueTarget()

	handler, ok := m[token]
	if !ok {
		return nil, fmt.Errorf("unknown token %v", token)
	}

	return handler, nil
}
