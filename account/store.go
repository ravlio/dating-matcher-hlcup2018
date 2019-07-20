package account

import (
	"errors"
	"github.com/ravlio/highloadcup2018/requests"
	"sync"
	// "sync/atomic"
)

var ErrUnexistingAccount = errors.New("account does not exists")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrPhoneAlreadyExists = errors.New("phone already exists")

type Store struct {
	//accs      map[uint32]*Account
	// слайс по индексу намного быстрее мэпы: 1ns vs 121ns
	accs []*Account
	// Если запрос придёт без фильтров, нужно всегда брать полный список
	// этот слайс отличается отсутствием дырок
	accsSlice []*Account
	lastID    uint32
	mx        sync.RWMutex
}

func NewAccounts() *Store {
	return &Store{accs: /*make(map[uint32]*Account)*/ make([]*Account, 1300000), accsSlice: make([]*Account, 0, 1300000)}
}

/*func (a *Store) LastID() uint32 {
	return atomic.LoadUint32(&a.lastID)
}*/

func (a *Store) LockAccountAndSet(acc *Account) error {
	acc.Mx.Lock()
	a.mx.Lock()
	a.accs[acc.ID] = acc
	a.mx.Unlock()
	return nil
}

func (a *Store) AddAccount(acc *Account) error {
	accID := acc.ID
	a.mx.Lock()
	if len(a.accs) <= int(accID) {
		a.accs = append(a.accs, make([]*Account, int(accID)+1-len(a.accs))...)
	}
	a.accs[acc.ID] = acc

	if cap(a.accsSlice) <= int(accID) {
		a.accsSlice = append(a.accsSlice, make([]*Account, 0, 5000)...)
	}
	a.accsSlice = append(a.accsSlice, acc)
	a.mx.Unlock()

	return nil
}

func (a *Store) SetAccount(acc *Account) error {
	a.mx.Lock()
	a.accs[acc.ID] = acc
	a.mx.Unlock()
	return nil
}

func (a *Store) MergeWithRequestAndReturn(id uint32, req *requests.AccountRequest) (*Account, error) {
	a.mx.Lock()
	acc := a.accs[id]
	err := acc.MergeWithRequest(req)
	a.mx.Unlock()
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (a *Store) GetCopyOfAccount(id uint32) (*Account, error) {
	if int(id) > len(a.accs) {
		return nil, ErrUnexistingAccount
	}
	a.mx.RLock()
	acc := a.accs[id]
	a.mx.RUnlock()

	if acc == nil {
		return nil, ErrUnexistingAccount
	}

	r := *acc
	return &r, nil
}

func (a *Store) GetAccount(id uint32) (*Account, error) {
	if int(id) > len(a.accs) {
		return nil, ErrUnexistingAccount
	}
	a.mx.RLock()
	acc := a.accs[id]
	a.mx.RUnlock()

	if acc == nil {
		return nil, ErrUnexistingAccount
	}

	return acc, nil
}

func (a *Store) GetAccountUnsafe(id uint32) (*Account, error) {
	if int(id) > len(a.accs) {
		return nil, ErrUnexistingAccount
	}
	acc := a.accs[id]

	if acc == nil {
		return nil, ErrUnexistingAccount
	}

	return acc, nil
}

func (a *Store) GetAccounts(ids []uint32) (Accounts, error) {
	a.mx.RLock()
	ret := make([]*Account, len(ids))
	for k, id := range ids {
		// вот так по простому, без проверок
		ret[k] = a.accs[id]
	}

	a.mx.RUnlock()

	return ret, nil

}

func (a *Store) CheckExistence(ids []uint32) error {
	a.mx.RLock()
	for _, id := range ids {
		if int(id) > len(a.accs) {
			a.mx.RUnlock()
			return ErrUnexistingAccount
		}

		ok := a.accs[id]
		if ok == nil {
			a.mx.RUnlock()
			return ErrUnexistingAccount
		}
	}

	a.mx.RUnlock()

	return nil

}

func (a *Store) GetSliceAndLock() Accounts {
	a.mx.RLock()
	return a.accsSlice
}

func (a *Store) RUnlock() {
	a.mx.RUnlock()
}

func (a *Store) SortSlice() {
	Sort(a.accsSlice, SortAsc)
}

func (a *Store) GetAccountAndLock(id uint32) (*Account, error) {
	if int(id) > len(a.accs) {
		return nil, ErrUnexistingAccount
	}

	a.mx.RLock()
	acc := a.accs[id]
	a.mx.RUnlock()

	if acc == nil {
		return nil, ErrUnexistingAccount
	}

	acc.Mx.Lock()
	return acc, nil
}

func (a *Store) UpdateAccount(id uint32, account requests.AccountRequest) error {
	return nil
}
