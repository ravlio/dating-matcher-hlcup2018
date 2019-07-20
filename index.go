package hl

import "sync"
import "fmt"
import "github.com/rs/zerolog/log"
import "github.com/ravlio/highloadcup2018/requests"
import "github.com/ravlio/highloadcup2018/requests/likes"
import "github.com/ravlio/highloadcup2018/indexer"
import "github.com/ravlio/highloadcup2018/account"
import "github.com/ravlio/highloadcup2018/dicts"
import "github.com/ravlio/highloadcup2018/metrics"
import "github.com/ravlio/highloadcup2018/slist"
import "time"

// import "time"

func (d *DB) PrintIndexStats() {
	//d.bitmap.phoneCode.Uint32PrintKeys()
	s := map[string]int{
		"email":       d.hash.email.Uint32Cardinality(),
		"emailDomain": d.bitmap.emailDomain.Uint32Cardinality(),
		"fname":       d.bitmap.fname.Uint32Cardinality(),
		"sname":       d.bitmap.sname.Uint32Cardinality(),
		"phone":       d.hash.phone.Int64Cardinality(),
		"phoneCode":   d.bitmap.phoneCode.Uint32Cardinality(),
		"country":     d.bitmap.country.Uint32Cardinality(),
		"city":        d.bitmap.city.Uint32Cardinality(),
		"birthYear":   d.bitmap.birthYear.Uint32Cardinality(),
		"joinedYear":  d.bitmap.joinedYear.Uint32Cardinality(),
		"interest":    d.bitmap.interest.Uint32Cardinality(),
		"like":        d.bitmap.like.Uint32Cardinality(),
	}

	for k, v := range s {
		fmt.Printf("|%20s|%6d|\n", k, v)
	}
}

func (d *DB) runIndexWorker(wg *sync.WaitGroup, dd *metrics.Duration) {
	wg.Done()
	var err error
	for {
		select {
		case acc := <-d.accChan:
			t := time.Now()
			err = d.makeMainIndexes(acc.msg, acc.id, acc.opType)
			if err != nil {
				log.Error().Err(err).Msg("Error indexing")
			}
			acc.msg.ReleaseToPool()
			dd.Write(t)
		case l := <-d.likesChan:
			t := time.Now()
			err = d.IndexLikes(l)
			dd.Write(t)
			if err != nil {
				log.Error().Err(err).Msg("Error indexing")
			}
		}
	}
}

