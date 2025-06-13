package language

import "notifications/pkg/util/strset"

const (
	RU = "ru"
	TJ = "tg"
	UZ = "uz"
	EN = "en"
)

type Language struct {
	RU string `json:"ru,omitempty"`
	TJ string `json:"tg,omitempty"`
	UZ string `json:"uz,omitempty"`
	EN string `json:"en,omitempty"`
}

type Fields struct {
	Key string
	Val string
}

func GetAll() []string {
	return []string{RU, EN, TJ, UZ}
}

func (l *Language) ValidAny() bool {
	return !strset.IsEmpty(l.RU) || !strset.IsEmpty(l.TJ) || !strset.IsEmpty(l.UZ) || !strset.IsEmpty(l.EN)
}

func (l *Language) Get(lang string) string {
	switch lang {
	case RU:
		return l.RU
	case TJ:
		return l.TJ
	case UZ:
		return l.UZ
	case EN:
		return l.EN
	default:
		return ""
	}
}

func (l *Language) GetAll() []string {
	return []string{l.RU, l.TJ, l.UZ, l.EN}
}

func (l *Language) GetAllWithLang() []Fields {
	return []Fields{
		{Key: RU, Val: l.RU},
		{Key: TJ, Val: l.TJ},
		{Key: UZ, Val: l.UZ},
		{Key: EN, Val: l.EN},
	}
}

func (l *Language) Set(lang, value string) {
	switch lang {
	case RU:
		l.RU = value
	case TJ:
		l.TJ = value
	case UZ:
		l.UZ = value
	case EN:
		l.EN = value
	}
}

func (l *Language) SetAll(value string) {
	l.RU = value
	l.TJ = value
	l.UZ = value
	l.EN = value
}

func (l *Language) Valid() {
	if strset.IsEmpty(l.TJ) {
		l.TJ = l.RU
	}
	if strset.IsEmpty(l.UZ) {
		l.UZ = l.RU
	}
	if strset.IsEmpty(l.EN) {
		l.EN = l.RU
	}
}
