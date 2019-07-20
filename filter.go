package hl

import (
	"github.com/ravlio/highloadcup2018/account"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/idx"
	"github.com/ravlio/highloadcup2018/querier"
	"github.com/ravlio/highloadcup2018/requests"
	"github.com/ravlio/highloadcup2018/slist"
)

func (d *DB) FilterAccounts(f *requests.Filter) (account.Accounts, error) {
	var accs account.Accounts
	var limit uint8
	if f.Limit.IsSet {
		limit = f.Limit.Val
	} else {
		limit = 50
	}

	// Набор работ. Работа у нас всегда одна — это вызвать f(in) out
	jset := querier.NewJobset(f.FilledCount())

	d.filterToJobset(f, jset)

	if jset.Count() > 0 {
		ids, err := d.querier.Exec(jset, false)

		if err != nil {
			return nil, errors.Wrap(err, "querier exec error")
		}
		if ids != nil {
			defer ids.ReleaseToPool()
		}

		// Сортировки то в фльтре нет, а я сделал зачем-то.
		/*if f.Order.IsSet && f.Order.Asc {
			// left limit
			if limit < len(ids) {
				return d.accounts.GetAccounts(ids[:limit])
			}

			// right limit
			if limit < len(ids) {
				ids=ids[len(ids)-limit:]
			}

			for left, right := 0, len(ids)-1; left < right; left, right = left+1, right-1 {
				ids[left], ids[right] = ids[right], ids[left]
			}

			return d.accounts.GetAccounts(ids)
		}*/

		// TODO предварительно делать limit и order слайсу. Поскольку id в итоговом битмэпе идут в desc порядке, можно заранее их вырезать и
		// отсортировать, только после этого уже лезть за аккаунтами
		ids = account.ReverseIDs(account.RightLimitIDs(ids, int(limit)))
		accs, err = d.accounts.GetAccounts(ids)
		if err != nil {
			return nil, errors.Wrap(err, "accounts error")
		}
		if accs == nil {
			return nil, nil
		}
	} else {
		// если нет ни одной джобы, значит берутся все аккаунты
		// лочим, так как будем работать с оригинальным датасетом
		// TODO передавать лимит и ордер
		accs = d.accounts.GetSliceAndLock()
		defer d.accounts.RUnlock()

		/*		if f.Order.IsSet && f.Order.Asc {
				accs=account.LeftLimitCopy(accs, int(limit))
			}*/

		return account.Reverse(account.RightLimitCopy(accs, int(limit))), nil
	}
	return accs, nil
}

