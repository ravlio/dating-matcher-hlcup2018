package hl

import "github.com/ravlio/highloadcup2018/requests"
import "github.com/ravlio/highloadcup2018/account"
import "github.com/ravlio/highloadcup2018/dicts"
import "github.com/ravlio/highloadcup2018/errors"

// import "github.com/valyala/fastjson"

/*
Добавление или обновление аккаутна. Универсальный метод, позволяющий обойтись дублированием кода
defer не используются в целях экономии
*/
func (d *DB) CreateAccount(accReq *requests.AccountRequest) error {
	// d.opChan <- opMsg{opType: AddAccount, msg: accReq, id: 1}
	// return nil
	var emailID uint32
	var accID = accReq.ID.Val

	// TODO сразу добавлять ID емейла в accReq, чтобы FillWithRequest повторно не лез в словари
	if accReq.Email.IsSet {
		emailID = dicts.Email.GetOrCreateValue(accReq.Email.Val)
	}
	if accReq.Email.IsSet {
		// Проверка уникальности емейла
		if _, f := d.hash.email.Uint32Get(emailID); f {
			return account.ErrEmailAlreadyExists
		}
	}

	if accReq.Phone.IsSet {
		// Проверка уникальности емейла
		if _, f := d.hash.phone.Int64Get(accReq.Phone.Int64); f {
			return account.ErrPhoneAlreadyExists
		}
	}

	var err error

	// TODO Сделать проверку и перевести акки со слайса на скиплит, так как встречаются разрывы
	/*if accReq.ID.IsSet && d.accounts.LastID()!=accReq.ID.Val-1 {
		return errors.New("new id is not in sequence")
	}*/

	acc := &account.Account{ID: accID}
	err = d.accounts.AddAccount(acc)
	if err != nil {
		return err
	}

	if accReq.Email.IsSet {
		d.hash.email.Uint32Set(emailID, accID)
	}

	if accReq.Phone.IsSet {
		d.hash.phone.Int64Set(accReq.Phone.Int64, accID)
	}

	d.accChan <- opMsg{opType: AddAccount, msg: accReq, id: accID}

	return nil
}

func (d *DB) UpdateAccount(accReq *requests.AccountRequest, accID uint32) error {
	var emailID uint32

	// Проверяем, есть ли акк
	if _, err := d.accounts.GetAccount(accID); err != nil {
		return err
	}

	if accReq.Email.IsSet {
		emailID = dicts.Email.GetOrCreateValue(accReq.Email.Val)
	}

	// Из-за возможных разрывов такой способ не подходит, нужно брать акк явно
	/*if d.accounts.LastID()<accID {
		return ErrUnexistingAccount
	}*/
	// Обновляем емейл, если он передан
	if accReq.Email.IsSet {
		// Емейл есть и он занят за кем-то другим
		eID, f := d.hash.email.Uint32Get(emailID)
		if f && eID != accID {
			return errors.New("accout have different email")
		} else if f && eID == accID {
			d.hash.email.Uint32Change(emailID, accID)
		} else if !f {
			d.hash.email.Uint32Set(emailID, accID)
		}
	}

	if accReq.Phone.IsSet {
		// Емейл есть и он занят за кем-то другим
		eID, f := d.hash.phone.Int64Get(accReq.Phone.Int64)
		if f && eID != accID {
			return errors.New("accout have different email")
		} else if f && eID == accID {
			d.hash.phone.Int64Change(accReq.Phone.Int64, accID)
		} else if !f {
			d.hash.phone.Int64Set(accReq.Phone.Int64, accID)
		}
	}

	d.accChan <- opMsg{opType: UpdateAccount, msg: accReq, id: accID}

	return nil
}
