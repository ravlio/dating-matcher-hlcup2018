package hl

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/account"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/idx"
	"github.com/ravlio/highloadcup2018/querier"
	"github.com/ravlio/highloadcup2018/requests"
)

func (d *DB) GroupAccounts(g *requests.Group) (account.RawGroups, error) {
	var fullscan bool
	var ids roaring.Arr
	var accs account.Accounts
	var acc *account.Account
	var err error
	var limit int
	if g.Limit.IsSet {
		limit = int(g.Limit.Val)
	} else {
		limit = 50
	}

	// Набор работ. Работа у нас всегда одна — это вызвать f(in) out
	jset := querier.NewJobset(1)

	d.groupToJobset(g, jset)

	if jset.Count() > 0 {
		ids, err = d.querier.Exec(jset, false)
		if err != nil {
			return nil, errors.Wrap(err, "querier exec error")
		}

		if ids == nil {
			return nil, nil
		}

		defer ids.ReleaseToPool()
	} else {
		fullscan = true
		accs = d.accounts.GetSliceAndLock()
		d.accounts.RUnlock()
	}

	groups := make(map[uint64]*account.Group)
	var pk uint64 // sex(1)+status(1)+interests(10000)+country(1000)+city(10000)

	var i int
	for {
		pk = 0

		if fullscan {
			if i >= len(accs) {
				break
			}

			acc = accs[i]

		} else {
			if i >= len(ids) {
				break
			}
			acc, err = d.accounts.GetAccountUnsafe(ids[i])
			if err != nil {
				return nil, errors.Wrap(err, "account error")
			}
		}

		i++

		if g.Keys.Has(requests.GroupSex) {
			if acc.Sex > 0 {
				pk += uint64(acc.Sex) // 1 байт
			} else {
				// continue // игнорим нулевые значения, так как задание не подразумевает нулевых значений размерностей
			}
		}

		if g.Keys.Has(requests.GroupStatus) {
			if acc.Status > 0 {
				pk += uint64(acc.Status) * 10 // 1 байт + 1 байт
			} else {
				// continue
			}
		}

		if g.Keys.Has(requests.GroupCountry) {
			if acc.Country > 0 {
				pk += uint64(acc.Country) * 100 // 7 байт + 3 байта
			} else {
				// continue
			}
		}

		if g.Keys.Has(requests.GroupCity) {
			if acc.City > 0 {
				pk += uint64(acc.City) * 100000 // 10 байт + 4 байта
			} else {
				// continue
			}
		}

		if g.Keys.Has(requests.GroupInterests) {
			if acc.Interests != nil {
				opk := pk // оригинальный ключ
				for _, in := range acc.Interests {
					pk := opk                     // сбрасываем клю до интересов
					pk += uint64(in) * 1000000000 // Добавляем интерес в ключ. 2 байта +5 байт (1-99999)

					d.incrementGroup(g, groups, pk, acc, in)
				}
			} else {
				// continue
			}
		} else {
			d.incrementGroup(g, groups, pk, acc, 0)
		}
	}

	ret := make([]*account.Group, len(groups))

	//  перемещаем группы в слайс
	i = 0
	for _, v := range groups {
		ret[i] = v
		i++
	}

	if !g.Order.IsSet || g.Order.Asc {
		account.GroupSort(ret, account.SortAsc)
	} else {
		account.GroupSort(ret, account.SortDesc)
	}

	// Второй этап сортировки. Находим всё с равными count, тянем строки и сортируем по ним
	var rgroups account.RawGroups

	if limit > len(ret) {
		rgroups = make(account.RawGroups, 0, limit)
	} else {
		rgroups = make(account.RawGroups, 0, len(ret))
	}

	var prevc uint32
	for i, r := range ret {
		// добавляем в группы все с равными количествами в том числе
		if i < limit || r.Count == prevc {
			rgroups = append(rgroups, account.MakeRawGroup(r))
		} else {
			break
		}
		prevc = r.Count
	}

	if !g.Order.IsSet || g.Order.Asc {
		account.RawGroupSort(rgroups, account.SortAsc, g.KeysRaw)
	} else {
		account.RawGroupSort(rgroups, account.SortDesc, g.KeysRaw)
	}

	if limit > len(rgroups) {
		return rgroups, nil
	}

	return rgroups[:limit], nil
}

func (d *DB) incrementGroup(g *requests.Group, groups map[uint64]*account.Group, pk uint64, acc *account.Account, in uint32) {
	cg, ok := groups[pk]
	if !ok {
		cg = &account.Group{
			Count: 1,
		}

		if g.Keys.Has(requests.GroupSex) {
			cg.Sex = acc.Sex
		}
		if g.Keys.Has(requests.GroupStatus) {
			cg.Status = acc.Status
		}
		if g.Keys.Has(requests.GroupCountry) {
			cg.Country = acc.Country
		}
		if g.Keys.Has(requests.GroupCity) {
			cg.City = acc.City
		}

		if g.Keys.Has(requests.GroupInterests) {

			cg.Interests = in
		}

		groups[pk] = cg
	} else {
		cg.Count++
	}
}

func (d *DB) groupToJobset(f *requests.Group, jset *querier.Jobset) {
	if f.Sex.IsSet {
		jset.Add(&querier.Job{
			Name: "sex_eq", // TODO переделать на константы
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.sex.SexGet(dicts.Sex(f.Sex.ID))
			}})
	}

	if f.Status.IsSet {
		jset.Add(&querier.Job{
			Name: "status_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.status.StatusGet(dicts.Status(f.Status.ID))
			}})
	}

	if f.Country.IsSet {
		jset.Add(&querier.Job{
			Name: "country_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.country.Uint32Get(f.Country.ID)
			}})
	}

	if f.City.IsSet {
		jset.Add(&querier.Job{
			Name: "city_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.city.Uint32Get(f.City.ID)
			}})
	}

	if f.Birth.IsSet {
		jset.Add(&querier.Job{
			Name: "birth_year",
			Bitmap: func() *idx.Bitmap {

				return d.bitmap.birthYear.Uint32Get(uint32(f.Birth.Val))
			}})
	}

	if f.Interests.IsSet {
		jset.Add(&querier.Job{
			Name: "interests_eq",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.interest.Uint32Get(uint32(f.Interests.ID))
			}})
	}

	if f.Likes.IsSet {
		jset.Add(&querier.Job{
			Name: "likes_eq",
			Bitmap: func() *idx.Bitmap {
				l, ok := d.likes[uint32(f.Likes.Val)]
				if !ok {
					return nil
				}

				return idx.NewBitmapFromSlice(l)

			}})
	}

	if f.Joined.IsSet {
		jset.Add(&querier.Job{
			Name: "joined_year",
			Bitmap: func() *idx.Bitmap {
				return d.bitmap.joinedYear.Uint32Get(uint32(f.Joined.Val))
			}})
	}

}
