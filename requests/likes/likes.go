package likes

import (
	"github.com/ravlio/highloadcup2018/errors"
	"github.com/ravlio/highloadcup2018/gojay"
)

type Likes []Like
type Like struct {
	Likee uint32
	Ts    int64
	Liker uint32
}

type LikesRequest struct {
	Likes Likes
}

func (p *Like) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "likee":
		err := dec.AddUint32(&p.Likee)
		if err != nil {
			return errors.Wrap(err, "error unmarshal like.likee")
		}

	case "ts":
		err := dec.AddInt64(&p.Ts)
		if err != nil {
			return errors.Wrap(err, "error unmarshal like.ts")
		}

	case "liker":
		err := dec.AddUint32(&p.Liker)
		if err != nil {
			return errors.Wrap(err, "error unmarshal like.liker")
		}
	default:
		return errors.New("unknown field")
	}

	return nil
}

func (p *Like) NKeys() int {
	return 3
}

func (i *Likes) UnmarshalJSONArray(dec *gojay.Decoder) error {
	l := Like{}
	if err := dec.Object(&l); err != nil {
		return err
	}
	*i = append(*i, l)

	return nil
}

func (p *LikesRequest) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	switch k {
	case "likes":
		if err := dec.Array(&p.Likes); err != nil {
			return err
		}
	default:
		return errors.New("unknown field")
	}

	return nil
}

func (p *LikesRequest) NKeys() int {
	return 1
}
