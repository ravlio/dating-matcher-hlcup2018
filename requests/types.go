package requests

import (
	"errors"
	"github.com/ravlio/highloadcup2018/dicts"
	"github.com/ravlio/highloadcup2018/utils"
	"strconv"
	"strings"
)

type Int32 struct {
	Val   int32
	IsSet bool
	ID    uint32
}

func (i Int32) IsEqual(j Int32) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetInt32(b []byte, i *Int32) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	v, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return err
	}

	r := int32(v)
	if int64(r) != v {
		return errors.New("int32 overflow")
	}

	i.Val = r
	i.IsSet = true
	return nil
}

type Int64 struct {
	Val   int64
	IsSet bool
	ID    uint32
}

func (i Int64) IsEqual(j Int64) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetInt64(b []byte, i *Int64) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	i.Val = v
	i.IsSet = true
	return nil
}

type Uint8 struct {
	Val   uint8
	IsSet bool
	ID    uint32
}

func (i Uint8) IsEqual(j Uint8) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetUint8(b []byte, i *Uint8) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	v, err := strconv.ParseInt(string(b), 10, 8)
	if err != nil {
		return err
	}

	r := uint8(v)
	if int64(r) != v {
		return errors.New("uint8 overflow")
	}

	i.Val = r
	i.IsSet = true
	return nil
}

type Uint16 struct {
	Val   uint16
	IsSet bool
	ID    uint32
}

func (i Uint16) IsEqual(j Uint16) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetUint16(b []byte, i *Uint16) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	v, err := strconv.ParseInt(string(b), 10, 16)
	if err != nil {
		return err
	}

	r := uint16(v)
	if int64(r) != v {
		return errors.New("uint16 overflow")
	}

	i.Val = r
	i.IsSet = true
	return nil
}

type Uint32 struct {
	Val   uint32
	IsSet bool
	ID    uint32
}

func (i Uint32) IsEqual(j Uint32) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetUint32(b []byte, i *Uint32) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	v, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return err
	}

	r := uint32(v)
	if int64(r) != v {
		return errors.New("uint22 overflow")
	}

	i.Val = r
	i.IsSet = true
	return nil
}

type String struct {
	Val   string
	IsSet bool
	ID    uint32
}

func (i String) IsEqual(j String) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetString(b []byte, i *String, maxmin ...int) error {
	return CheckAndSetStringS(string(b), i, maxmin...)

}

func CheckAndSetStringS(b string, i *String, maxmin ...int) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	if maxmin != nil && len(b) > maxmin[0] {
		return errors.New("too long")
	}

	i.Val = string(b)
	i.IsSet = true
	return nil
}

type Uint32Array struct {
	Val   []uint32
	IsSet bool
}

func CheckAndSetUint32Array(b []byte, dist *Uint32Array) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}

	strarr := strings.Split(string(b), ",")
	if len(strarr) == 0 {
		return errors.New("empty val")
	}

	for _, v := range strarr {
		if len(strarr) == 0 {
			return errors.New("empty val")
		}

		d, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}

		r := uint32(d)
		if int64(r) != d {
			return errors.New("uint322 overflow")
		}

		dist.Val = append(dist.Val, r)
	}

	dist.IsSet = true

	return nil
}

type StringArray struct {
	Val   []string
	ID    []uint32
	IsSet bool
}

func CheckAndSetStringArray(b []byte, i *StringArray, maxmin ...int) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}

	v := strings.Split(string(b), ",")
	if len(v) == 0 {
		return errors.New("empty val")
	}

	for _, z := range v {
		if len(z) == 0 {
			return errors.New("empty val")
		}

		if maxmin != nil {
			if len(z) > maxmin[0] {
				return errors.New("too long")
			}
		}

		i.Val = append(i.Val, z)
	}

	i.IsSet = true

	return nil
}

type YesNo struct {
	Val   bool
	IsSet bool
}

func (i YesNo) IsEqual(j YesNo) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetYesNo(b []byte, i *YesNo) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	if len(b) > 1 {
		return errors.New("arg too long")
	}

	// если null=0, то есть выводить всех, у кого нулевое значение
	if b[0] == '0' {
		i.Val = true
	} else if b[0] == '1' {
		i.Val = false
	} else {
		return errors.New("wrong null")
	}
	i.IsSet = true

	return nil
}

type Order struct {
	Asc   bool
	IsSet bool
}

func (i Order) IsEqual(j Order) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Asc == j.Asc
}

