package account

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/gojay"
	"github.com/ravlio/highloadcup2018/requests"
	"github.com/ravlio/highloadcup2018/utils"
	"sort"
	"strconv"
	"sync"
	"time"
	// "github.com/valyala/fastjson"
)

type Like struct {
	ID uint32
	Ts uint32
}

type Interests []uint32
type Likes []Like

type Premium struct {
	Start  int32
	Finish int32
}

// Временная структура, призванная перекинуть некоторые строковые данные из запроса до индексатора
// которая после сохранения акка отправляется в GC
type Vars struct {
	Email string
}
type Account struct {
	ID          uint32
	Email       uint32
	EmailDomain uint32
	Fname       uint32
	Sname       uint32
	Phone       int64
	//PhoneCountryCode uint16
	//PhoneCode        uint16
	//Phone         uint32
	SexB          bool
	Sex           dicts.Sex
	Birth         int32
	Country       uint32
	City          uint32
	Joined        int32
	Interests     Interests
	PremiumStart  int32
	PremiumFinish int32
	Status        dicts.Status
	Likes         Likes
	Mx            sync.RWMutex
}

/*func (a *Account) FromFastJson(fj fastjson.Value) {
	if v:=fj.GetUint("id");v>0 {
		a.ID=uint32(v)
	}

	if v:=fj.GetStringBytes("email");len(v)>0 {
		a.Email = dicts.Email.GetOrCreateValue(string(v))
		domain, err := utils.ParseEmailDomain(a.Email)
		if err != nil {
			return errors.Wrap(err, "error parsing email domain")
		}
	}

	if v:=fj.GetStringBytes("fname");len(v)>0 {
		a.Fname = dicts.Fname.GetOrCreateValue(string(v))
	}

	if v:=fj.GetStringBytes("sname");len(v)>0 {
		a.Sname = dicts.Sname.GetOrCreateValue(string(v))
	}

	if v:=fj.GetStringBytes("phone");len(v)>0 {
		phone,err:=strconv.ParseUint(b[:1]+b[2:5]+b[6:], 10, 64)
		if err != nil {
			return errors.New("error parsing phone number")
		}

		a.Phone = int64(phone)
	}

	if v:=fj.GetStringBytes("sex");len(v)>0 {
		a.Sname,err = dicts.StringToSex(string(v))

		if err != nil {
			return errors.Wrap(err, "error parsing sex")
		}
	}

	if v:=fj.GetInt("birth");v!=0 {
		if v==0 {
		a.Birth=1
		} else {
			a.Birth=int32(v)
		}
	}

	if v:=fj.GetStringBytes("country");len(v)>0 {
		a.Country = dicts.Country.GetOrCreateValue(string(v))
	}

	if v:=fj.GetStringBytes("city");len(v)>0 {
		a.City = dicts.City.GetOrCreateValue(string(v))
	}

	if v:=fj.GetInt("joined");v!=0 {
		a.Joined=int32(v)
	}

	if v:=fj.GetArray("interests");len(v)>0 {
		a.Interests=make(Interests,len(v))
		for k,i:=range v {
			if s:=i.String();len(s)>0{
				a.Interests[k]=s
			} else {
				return errors.New("empty interest")
			}
		}

	}

	if v:=fj.GetObject("premium");v!=nil {
	)
		if s:=v.Get("start");s!=nil {
			i,err:=s.Int()
			if err!=nil {
				return err
			}

			a.PremiumStart=int32(i)
		}

		if s:=v.Get("finish");s!=nil {
			i,err:=s.Int()
			if err!=nil {
				return err
			}

			a.PremiumFinish=int32(i)
		}
	}
}*/
/*
Метод нужен, чтобы перед индексацией смерджить уже имеющийся аккаунт (там обычно только ид и емейл или телефон)
с остальными данными. В целях оптимизации ресурсов, данные мерджатся уже после обработки post-запроса,
много данных берётся из словарей.
*/
func (a *Account) MergeWithRequest(req *requests.AccountRequest) error {
	var err error

	if req.ID.IsSet {
		a.ID = req.ID.Val
	}

	if req.Email.IsSet {
		a.Email = dicts.Email.GetOrCreateValue(req.Email.Val)

		domain, err := utils.ParseEmailDomain(req.Email.Val)
		if err != nil {
			return errors.Wrap(err, "error parsing email domain")
		}

		a.EmailDomain = dicts.EmailDomain.GetOrCreateValue(domain)
	}
	if req.Fname.IsSet {
		a.Fname = dicts.Fname.GetOrCreateValue(req.Fname.Val)
	}

	if req.Sname.IsSet {
		a.Sname = dicts.Sname.GetOrCreateValue(req.Sname.Val)
	}

	if req.Phone.IsSet {
		a.Phone = req.Phone.Int64
	}

	if req.Sex.IsSet {
		a.Sex, err = dicts.StringToSex(req.Sex.Val)
		if err != nil {
			return errors.Wrap(err, "error parsing sex")
		}
	}

	if req.Birth.IsSet {
		if req.Birth.Val == 0 { // нулём может быть только отсутствие значения :)
			a.Birth = 1
		} else {
			a.Birth = req.Birth.Val
		}
	}

	if req.Country.IsSet {
		a.Country = dicts.Country.GetOrCreateValue(req.Country.Val)
	}

	if req.City.IsSet {
		a.City = dicts.City.GetOrCreateValue(req.City.Val)
	}

	if req.Joined.IsSet {
		a.Joined = req.Joined.Val
	}

	if len(req.Interests) > 0 {
		a.Interests = make([]uint32, len(req.Interests))
		for k, v := range req.Interests {
			a.Interests[k] = dicts.Interest.GetOrCreateValue(v)
		}
	}
	/*
		if len(req.Likes)>0 {
			a.Likes = make([]Like, len(req.Likes))
			for k, v := range req.Likes {
				a.Likes[k] = Like{Ts: uint32(v.Ts), ID: v.ID}
			}
		}*/

	if req.Premium.Start.IsSet {
		a.PremiumStart = req.Premium.Start.Val
	}

	if req.Premium.Finish.IsSet {
		a.PremiumFinish = req.Premium.Finish.Val
	}

	if req.Status.IsSet {
		a.Status, err = dicts.StringToStatus(req.Status.Val)

		if err != nil {
			return errors.Wrap(err, "error parsing status")
		}
	}

	return nil
}