func (d *DB) filterToJobset(f *requests.Filter, jset *querier.Jobset) {
	if f.SexEq.IsSet {
		jset.Add(&querier.Job{
			Name: "sex_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.sex.SexGet(dicts.Sex(f.SexEq.ID))
			}})
	}

	if f.EmailEq.IsSet {
		jset.Add(&querier.Job{
			Name: "email_eq",
			Hash: func() (uint32, bool) {
				return d.hash.email.Uint32Get(f.EmailEq.ID)
			}})
	}

	if f.EmailDomain.IsSet {
		jset.Add(&querier.Job{
			Name: "email_domain",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.emailDomain.Uint32Get(f.EmailDomain.ID)
			}})
	}

	if f.EmailGt.IsSet {
		jset.Add(&querier.Job{
			Name: "email_gt",
			Bitmap: func() *idx.Bitmap {
				r, _, _ := d.skiplist.emailGt.StringSelect(f.EmailGt.Val)
				return r
			}})
	}

	if f.EmailLt.IsSet {
		jset.Add(&querier.Job{
			Name: "email_lt",
			Bitmap: func() *idx.Bitmap {
				r, _, _ := d.skiplist.emailLt.StringSelect(f.EmailLt.Val)
				return r
			}})
	}

	if f.StatusEq.IsSet {
		jset.Add(&querier.Job{
			Name: "status_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.status.StatusGet(dicts.Status(f.StatusEq.ID))
			}})
	}

	if f.StatusNeq.IsSet {
		jset.Add(&querier.Job{
			Name: "status_neq",
			Bitmap: func() *idx.Bitmap {
				var base *idx.Bitmap

				for i := 1; i <= 3; i++ {
					if f.StatusNeq.ID == uint32(i) {
						continue
					}
					base, _ = base.Or(d.bitmap.status.StatusGet(dicts.Status(i)))
				}

				return base
			}})
	}

	if f.FnameEq.IsSet {
		jset.Add(&querier.Job{
			Name: "frame_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.fname.Uint32Get(f.FnameEq.ID)
			}})
	}

	if f.FnameAny.IsSet {
		jset.Add(&querier.Job{
			Name: "fname_any",
			Bitmap: func() *idx.Bitmap {
				var base *idx.Bitmap

				for _, v := range f.FnameAny.ID {
					base, _ = base.Or(d.bitmap.fname.Uint32Get(v))
				}

				return base
			}})
	}

	if f.FnameNull.IsSet {
		jset.Add(&querier.Job{
			Name: "fname_null",
			Bitmap: func() *idx.Bitmap {
				if f.FnameNull.Val {
					return d.bitmap.fnameY
				} else {
					return d.bitmap.fnameN
				}

			}})
	}

	if f.SnameEq.IsSet {
		jset.Add(&querier.Job{
			Name: "sname_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.sname.Uint32Get(f.SnameEq.ID)
			}})
	}

	if f.SnameStarts.IsSet {
		jset.Add(&querier.Job{
			Name: "sname_starts",
			Bitmap: func() *idx.Bitmap {
				return d.trie.sname.Select(f.SnameStarts.Val)
			}})
	}

	if f.SnameNull.IsSet {
		jset.Add(&querier.Job{
			Name: "sname_null",
			Bitmap: func() *idx.Bitmap {
				if f.SnameNull.Val { //sname_null=0, то есть фамилия должна быть
					return d.bitmap.snameY
				} else {
					return d.bitmap.snameN
				}

			}})
	}

	if f.PhoneEq.IsSet {
		jset.Add(&querier.Job{
			Name: "phone_eq",
			Hash: func() (uint32, bool) {
				return d.hash.phone.Int64Get(f.PhoneEq.Int64)
			}})
	}

	if f.PhoneCode.IsSet {
		jset.Add(&querier.Job{
			Name: "phone_code",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.phoneCode.Uint32Get(uint32(f.PhoneCode.Val))
			}})
	}

	if f.PhoneNull.IsSet {
		jset.Add(&querier.Job{
			Name: "phone_null",
			Bitmap: func() *idx.Bitmap {
				if f.PhoneNull.Val {
					return d.bitmap.phoneY
				} else {
					return d.bitmap.phoneN
				}

			}})
	}

	if f.CountryEq.IsSet {
		jset.Add(&querier.Job{
			Name: "country_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.country.Uint32Get(f.CountryEq.ID)
			}})
	}

	if f.CountryNull.IsSet {
		jset.Add(&querier.Job{
			Name: "country_null",
			Bitmap: func() *idx.Bitmap {
				if f.CountryNull.Val {
					return d.bitmap.countryY
				} else {
					return d.bitmap.countryN
				}

			}})
	}

	if f.CityEq.IsSet {
		jset.Add(&querier.Job{
			Name: "city_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.city.Uint32Get(f.CityEq.ID)
			}})
	}

	if f.CityAny.IsSet {
		jset.Add(&querier.Job{
			Name: "city_any",
			Bitmap: func() *idx.Bitmap {
				var base *idx.Bitmap

				for _, v := range f.CityAny.ID {
					base, _ = base.Or(d.bitmap.city.Uint32Get(v))
				}

				return base
			}})
	}

	if f.CityNull.IsSet {
		jset.Add(&querier.Job{
			Name: "city_null",
			Bitmap: func() *idx.Bitmap {
				if f.CityNull.Val {
					return d.bitmap.cityY
				} else {
					return d.bitmap.cityN
				}

			}})
	}

	if f.BirthGt.IsSet {
		jset.Add(&querier.Job{
			Name: "birth_gt",
			Bitmap: func() *idx.Bitmap {
				r, _, _ := d.skiplist.birthGt.Int32Select(int32(f.BirthGt.Val))
				return r
			}})
	}

	if f.BirthLt.IsSet {
		jset.Add(&querier.Job{
			Name: "birth_lt",
			Bitmap: func() *idx.Bitmap {
				r, _, _ := d.skiplist.birthLt.Int32Select(int32(f.BirthLt.Val))
				return r
			}})
	}

	if f.BirthYear.IsSet {
		jset.Add(&querier.Job{
			Name: "birth_year",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.birthYear.Uint32Get(uint32(f.BirthYear.Val))
			}})
	}

	if f.InterestsAny.IsSet {
		jset.Add(&querier.Job{
			Name: "interests_any",
			Bitmap: func() *idx.Bitmap {
				var base *idx.Bitmap

				for _, v := range f.InterestsAny.ID {
					base, _ = base.Or(d.bitmap.interest.Uint32Get(v))
				}

				return base
			}})
	}

	if f.InterestsContains.IsSet {
		jset.Add(&querier.Job{
			Name: "interests_contains",
			Bitmap: func() *idx.Bitmap {
				var base *idx.Bitmap
				var ok bool

				for _, v := range f.InterestsContains.ID {
					base, ok = base.And(d.bitmap.interest.Uint32Get(v))
					if !ok {
						return nil
					}
				}

				return base
			}})
	}

	if f.InterestsEq.IsSet {
		jset.Add(&querier.Job{
			Name: "interests_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.interest.Uint32Get(f.InterestsEq.ID)
			}})
	}

	if f.LikesContains.IsSet {
		jset.Add(&querier.Job{
			Name: "likes_contains",
			Bitmap: func() *idx.Bitmap {
				var l []uint32

				for k, v := range f.LikesContains.Val {
					vl, ok := d.likes[v]
					if !ok {
						return nil
					}

					if k == 0 {
						l = vl
						continue
					}

					l = slist.And(l, vl)
					if l == nil {

						return nil
					}
				}

				return idx.NewBitmapFromSlice(l)
			}})
	}

	if f.PremiumNow.IsSet && f.PremiumNow.Val {
		jset.Add(&querier.Job{
			Name: "premium_now",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.premium
			}})
	}

	if f.PremiumNull.IsSet {
		jset.Add(&querier.Job{
			Name: "premium_null",
			Bitmap: func() *idx.Bitmap {
				if f.PremiumNull.Val {
					return d.bitmap.premiumY
				} else {
					return d.bitmap.premiumN
				}
			}})
	}
}
