package requests

import (
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/gojay"
	"github.com/valyala/fastjson"
	"strconv"
	"sync"
)

type Like struct {
	ID uint32
	Ts int64
}

type Premium struct {
	Start  Int32
	Finish Int32
}

type Interests []string
type Likes []Like

type AccountRequest struct {
	ID        Uint32
	Email     String
	Fname     String
	Sname     String
	Phone     Phone
	Sex       String
	Birth     Int32
	Country   String
	City      String
	Joined    Int32
	Interests Interests
	Premium   Premium
	Likes     Likes
	Status    String
}

var accountRequestPool = sync.Pool{
	New: func() interface{} {
		return &AccountRequest{
			Interests: make(Interests, 0, 20),
			Likes:     make(Likes, 0, 100),
		}
	},
}

func (a *AccountRequest) ReleaseToPool() {
	a.Birth.IsSet = false
	a.City.IsSet = false
	a.Country.IsSet = false
	a.Email.IsSet = false
	a.Fname.IsSet = false
	a.ID.IsSet = false
	a.Joined.IsSet = false
	a.Phone.IsSet = false
	a.Sex.IsSet = false
	a.Sname.IsSet = false
	a.Status.IsSet = false
	if len(a.Interests) > 0 {
		a.Interests = a.Interests[:0]
	}

	if len(a.Likes) > 0 {
		a.Likes = a.Likes[:0]
	}

	a.Premium.Start.IsSet = false
	a.Premium.Finish.IsSet = false

	accountRequestPool.Put(a)
}

func AccountRequestPoolGet() *AccountRequest {
	return accountRequestPool.Get().(*AccountRequest)
}

func (p *Premium) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "start":
		err := dec.AddInt32(&p.Start.Val)
		if err != nil {
			return errors.Wrap(err, "error unmarshal premium.start")
		}

		if p.Start.Val > 0 {
			p.Start.IsSet = true
		}
	case "finish":
		err := dec.AddInt32(&p.Finish.Val)
		if err != nil {
			return errors.Wrap(err, "error unmarshal premium.finish")
		}

		if p.Finish.Val > 0 {
			p.Finish.IsSet = true
		}
	default:
		return errors.New("unknown field")
	}
	return nil
}

func (p *Premium) NKeys() int {
	return 2
}

func (i *Interests) UnmarshalJSONArray(dec *gojay.Decoder) error {
	str := ""
	if err := dec.String(&str); err != nil {
		return err
	}

	if len(str) == 0 {
		return errors.New("empty interest")
	}

	if len(str) > 100 {
		return errors.New("interest too long")
	}
	*i = append(*i, str)

	return nil
}

func (i *Likes) UnmarshalJSONArray(dec *gojay.Decoder) error {
	l := Like{}
	if err := dec.Object(&l); err != nil {
		return err
	}
	*i = append(*i, l)

	return nil
}

func (i *Like) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "id":
		err := dec.AddUint32(&i.ID)
		if err != nil {
			return errors.Wrap(err, "error unmarshal like.id")
		}

		if i.ID <= 0 {
			return errors.New("wrong like.id")
		}
	case "ts":
		err := dec.AddInt64(&i.Ts)
		if err != nil {
			return errors.Wrap(err, "error unmarshal like.ts")
		}
	default:
		return errors.New("unknown field")
	}
	return nil
}

func (i *Like) NKeys() int {
	return 2
}

func (a *AccountRequest) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "id":
		a.ID.IsSet = true
		err := dec.AddUint32(&a.ID.Val)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.id")
		}

	case "email":
		// TODO убрать лишние аллокации, кидать сразу в a.Email.Val
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.email")
		}
		err = CheckAndSetEmailS(s, &a.Email, 100)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetEmailS account.email")
		}

	case "fname":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.fname")
		}
		err = CheckAndSetStringS(s, &a.Fname, 50)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetStringS account.fname")
		}

	case "sname":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.sname")
		}

		err = CheckAndSetStringS(s, &a.Sname, 50)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetStringS account.sname")
		}

	case "phone":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.phone")
		}

		err = CheckAndSetPhoneS(s, &a.Phone)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetPhoneS account.phone")
		}
	case "sex":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.sex")
		}

		err = CheckAndSetSexS(s, &a.Sex)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetSexS account.sex")
		}

	case "birth":
		a.Birth.IsSet = true
		err := dec.AddInt32(&a.Birth.Val)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.birth")
		}

	case "country":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.country")
		}
		err = CheckAndSetStringS(s, &a.Country, 50)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetStringS account.country")
		}

	case "city":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.city")
		}
		err = CheckAndSetStringS(s, &a.City, 50)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetStringS account.city")
		}

	case "joined":
		a.Joined.IsSet = true

		err := dec.AddInt32(&a.Joined.Val)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.joined")
		}

	case "interests":
		err := dec.DecodeArray(&a.Interests)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.interests")
		}

	case "likes":
		err := dec.DecodeArray(&a.Likes)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.likes")
		}

	case "premium":
		if err := dec.Object(&a.Premium); err != nil {
			return errors.Wrap(err, "error unmarshal account.premium")
		}

	case "status":
		var s string
		err := dec.AddString(&s)
		if err != nil {
			return errors.Wrap(err, "error unmarshal account.status")
		}
		err = CheckAndSetStatusS(s, &a.Status)
		if err != nil {
			return errors.Wrap(err, "error CheckAndSetStatusS account.status")
		}
	default:
		return errors.New("unknown field")

	}
	return nil
}