/*func (a *Account) GetPhone() int64 {
	return utils.GetPhone(a.PhoneCountryCode, a.PhoneCode, a.Phone)
}*/

func (a *Account) GetBirthYear() uint32 {
	if a.Birth == 0 {
		return 0
	}

	return uint32(time.Unix(int64(a.Birth), 0).UTC().Year())
}

func (a *Account) GetJoinedYear() uint32 {
	if a.Joined == 0 {
		return 0
	}

	return uint32(time.Unix(int64(a.Joined), 0).UTC().Year())
}

// TODO do it
// Коряво, надо из инта тащить инт, без строк
func (a *Account) GetPhoneCode() uint16 {
	if a.Phone == 0 {
		return 0
	}
	s := strconv.Itoa(int(a.Phone))
	r, _ := strconv.Atoi(s[1:4])

	return uint16(r)
}

type SortOrder int

const (
	SortDesc = -1
	SortAsc  = 1
)

type Accounts []*Account

// Нужно, чтобы пробросить маршалеру необходимые поля
type AccountsContainer struct {
	Accounts Accounts
	Fields   []string
}

func Sort(a Accounts, s SortOrder) {
	if s == SortAsc {
		sort.Slice(a, func(i, j int) bool {
			return a[i].ID < a[j].ID
		})
	} else {
		sort.Slice(a, func(i, j int) bool {
			return a[i].ID > a[j].ID
		})
	}
}

func LeftLimit(a Accounts, l int) Accounts {
	if l > len(a) {
		return a
	}

	return a[:l]
}

