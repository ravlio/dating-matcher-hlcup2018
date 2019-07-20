package requests

import "github.com/valyala/fasthttp"
import "github.com/ravlio/highloadcup2018/errors"
import "github.com/ravlio/highloadcup2018/utils"
import "github.com/ravlio/highloadcup2018/dicts"
import "bytes"
import "strings"
import "strconv"

type GroupKeys uint16

const (
	GroupSex GroupKeys = 1 << iota
	GroupStatus
	GroupInterests
	GroupCountry
	GroupCity
)

func (g *GroupKeys) Set(flag GroupKeys)     { *g = *g | flag }
func (g *GroupKeys) Clear(flag GroupKeys)   { *g = *g &^ flag }
func (g *GroupKeys) Toggle(flag GroupKeys)  { *g = *g ^ flag }
func (g GroupKeys) Has(flag GroupKeys) bool { return g&flag != 0 }

type Group struct {
	Sex    String
	Status String

	Country String
	City    String

	Birth     Uint16
	Interests String
	Likes     Uint32
	Joined    Uint16
	Order     Order
	Limit     Uint8
	Keys      GroupKeys
	KeysRaw   []string
	QueryID   uint32
}

func (f *Group) FillRequest(ctx *fasthttp.RequestCtx) (error, int) {
	if !ctx.Request.URI().QueryArgs().Has("limit") {
		return errors.New("empty 'limit' arg"), fasthttp.StatusBadRequest
	}

	if !ctx.Request.URI().QueryArgs().Has("order") {
		return errors.New("empty 'order' arg"), fasthttp.StatusBadRequest
	}

	if !ctx.Request.URI().QueryArgs().Has("keys") {
		return errors.New("empty 'keys' arg"), fasthttp.StatusBadRequest
	}

	for _, arg := range utils.GetURIArgs(ctx) {
		if bytes.Compare(arg.Key, []byte(`keys`)) == 0 {
			if len(arg.Key) == 0 {
				return errors.New("empty key"), fasthttp.StatusBadRequest
			}

			v := strings.Split(string(arg.Value), ",")
			f.KeysRaw = v
			for _, z := range v {
				kp := &f.Keys

				switch z {
				case "sex":
					kp.Set(GroupSex)
				case "status":
					kp.Set(GroupStatus)
				case "interests":
					kp.Set(GroupInterests)
				case "country":
					kp.Set(GroupCountry)
				case "city":
					kp.Set(GroupCity)
				default:
					return errors.New("wrong key error"), fasthttp.StatusBadRequest
				}
			}

			continue
		}

		if bytes.Compare(arg.Key, []byte(`sex`)) == 0 {
			if err := CheckAndSetSex(arg.Value, &f.Sex); err != nil {
				return errors.Wrap(err, "sex arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`status`)) == 0 {
			if err := CheckAndSetStatus(arg.Value, &f.Status); err != nil {
				return errors.Wrap(err, "status arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`country`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.Country, 50); err != nil {
				return errors.Wrap(err, "country arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`city`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.City, 50); err != nil {
				return errors.Wrap(err, "city arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`birth`)) == 0 {
			if err := CheckAndSetUint16(arg.Value, &f.Birth); err != nil {
				return errors.Wrap(err, "birth arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`interests`)) == 0 {
			if err := CheckAndSetString(arg.Value, &f.Interests, 100); err != nil {
				return errors.Wrap(err, "interests arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`likes`)) == 0 {
			if err := CheckAndSetUint32(arg.Value, &f.Likes); err != nil {
				return errors.Wrap(err, "likes arg error"), fasthttp.StatusBadRequest
			}
			continue
		}
		if bytes.Compare(arg.Key, []byte(`joined`)) == 0 {
			if err := CheckAndSetUint16(arg.Value, &f.Joined); err != nil {
				return errors.Wrap(err, "joined arg error"), fasthttp.StatusBadRequest
			}
			continue
		}

		if bytes.Compare(arg.Key, []byte(`order`)) == 0 {
			if err := CheckAndSetOrder(arg.Value, &f.Order); err != nil {
				return errors.Wrap(err, "order arg error"), fasthttp.StatusBadRequest
			}
			continue
		}

		if bytes.Compare(arg.Key, []byte(`limit`)) == 0 {
			if err := CheckAndSetUint8(arg.Value, &f.Limit); err != nil {
				return errors.Wrap(err, "limit arg error"), fasthttp.StatusBadRequest
			}

			if f.Limit.Val > 50 {
				return errors.New("expected limit <=50"), fasthttp.StatusBadRequest
			}
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
		return errors.New("unknown group query arg"), fasthttp.StatusBadRequest
	}

	return nil, 0
}

func (f *Group) FillWithDicts() error {
	var ok bool

	if f.Sex.IsSet {
		var err error
		i, err := dicts.StringToSex(f.Sex.Val)
		if err != nil {
			return errors.Wrap(err, "error filling sex")
		}

		f.Sex.ID = uint32(i)
	}

	if f.Status.IsSet {
		var err error
		i, err := dicts.StringToStatus(f.Status.Val)
		if err != nil {
			return errors.Wrap(err, "error filling status")
		}

		f.Status.ID = uint32(i)
	}

	if f.Country.IsSet {
		var ok bool
		if f.Country.ID, ok = dicts.Country.GetValue(f.Country.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in country dict")
		}
	}

	if f.City.IsSet {
		var ok bool
		if f.City.ID, ok = dicts.City.GetValue(f.City.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in city dict")
		}
	}

	if f.Interests.IsSet {
		if f.Interests.ID, ok = dicts.Interest.GetValue(f.Interests.Val); !ok {
			return errors.Wrap(ErrNotFoundInDict, "not found in interests dict")
		}
	}

	return nil
}
