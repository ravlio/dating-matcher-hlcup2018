package main

type Account struct{}

func (a *Account) MergeWithRequest(msg *AccountRequest) (*Account, error) {
	return a, nil
}

type Accounts struct {
	accs map[uint32]*Account
}

func (a *Accounts) MergeWithMessageAndReturn(id uint32, msg *AccountRequest) (*Account, error) {
	r, err := a.accs[id].MergeWithRequest(msg)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (a *Accounts) GetCopyOfAccount(id uint32) (*Account, error) {
	acc, _ := a.accs[id]

	r := *acc
	return &r, nil
}

func (a *Accounts) GetAccount(id uint32) (*Account, error) {
	acc, _ := a.accs[id]

	return acc, nil
}

func main() {
	accs := &Accounts{accs: map[uint32]*Account{1: {}}}
	acc, err := accs.MergeWithMessageAndReturn(1, &AccountRequest{})

	acc, err = accs.GetCopyOfAccount(1)

	_ = acc
	_ = err

	acc, err = accs.GetAccount(1)

	cp := *acc
	_ = cp

}
