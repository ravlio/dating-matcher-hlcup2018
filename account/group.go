package account

import "github.com/ravlio/highloadcup2018/dicts"
import "github.com/ravlio/highloadcup2018/gojay"
import "sort"

type Group struct {
	Sex       dicts.Sex
	Status    dicts.Status
	Interests uint32
	Country   uint32
	City      uint32
	Count     uint32
}

type Groups []*Group

type RawGroup struct {
	Sex       string
	Status    string
	Interests string
	Country   string
	City      string
	Count     uint32
}

type RawGroups []*RawGroup

type GroupsContainer struct {
	Groups RawGroups
}

func MakeRawGroup(g *Group) *RawGroup {
	rg := &RawGroup{}

	if g.Sex > 0 {
		rg.Sex, _ = dicts.SexToString(g.Sex)
	}

	if g.Status > 0 {
		rg.Status, _ = dicts.StatusToString(g.Status)
	}

	if g.Interests != 0 {
		rg.Interests = dicts.Interest.GetKey(g.Interests)
		if rg.Interests == "" {
			println(g.Interests)
		}
	}

	if g.Country != 0 {
		rg.Country = dicts.Country.GetKey(g.Country)
	}

	if g.City != 0 {
		rg.City = dicts.City.GetKey(g.City)
	}

	rg.Count = g.Count

	return rg
}

func RawGroupSort(a RawGroups, s SortOrder, keys []string) {
	if s == SortAsc {
		sort.Slice(a, func(i, j int) bool {
			if a[i].Count < a[j].Count {
				return true
			}

			if a[i].Count > a[j].Count {
				return false
			}
			for _, k := range keys {
				switch k {
				case "sex":
					if a[i].Sex < a[j].Sex {
						return true
					}
					if a[i].Sex > a[j].Sex {
						return false
					}
				case "status":
					if a[i].Status < a[j].Status {
						return true
					}
					if a[i].Status > a[j].Status {
						return false
					}
				case "interests":
					if a[i].Interests < a[j].Interests {
						return true
					}
					if a[i].Interests > a[j].Interests {
						return false
					}
				case "country":
					if a[i].Country < a[j].Country {
						return true
					}

					if a[i].Country > a[j].Country {
						return false
					}
				case "city":
					if a[i].City < a[j].City {
						return true
					}

					if a[i].City > a[j].City {
						return false
					}
				}
			}

			return false
		})
	} else {
		sort.Slice(a, func(i, j int) bool {
			if a[i].Count > a[j].Count {
				return true
			}

			if a[i].Count < a[j].Count {
				return false
			}
			for _, k := range keys {
				switch k {
				case "sex":
					if a[i].Sex > a[j].Sex {
						return true
					}
					if a[i].Sex < a[j].Sex {
						return false
					}
				case "status":
					if a[i].Status > a[j].Status {
						return true
					}
					if a[i].Status < a[j].Status {
						return false
					}
				case "interests":
					if a[i].Interests > a[j].Interests {
						return true
					}
					if a[i].Interests < a[j].Interests {
						return false
					}
				case "country":
					if a[i].Country > a[j].Country {
						return true
					}

					if a[i].Country < a[j].Country {
						return false
					}
				case "city":
					if a[i].City > a[j].City {
						return true
					}

					if a[i].City < a[j].City {
						return false
					}
				}
			}

			return false
		})
	}
}

func GroupSort(a Groups, s SortOrder) {
	if s == SortAsc {
		sort.Slice(a, func(i, j int) bool {
			return a[i].Count < a[j].Count
		})
	} else {
		sort.Slice(a, func(i, j int) bool {
			return a[i].Count > a[j].Count
		})
	}
}

func (a *RawGroup) MarshalJSONObject(enc *gojay.Encoder) {
	if len(a.Sex) > 0 {
		enc.StringKey("sex", a.Sex)
	}

	if len(a.Status) > 0 {
		enc.StringKeyNoescape("status", escape(a.Status))
	}

	if len(a.Interests) > 0 {
		enc.StringKeyNoescape("interests", escape(a.Interests))
	}

	if len(a.Country) > 0 {
		enc.StringKeyNoescape("country", escape(a.Country))
	}

	if len(a.City) > 0 {
		enc.StringKeyNoescape("city", escape(a.City))
	}

	enc.Uint32Key("count", a.Count)
}

func (a *RawGroup) IsNil() bool {
	return a == nil
}

func (c *GroupsContainer) MarshalJSONArray(enc *gojay.Encoder) {
	for _, v := range c.Groups {
		enc.Object(v)
	}
}

func (c *GroupsContainer) IsNil() bool {
	return len(c.Groups) == 0
}
