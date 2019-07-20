package requests

type FilterBits uint16

const (
	FilterSex FilterBits = 1 << iota
	FilterEmailDomain
	FilterStatus
	FilterFname
	FilterSname
	FilterPhoneCode
	FilterCountry
	FilterCity
	FilterBirthYear
	FilterPremiumNow
)

func Set(b, flag FilterBits) FilterBits    { return b | flag }
func Clear(b, flag FilterBits) FilterBits  { return b &^ flag }
func Toggle(b, flag FilterBits) FilterBits { return b ^ flag }
func Has(b, flag FilterBits) bool          { return b&flag != 0 }
