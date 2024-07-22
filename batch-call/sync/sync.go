package sync

import "fmt"

type SyncTarget interface {
	GetToken() string
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

func (at *AccountTarget) GetToken() string {
	return "Account"
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

func (vt *ValueTarget) GetToken() string {
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