func (a *AccountRequest) NKeys() int {
	return 14
}

func (a *AccountRequest) FromFastJson(fj *fastjson.Value) error {
	if v := fj.GetUint("id"); v > 0 {
		a.ID.IsSet = true
		a.ID.Val = uint32(v)
	}

	if v := fj.GetStringBytes("email"); len(v) > 0 {
		a.Email.IsSet = true
		a.Email.Val = string(v)
	}

	if v := fj.GetStringBytes("fname"); len(v) > 0 {
		a.Fname.IsSet = true
		a.Fname.Val = string(v)
	}

	if v := fj.GetStringBytes("sname"); len(v) > 0 {
		a.Sname.IsSet = true
		a.Sname.Val = string(v)
	}

	if v := fj.GetStringBytes("phone"); len(v) > 0 {
		a.Phone.IsSet = true
		a.Phone.Val = string(v)
		phone, err := strconv.ParseUint(a.Phone.Val[:1]+a.Phone.Val[2:5]+a.Phone.Val[6:], 10, 64)
		if err != nil {
			return errors.New("error parsing phone number")
		}

		a.Phone.Int64 = int64(phone)
	}

	if v := fj.GetStringBytes("sex"); len(v) > 0 {
		a.Sex.IsSet = true
		a.Sex.Val = string(v)
	}

	if v := fj.GetInt("birth"); v != 0 {
		a.Birth.IsSet = true
		a.Birth.Val = int32(v)
	}

	if v := fj.GetStringBytes("country"); len(v) > 0 {
		a.Country.IsSet = true
		a.Country.Val = string(v)
	}

	if v := fj.GetStringBytes("city"); len(v) > 0 {
		a.City.IsSet = true
		a.City.Val = string(v)
	}

	if v := fj.GetInt("joined"); v != 0 {
		a.Joined.IsSet = true
		a.Joined.Val = int32(v)
	}

	if v := fj.GetArray("interests"); len(v) > 0 {
		for _, i := range v {
			sb, err := i.StringBytes()
			if err != nil {
				return err
			}

			if len(sb) > 0 {
				a.Interests = append(a.Interests, string(sb))
			} else {
				return errors.New("empty interest")
			}
		}

	}

	if v := fj.GetArray("likes"); len(v) > 0 {
		for _, i := range v {
			l, err := i.Object()
			if err != nil {
				return err
			}

			like := Like{}
			if s := l.Get("id"); s != nil {
				i, err := s.Uint()
				if err != nil {
					return err
				}

				like.ID = uint32(i)
			}

			if s := l.Get("ts"); s != nil {
				i, err := s.Int()
				if err != nil {
					return err
				}

				like.Ts = int64(i)
			}

			a.Likes = append(a.Likes, like)
		}
	}

	if v := fj.GetObject("premium"); v != nil {
		if s := v.Get("start"); s != nil {
			i, err := s.Int()
			if err != nil {
				return err
			}
			a.Premium.Start.IsSet = true
			a.Premium.Start.Val = int32(i)
		}

		if s := v.Get("finish"); s != nil {
			i, err := s.Int()
			if err != nil {
				return err
			}

			a.Premium.Finish.IsSet = true
			a.Premium.Finish.Val = int32(i)
		}
	}

	if v := fj.GetStringBytes("status"); len(v) > 0 {
		a.Status.IsSet = true
		a.Status.Val = string(v)
	}

	return nil
}