func CheckAndSetOrder(b []byte, i *Order) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	if len(b) > 2 {
		return errors.New("arg too long")
	}

	if b[0] == '1' {
		i.Asc = true
	} else if b[0] == '-' && b[1] == '1' {
		i.Asc = false
	} else {
		return errors.New("wrong order")
	}
	i.IsSet = true

	return nil
}

type Phone struct {
	Val   string
	Int64 int64
	IsSet bool
	ID    uint32
}

func (i Phone) IsEqual(j Phone) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetPhone(b []byte, i *Phone) error {
	return CheckAndSetPhoneS(string(b), i)
}

// TODO придумать что-нибудь больше лучше. Хранить первой цифрой длину кода страны, второй - длину кода города
func CheckAndSetPhoneS(b string, i *Phone) error {
	if len(b) == 0 {
		return errors.New("empty phone")
	}

	phone, err := strconv.ParseUint(b[:1]+b[2:5]+b[6:], 10, 64)
	if err != nil {
		return errors.New("error parsing phone number")
	}
	/*a,err:=strconv.Atoi(p[:1])
	if err!=nil {
		return errors.New("wrong phone ")
	}
	b,err:=strconv.Atoi(p[1:4])
	if err!=nil {
		return 0,0,0
	}
	c,err:=strconv.Atoi(p[4:])
	if err!=nil {
		return 0,0,0
	}

	p1 := strings.SplitN(b, "(", 2)
	if len(p1) < 2 {
		return errors.New("wrong phone number")
	}
	p2 := strings.SplitN(p1[1], ")", 2)
	if len(p2) < 2 {
		return errors.New("wrong phone number")
	}

	if len(p1[0])!=1 {
		return errors.New("wrong phone country code length")
	}

	if len(p2[0])!=3 {
		return errors.New("wrong phone city code length")
	}

	if len(p2[1])!=7 {
		return errors.New("wrong phone length")
	}

	// TODO очень плохо тут поступаю
	phone,err:=strconv.ParseUint(p1[0]+p2[0]+p2[1], 10, 64)

	if err != nil {
		return errors.New("error parsing phone number")
	}*/

	i.IsSet = true
	i.Val = b
	i.Int64 = int64(phone)

	return nil
}

type Bool struct {
	Val   bool
	IsSet bool
}

func (i Bool) IsEqual(j Bool) bool {
	if i.IsSet != j.IsSet {
		return false
	}

	return i.Val == j.Val
}

func CheckAndSetBool(b []byte, i *Bool) error {
	if len(b) == 0 {
		return errors.New("empty arg")
	}
	if len(b) > 1 {
		return errors.New("arg too long")
	}

	if b[0] == '0' {
		i.Val = false
	} else if b[0] == '1' {
		i.Val = true
	} else {
		return errors.New("wrong bool")
	}
	i.IsSet = true

	return nil
}

func CheckAndSetEmail(b []byte, i *String) error {
	err := CheckAndSetString(b, i)
	if err != nil {
		return err
	}

	// TODO здесь валидатор емейла. Попроще и побыстрее
	return nil
}

func CheckAndSetEmailS(b string, i *String, maxmin ...int) error {
	if _, err := utils.ParseEmailDomain(string(b)); err != nil {
		return err
	}
	err := CheckAndSetStringS(b, i, maxmin...)
	if err != nil {
		return err
	}

	// TODO здесь валидатор емейла. Попроще и побыстрее
	return nil
}

func CheckAndSetEmailDomain(b []byte, i *String, maxmin ...int) error {
	err := CheckAndSetString(b, i, maxmin...)
	if err != nil {
		return err
	}

	// TODO здесь валидатор емейла. Попроще и побыстрее
	return nil
}

func CheckAndSetSex(b []byte, i *String) error {
	return CheckAndSetSexS(string(b), i)
}
func CheckAndSetSexS(b string, i *String) error {
	err := CheckAndSetStringS(b, i)
	if err != nil {
		return err
	}

	if _, err := dicts.StringToSex(i.Val); err != nil {
		return err
	}

	return nil
}

func CheckAndSetStatus(b []byte, i *String) error {
	return CheckAndSetStatusS(string(b), i)
}
func CheckAndSetStatusS(b string, i *String) error {
	err := CheckAndSetStringS(b, i)
	if err != nil {
		return err
	}

	if _, err := dicts.StringToStatus(i.Val); err != nil {
		return err
	}

	return nil
}
