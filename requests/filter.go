package requests

import (
	"bytes"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/utils"
	"github.com/valyala/fasthttp"
	"strconv"
)

var ErrNotFoundInDict = errors.New("not found in dict")

type FilterOrder int8

const (
	OrderAsc  FilterOrder = 1
	OrderDesc FilterOrder = -1
)

type Filter struct {
	SexEq       String
	EmailEq     String
	EmailDomain String
	EmailLt     String
	EmailGt     String
	StatusEq    String
	StatusNeq   String

	FnameEq   String
	FnameAny  StringArray
	FnameNull YesNo

	SnameEq     String
	SnameStarts String
	SnameNull   YesNo

	PhoneEq   Phone
	PhoneCode Uint16
	PhoneNull YesNo

	CountryEq   String
	CountryNull YesNo
	CityEq      String
	CityAny     StringArray
	CityNull    YesNo

	BirthLt           Int64
	BirthGt           Int64
	BirthYear         Uint16
	InterestsEq       String
	InterestsContains StringArray
	InterestsAny      StringArray

	LikesContains Uint32Array
	PremiumNow    Bool
	PremiumNull   YesNo
	Limit         Uint8
	filledFields  []string // поля, которые будут выводиться в ответном json
	fieldsMap     map[string]struct{}
	QueryID       uint32
}

func (f *Filter) FilledCount() int {
	return len(f.filledFields)
}

func (f *Filter) FilledFields() []string {
	return f.filledFields
}

func (f *Filter) FieldsMap() map[string]struct{} {
	return f.fieldsMap
}

func appendField(dst []string, f string) []string {
	for _, v := range dst {
		if v == f {
			return dst
		}
	}

	return append(dst, f)
}