func (d *DB) makeMainIndexes(msg *requests.AccountRequest, accID uint32, op opType) error {
	var newAcc *account.Account
	var err error

	a, err := d.accounts.GetAccount(accID)
	if err != nil {
		return err
	}

	// Делаем "иммутабильную" копию текущей версии аккаунта
	// Так как после этого MergeWithRequestAndReturn изменит аккаунт
	curAcc := *a

	// Сохраняем сделанные изменения и возвращаем смердженный аккаунт
	// который является текущей версией
	newAcc, err = d.accounts.MergeWithRequestAndReturn(accID, msg)

	if err != nil {
		return err
	}
	// Email

	jobset := indexer.NewJobset()
	/*if msg.Email.IsSet && (op == AddAccount || curAcc.Email != newAcc.Email)*/ {
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.Email,
			NewUint32: newAcc.Email,
			Add: func() error {
				d.skiplist.emailGt.StringInsert(msg.Email.Val, accID)
				d.skiplist.emailLt.StringInsert(msg.Email.Val, accID)

				return nil
			},
			DeleteUint32: func(id uint32) error {
				d.skiplist.emailGt.StringDelete(dicts.Email.GetKey(id), accID)

				d.skiplist.emailLt.StringDelete(dicts.Email.GetKey(id), accID)
				return nil
			},
		})
	}

	/*if msg.Email.IsSet && (op == AddAccount || curAcc.EmailDomain != newAcc.EmailDomain)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.EmailDomain,
			NewUint32: newAcc.EmailDomain,
			AddUint32: func(domain uint32) error {
				d.bitmap.emailDomain.Uint32GetOrCreate(domain).Add(accID)
				return nil
			},
			DeleteUint32: func(domain uint32) error {
				d.bitmap.emailDomain.Uint32GetOrCreate(domain).Remove(accID)
				return nil
			},
		})
	}

	/*if msg.Fname.IsSet && (op == AddAccount || curAcc.Fname != newAcc.Fname)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.Fname,
			NewUint32: newAcc.Fname,
			AddUint32: func(fname uint32) error {
				d.bitmap.fname.Uint32GetOrCreate(fname).Add(accID)
				return nil
			},
			DeleteUint32: func(fname uint32) error {
				d.bitmap.fname.Uint32GetOrCreate(fname).Remove(accID)
				return nil
			},
			AddYes: func() error {
				d.bitmap.fnameY.Add(accID)
				return nil
			},
			DeleteYes: func() error {
				d.bitmap.fnameY.Remove(accID)
				return nil
			},
			AddNo: func() error {
				d.bitmap.fnameN.Add(accID)
				return nil
			},
			DeleteNo: func() error {
				d.bitmap.fnameN.Remove(accID)
				return nil

			},
		})
	}

	/*if msg.Sname.IsSet && (op == AddAccount || curAcc.Sname != newAcc.Sname)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.Sname,
			NewUint32: newAcc.Sname,
			AddUint32: func(sname uint32) error {
				d.bitmap.sname.Uint32GetOrCreate(sname).Add(accID)
				d.trie.sname.Insert(msg.Sname.Val, accID)

				return nil
			},
			DeleteUint32: func(sname uint32) error {
				d.bitmap.sname.Uint32GetOrCreate(sname).Remove(accID)
				d.trie.sname.Delete(dicts.Sname.GetKey(sname), accID)
				return nil
			},
			AddYes: func() error {
				d.bitmap.snameY.Add(accID)
				return nil
			},
			DeleteYes: func() error {
				d.bitmap.snameY.Remove(accID)
				return nil
			},
			AddNo: func() error {
				d.bitmap.snameN.Add(accID)
				return nil

			},
			DeleteNo: func() error {
				d.bitmap.snameN.Remove(accID)
				return nil
			},
		})
	}

	/*if msg.Phone.IsSet && (op == AddAccount || curAcc.Phone != newAcc.Phone)*/
	{
		jobset.Add(&indexer.Job{
			VarType:  indexer.VarInt64,
			CurInt64: curAcc.Phone,
			NewInt64: newAcc.Phone,
			/*AddInt64: func(phone int64) error {
				d.hash.phone.Int64Set(phone, accID)
				return nil
			},
			DeleteInt64: func(phone int64) error {
				d.hash.phone.Int64Delete(phone)
				return nil
			},*/
			AddYes: func() error {
				d.bitmap.phoneY.Add(accID)
				return nil
			},
			DeleteYes: func() error {
				d.bitmap.phoneY.Remove(accID)
				return nil
			},
			AddNo: func() error {
				d.bitmap.phoneN.Add(accID)
				return nil

			},
			DeleteNo: func() error {
				d.bitmap.phoneN.Remove(accID)
				return nil
			},
		})
	}

	cph := curAcc.GetPhoneCode()
	nph := newAcc.GetPhoneCode()
	/*if msg.Phone.IsSet && (op == AddAccount || curAcc.GetPhoneCode() != newAcc.GetPhoneCode())*/ {
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: uint32(cph),
			NewUint32: uint32(nph),
			AddUint32: func(code uint32) error {
				d.bitmap.phoneCode.Uint32GetOrCreate(code).Add(accID)
				return nil
			},
			DeleteUint32: func(code uint32) error {
				d.bitmap.phoneCode.Uint32GetOrCreate(code).Remove(accID)
				return nil
			},
		})
	}

	/*if msg.Sex.IsSet && (op == AddAccount || curAcc.Sex != newAcc.Sex)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: uint32(curAcc.Sex),
			NewUint32: uint32(newAcc.Sex),
			AddUint32: func(sex uint32) error {
				d.bitmap.sex.SexGet(dicts.Sex(sex)).Add(accID)
				return nil
			},
			DeleteUint32: func(sex uint32) error {
				d.bitmap.sex.SexGet(dicts.Sex(sex)).Remove(accID)
				return nil
			},
		})
	}

	/*if msg.Birth.IsSet */
	{
		/*if op == AddAccount || curAcc.Birth != newAcc.Birth */ {
			jobset.Add(&indexer.Job{
				VarType:  indexer.VarInt32,
				CurInt32: curAcc.Birth,
				NewInt32: newAcc.Birth,
				AddInt32: func(birth int32) error {
					d.skiplist.birthLt.Int32Insert(birth, accID)
					d.skiplist.birthGt.Int32Insert(birth, accID)
					return nil
				},
				DeleteInt32: func(birth int32) error {
					d.skiplist.birthLt.Int32Delete(birth, accID)
					d.skiplist.birthGt.Int32Delete(birth, accID)
					return nil
				},
			})
		}

		curBirthYear := curAcc.GetBirthYear()
		newBirthYear := newAcc.GetBirthYear()
		/*if op == AddAccount || curBirthYear != newBirthYear */ {

			jobset.Add(&indexer.Job{
				VarType:   indexer.VarUint32,
				CurUint32: curBirthYear,
				NewUint32: newBirthYear,
				AddUint32: func(year uint32) error {
					d.bitmap.birthYear.Uint32GetOrCreate(year).Add(accID)
					return nil
				},
				DeleteUint32: func(year uint32) error {
					d.bitmap.birthYear.Uint32GetOrCreate(year).Remove(accID)
					return nil
				},
			})
		}
	}

	/*if msg.City.IsSet && (op == AddAccount || curAcc.City != newAcc.City)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.City,
			NewUint32: newAcc.City,
			AddUint32: func(city uint32) error {
				d.bitmap.city.Uint32GetOrCreate(city).Add(accID)
				return nil
			},
			DeleteUint32: func(city uint32) error {
				d.bitmap.city.Uint32GetOrCreate(city).Remove(accID)
				return nil
			},
			AddYes: func() error {
				d.bitmap.cityY.Add(accID)
				return nil
			},
			DeleteYes: func() error {
				d.bitmap.cityY.Remove(accID)
				return nil
			},
			AddNo: func() error {
				d.bitmap.cityN.Add(accID)
				return nil
			},
			DeleteNo: func() error {
				d.bitmap.cityN.Remove(accID)
				return nil
			},
		})
	}

	/*if msg.Country.IsSet && (op == AddAccount || curAcc.Country != newAcc.Country)*/
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curAcc.Country,
			NewUint32: newAcc.Country,
			AddUint32: func(country uint32) error {
				d.bitmap.country.Uint32GetOrCreate(country).Add(accID)
				return nil
			},
			DeleteUint32: func(country uint32) error {
				d.bitmap.country.Uint32GetOrCreate(country).Remove(accID)
				return nil
			},
			AddYes: func() error {
				d.bitmap.countryY.Add(accID)
				return nil
			},
			DeleteYes: func() error {
				d.bitmap.countryY.Remove(accID)
				return nil
			},
			AddNo: func() error {
				d.bitmap.countryN.Add(accID)
				return nil
			},
			DeleteNo: func() error {
				d.bitmap.countryN.Remove(accID)
				return nil
			},
		})
	}

	curJoinedYear := curAcc.GetJoinedYear()
	newJoinedYear := newAcc.GetJoinedYear()

	/*if op == AddAccount || curJoinedYear != curJoinedYear */
	{
		jobset.Add(&indexer.Job{
			VarType:   indexer.VarUint32,
			CurUint32: curJoinedYear,
			NewUint32: newJoinedYear,
			AddUint32: func(year uint32) error {
				d.bitmap.joinedYear.Uint32GetOrCreate(year).Add(accID)
				return nil
			},
			DeleteUint32: func(year uint32) error {
				d.bitmap.joinedYear.Uint32GetOrCreate(year).Remove(accID)
				return nil
			},
		})
	}

	ct := d.curTime.Unix()

	jobset.Add(&indexer.Job{
		VarType: indexer.VarBool,
		CurCond: curAcc.PremiumStart > 0 &&
			curAcc.PremiumFinish > 0 &&
			ct >= int64(curAcc.PremiumStart) &&
			ct <= int64(curAcc.PremiumFinish),
		NewCond: newAcc.PremiumStart > 0 &&
			newAcc.PremiumFinish > 0 &&
			ct >= int64(newAcc.PremiumStart) &&
			ct <= int64(newAcc.PremiumFinish),
		Add: func() error {
			d.bitmap.premium.Add(accID)
			return nil
		},
		Delete: func() error {
			d.bitmap.premium.Remove(accID)
			return nil
		},
	})

	jobset.Add(&indexer.Job{
		VarType: indexer.VarBool,
		CurCond: curAcc.PremiumStart > 0 && curAcc.PremiumFinish > 0,
		NewCond: newAcc.PremiumStart > 0 && newAcc.PremiumFinish > 0,
		AddYes: func() error {
			d.bitmap.premiumY.Add(accID)
			return nil
		},
		DeleteYes: func() error {
			d.bitmap.premiumY.Remove(accID)
			return nil
		},
		AddNo: func() error {
			d.bitmap.premiumN.Add(accID)
			return nil
		},
		DeleteNo: func() error {
			d.bitmap.premiumN.Remove(accID)
			return nil
		},
	})

	/*if msg.Status.IsSet && (op == AddAccount || curAcc.Status != newAcc.Status) */
	{
		jobset.Add(&indexer.Job{
			CurUint32: uint32(curAcc.Status),
			NewUint32: uint32(newAcc.Status),
			AddUint32: func(status uint32) error {
				d.bitmap.status.StatusGet(dicts.Status(status)).Add(accID)
				// d.bitmap.statusNeq.StatusGet(dicts.Status(status)).Remove(accID)
				return nil
			},
			DeleteUint32: func(status uint32) error {
				d.bitmap.status.StatusGet(dicts.Status(status)).Remove(accID)
				// d.bitmap.statusNeq.StatusGet(dicts.Status(status)).Add(accID)
				return nil
			},
		})
	}

	jobset.Add(&indexer.Job{
		VarType:        indexer.VarUint32Slice,
		CurUint32Slice: curAcc.Interests,
		NewUint32Slice: newAcc.Interests,
		AddUint32: func(id uint32) error {
			d.bitmap.interest.Uint32GetOrCreate(id).Add(accID)
			return nil
		},
		DeleteUint32: func(id uint32) error {
			d.bitmap.interest.Uint32GetOrCreate(id).Remove(accID)
			return nil
		},
	})

	if len(msg.Likes) > 0 {
		var l = make(likes.Likes, 0, len(msg.Likes))

		for _, v := range msg.Likes {
			l = append(l, likes.Like{
				Liker: accID,
				Likee: v.ID,
				Ts:    v.Ts,
			})
		}

		// Минуя канал индексации лайков, так как мы итак уже из каналв
		d.IndexLikes(l)
	}

	if op == AddAccount {
		return d.indexer.Insert(jobset)
	} else {
		return d.indexer.Update(jobset)
	}
}

func (d *DB) ValidateLikes(ls likes.Likes) error {
	ex := map[uint32]struct{}{}

	accs := make([]uint32, 0, len(ls)*2)
	for _, l := range ls {
		if _, ok := ex[l.Likee]; !ok {
			accs = append(accs, l.Likee)
		}
		if _, ok := ex[l.Liker]; !ok {
			accs = append(accs, l.Liker)
		}
	}

	if err := d.accounts.CheckExistence(accs); err != nil {
		return err
	}

	return nil
}

func (db *DB) AddLikes(ls likes.Likes) {
	db.likesChan <- ls
}

func (d *DB) IndexLikes(ls likes.Likes) error {
	d.lmx.Lock()
	for _, l := range ls {
		dl, ok := d.likes[l.Likee]
		if !ok {
			dl := make([]uint32, 0, 34)
			slist.Insert(&dl, l.Liker)
			d.likes[l.Likee] = dl

		} else {
			// m=append(m,l.Liker)
			slist.Insert(&dl, l.Liker)
			d.likes[l.Likee] = dl
		}
	}

	d.lmx.Unlock()
	return nil
}