func LeftLimitCopy(a Accounts, l int) Accounts {
	if l > len(a) {
		return a
	}

	r := make([]*Account, l)

	copy(r, a[:l])
	return r
}

func RightLimit(a Accounts, l int) Accounts {
	if l > len(a) {
		return a
	}

	return a[len(a)-l:]
}

func RightLimitIDs(a roaring.Arr, l int) roaring.Arr {
	if l > len(a) {
		return a
	}

	return a[len(a)-l:]
}

// Чтобы не изменять основной слайс
func RightLimitCopy(a Accounts, l int) Accounts {
	if l > len(a) {
		return a
	}

	r := make([]*Account, l)

	copy(r, a[len(a)-l:])

	return r

}

func Reverse(a Accounts) Accounts {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}

	return a
}

func ReverseIDs(a roaring.Arr) roaring.Arr {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}

	return a
}

func (i Interests) MarshalJSONArray(enc *gojay.Encoder) {
	for _, e := range i {
		enc.StringNoescape(escape(dicts.Interest.GetKey(e)))
	}
}

func (i Interests) IsNil() bool {
	return len(i) == 0
}

func (l Likes) MarshalJSONArray(enc *gojay.Encoder) {
	for _, v := range l {
		enc.Object(&v)
	}
}

func (l Likes) IsNil() bool {
	return len(l) == 0
}

func (l *Like) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Uint32Key("id", l.ID)
	enc.Uint32Key("ts", l.Ts)
}
func (l *Like) IsNil() bool {
	return l == nil
}

func (p *Premium) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Int32Key("start", p.Start)
	enc.Int32Key("finish", p.Finish)
}
func (p *Premium) IsNil() bool {
	return p == nil
}

// TODO найти нормальное решение
func escape(s string) string {
	r := strconv.QuoteToASCII(s)
	return r[1 : len(r)-1]
}

func (a *Account) MarshalJSONObject(enc *gojay.Encoder) {
	enc.Uint32Key("id", a.ID)

	if a.Birth != 0 {
		enc.Int32Key("birth", a.Birth)
	}

	if a.Email > 0 {
		enc.StringKey("email", dicts.Email.GetKey(a.Email))
	}

	if a.Sex > 0 {
		s, _ := dicts.SexToString(a.Sex)
		enc.StringKey("sex", s)
	}

	if a.Fname > 0 {
		enc.StringKeyNoescape("fname", escape(dicts.Fname.GetKey(a.Fname)))
	}
	if a.Sname > 0 {
		enc.StringKeyNoescape("sname", escape(dicts.Sname.GetKey(a.Sname)))
	}

	if a.Phone > 0 {
		enc.StringKey("phone", utils.PhoneToString(a.Phone))
	}

	if a.Country > 0 {
		enc.StringKeyNoescape("country", escape(dicts.Country.GetKey(a.Country)))
	}

	if a.City > 0 {
		enc.StringKeyNoescape("city", escape(dicts.City.GetKey(a.City)))
	}

	if a.Joined > 0 {
		enc.Int32Key("joined", a.Joined)
	}

	if len(a.Interests) > 0 {
		enc.AddArrayKey("interests", a.Interests)
	}

	if a.PremiumStart > 0 {
		p := &Premium{Start: a.PremiumStart, Finish: a.PremiumFinish}
		enc.ObjectKey("premium", p)
	}

	if len(a.Likes) > 0 {
		enc.AddArrayKey("likes", a.Likes)
	}

	if a.Status > 0 {
		s, _ := dicts.StatusToString(a.Status)
		enc.StringKeyNoescape("status", escape(s))
	}
}

func (a *Account) IsNil() bool {
	return a == nil
}

func (c *AccountsContainer) MarshalJSONArray(enc *gojay.Encoder) {
	for _, v := range c.Accounts {
		enc.ObjectWithKeys(v, c.Fields)
	}
}

func (c *AccountsContainer) IsNil() bool {
	return len(c.Accounts) == 0
}