func (f *Filter) FillRequest(ctx *fasthttp.RequestCtx) (error, int) {
	if !ctx.Request.URI().QueryArgs().Has("limit") {
		return errors.New("empty 'limit' arg"), fasthttp.StatusBadRequest
	}

	ff := make([]string, 0)

	for _, arg := range utils.GetURIArgs(ctx) {
		if bytes.Compare(arg.Key, []byte(`sex_eq`)) == 0 {
			if err := CheckAndSetSex(arg.Value, &f.SexEq); err != nil {
				return errors.Wrap(err, "sex_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "sex")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`email_domain`)) == 0 {
			if err := CheckAndSetEmailDomain(arg.Value, &f.EmailDomain, 100); err != nil {
				return errors.Wrap(err, "email_domain arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "email")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`email_lt`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.EmailLt, 100); err != nil {
				return errors.Wrap(err, "email_lt arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "email")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`email_gt`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.EmailGt, 100); err != nil {
				return errors.Wrap(err, "email_gt arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "email")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`status_eq`)) == 0 {
			if err := CheckAndSetStatus(arg.Value, &f.StatusEq); err != nil {
				return errors.Wrap(err, "status_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "status")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`status_neq`)) == 0 {
			if err := CheckAndSetStatus(arg.Value, &f.StatusNeq); err != nil {
				return errors.Wrap(err, "status_neq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "status")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`fname_eq`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.FnameEq, 50); err != nil {
				return errors.Wrap(err, "fname_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "fname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`fname_any`)) == 0 {
			if err := CheckAndSetStringArray(arg.Value, &f.FnameAny, 50); err != nil {
				return errors.Wrap(err, "fname_any arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "fname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`fname_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.FnameNull); err != nil {
				return errors.Wrap(err, "fname_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "fname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`sname_eq`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.SnameEq, 50); err != nil {
				return errors.Wrap(err, "sname_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "sname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`sname_starts`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.SnameStarts, 50); err != nil {
				return errors.Wrap(err, "sname_starts arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "sname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`sname_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.SnameNull); err != nil {
				return errors.Wrap(err, "sname_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "sname")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`phone_code`)) == 0 {
			if err := CheckAndSetUint16(arg.Value, &f.PhoneCode); err != nil {
				return errors.Wrap(err, "phone_code arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "phone")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`phone_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.PhoneNull); err != nil {
				return errors.Wrap(err, "phone_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "phone")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`country_eq`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.CountryEq, 50); err != nil {
				return errors.Wrap(err, "country_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "country")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`country_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.CountryNull); err != nil {
				return errors.Wrap(err, "country_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "country")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`city_eq`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.CityEq, 50); err != nil {
				return errors.Wrap(err, "city_eq arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "city")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`city_any`)) == 0 {
			if err := CheckAndSetStringArray(arg.Value, &f.CityAny, 50); err != nil {
				return errors.Wrap(err, "city_any arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "city")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`city_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.CityNull); err != nil {
				return errors.Wrap(err, "city_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "city")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`birth_lt`)) == 0 {
			if err := CheckAndSetInt64(arg.Value, &f.BirthLt); err != nil {
				return errors.Wrap(err, "birth_lt arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "birth")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`birth_gt`)) == 0 {
			if err := CheckAndSetInt64(arg.Value, &f.BirthGt); err != nil {
				return errors.Wrap(err, "birth_gt arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "birth")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`birth_year`)) == 0 {
			if err := CheckAndSetUint16(arg.Value, &f.BirthYear); err != nil {
				return errors.Wrap(err, "birth_year arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "birth")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`interests_contains`)) == 0 {
			if err := CheckAndSetStringArray(arg.Value, &f.InterestsContains, 100); err != nil {
				return errors.Wrap(err, "interests_contains arg error"), fasthttp.StatusBadRequest
			}
			// ff=appendField(ff,"interests")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`interests_any`)) == 0 {
			if err := CheckAndSetStringArray(arg.Value, &f.InterestsAny, 100); err != nil {
				return errors.Wrap(err, "sex_eq arg error"), fasthttp.StatusBadRequest
			}
			// ff=appendField(ff,"interests")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`likes_contains`)) == 0 {
			if err := CheckAndSetUint32Array(arg.Value, &f.LikesContains); err != nil {
				return errors.Wrap(err, "likes_contains arg error"), fasthttp.StatusBadRequest
			}
			// ff=appendField(ff,"likes")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`premium_now`)) == 0 {
			if err := CheckAndSetBool(arg.Value, &f.PremiumNow); err != nil {
				return errors.Wrap(err, "premium_now arg error"), fasthttp.StatusBadRequest
			}

			if !f.PremiumNow.Val {
				return errors.New("expected premium=0"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "premium")
			ff = appendField(ff, "start") // костыльненько определяем, что поля внутри premium{} также можно энкодить json-энкодеру
			ff = appendField(ff, "finish")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`premium_null`)) == 0 {
			if err := CheckAndSetYesNo(arg.Value, &f.PremiumNull); err != nil {
				return errors.Wrap(err, "premium_null arg error"), fasthttp.StatusBadRequest
			}
			ff = appendField(ff, "premium")
			ff = appendField(ff, "start")
			ff = appendField(ff, "finish")
			continue
		}
		if bytes.Compare(arg.Key, []byte(`limit`)) == 0 {
			if err := CheckAndSetUint8(arg.Value, &f.Limit); err != nil {
				return errors.Wrap(err, "limit arg error"), fasthttp.StatusBadRequest
			}

			if f.Limit.Val > 50 {
				return errors.New("expected limit <=50"), fasthttp.StatusBadRequest
			}
			// ff=append(ff,"limit")
			continue
		}

		if bytes.Compare(arg.Key, []byte(`query_id`)) == 0 {
			v, err := strconv.ParseInt(string(arg.Value), 10, 32)
			if err != nil {
				continue
			}

			r := uint32(v)
			if int64(r) != v {
				continue
			}
			f.QueryID = r

			continue
		}

		// неизвестный аргумент
		return errors.New("unknown filter query arg"), fasthttp.StatusBadRequest
	}

	ff = appendField(ff, "email")
	ff = appendField(ff, "id")
	f.filledFields = ff
	/*f.fieldsMap=make(map[string]struct{})
	for _,v:=range ff {
		f.fieldsMap[v]=struct{}{}
	}*/

	return nil, 0
}

/* Заполняем словари на все предикты, определяющие полные соответствия
Если в словаре пусто, значит и фильтр ничего не найдёт, можно сразу возвращать 404
Сразу не заполняем словари, так как в конце может возникнуть ошибка валидации, а мы уже потратили время на хождение
в словари */
func (f *Filter) FillWithDicts() error {

	if f.SexEq.IsSet {
		var err error
		i, err := dicts.StringToSex(f.SexEq.Val)
		if err != nil {
			return errors.Wrap(err, "error filling sex")
		}

		f.SexEq.ID = uint32(i)
	}
	if f.EmailEq.IsSet {
		var ok bool
		if f.EmailEq.ID, ok = dicts.Email.GetValue(f.EmailEq.Val); !ok {
			// Если в словаре ничего нет, то сразу возвращаем 404
			return errors.Wrap(ErrNotFoundInDict, "not found in email dict")
		}
	}
	if f.EmailDomain.IsSet {
		var ok bool
		if f.EmailDomain.ID, ok = dicts.EmailDomain.GetValue(f.EmailDomain.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in emailDomain dict")
		}
	}

	if f.StatusEq.IsSet {
		var err error
		i, err := dicts.StringToStatus(f.StatusEq.Val)
		if err != nil {
			return errors.Wrap(err, "error filling status")
		}

		f.StatusEq.ID = uint32(i)
	}

	if f.StatusNeq.IsSet {
		var err error
		i, err := dicts.StringToStatus(f.StatusNeq.Val)
		if err != nil {
			return errors.Wrap(err, "error filling status")
		}

		f.StatusNeq.ID = uint32(i)
	}

	if f.FnameEq.IsSet {
		var ok bool
		if f.FnameEq.ID, ok = dicts.Fname.GetValue(f.FnameEq.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in fname dict")
		}
	}

	if f.FnameAny.IsSet {
		f.FnameAny.ID = make([]uint32, 0, len(f.FnameAny.Val))
		i := 0
		for _, v := range f.FnameAny.Val {
			fn, ok := dicts.Fname.GetValue(v)
			if ok {
				f.FnameAny.ID = append(f.FnameAny.ID, fn)
				i++
			}
		}

		if len(f.FnameAny.ID) == 0 {
			return errors.Wrap(ErrNotFoundInDict, "no fname was found")
		}
		f.FnameAny.ID = f.FnameAny.ID[:i]
	}

	if f.CountryEq.IsSet {
		var ok bool
		if f.CountryEq.ID, ok = dicts.Country.GetValue(f.CountryEq.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in country dict")
		}
	}

	if f.CityAny.IsSet {
		f.CityAny.ID = make([]uint32, 0, len(f.CityAny.Val))
		i := 0
		for _, v := range f.CityAny.Val {
			fn, ok := dicts.City.GetValue(v)
			if ok {
				f.CityAny.ID = append(f.CityAny.ID, fn)
				i++
			}
		}

		if len(f.CityAny.ID) == 0 {
			return errors.Wrap(ErrNotFoundInDict, "no city was found")
		}

		f.CityAny.ID = f.CityAny.ID[:i]
	}

	if f.CityEq.IsSet {
		var ok bool
		if f.CityEq.ID, ok = dicts.City.GetValue(f.CityEq.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in city dict")
		}
	}

	if f.InterestsContains.IsSet {
		f.InterestsContains.ID = make([]uint32, len(f.InterestsContains.Val))
		for k, v := range f.InterestsContains.Val {
			fn, ok := dicts.Interest.GetValue(v)
			if !ok {
				return errors.Wrap(ErrNotFoundInDict, "not found in interest dict")
			}

			f.InterestsContains.ID[k] = fn
		}
	}

	if f.InterestsAny.IsSet {
		f.InterestsAny.ID = make([]uint32, 0, len(f.InterestsAny.Val))
		i := 0
		for _, v := range f.InterestsAny.Val {
			fn, ok := dicts.Interest.GetValue(v)
			if ok {
				f.InterestsAny.ID = append(f.InterestsAny.ID, fn)
				i++
			}
		}

		if len(f.InterestsAny.ID) == 0 {
			return errors.Wrap(ErrNotFoundInDict, "no interest was found")
		}
		f.InterestsAny.ID = f.InterestsAny.ID[:i]
	}

	return nil
}
