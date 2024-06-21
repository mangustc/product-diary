package localization

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmg-c/product-diary/logger"
)

const Divider = "|"

type (
	Msg    uint16
	Locale uint8
)

const (
	MsgErrorTest            Msg = 0
	MsgErrorInternalServer  Msg = 1
	MsgErrorGetUserNoInfo   Msg = 2
	MsgErrorGetUserNotFound Msg = 3
	MsgErrorCodeWrong       Msg = 4
	MsgErrorEmailWrong      Msg = 5
	MsgLoginInfoSent        Msg = 6
	MsgEmailPlaceholder     Msg = 7
	MsgCodePlaceholder      Msg = 8
	MsgSignIn               Msg = 9
	MsgLogIn                Msg = 10
	MsgFindUser             Msg = 11
	MsgUserControl          Msg = 12
	MsgUserList             Msg = 13
	MsgUserByID             Msg = 14
	MsgProfileInfo          Msg = 15
	MsgLogOut               Msg = 16
)

const (
	LocaleEnUS Locale = 0
	LocaleRuRU Locale = 1
)

type translateFunction func(locale Locale, args []string) string

var translations map[Msg]translateFunction = map[Msg]translateFunction{
	MsgErrorTest: func(locale Locale, args []string) string {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err.Error())
		}

		switch locale {
		case LocaleRuRU:
			if count == 1 {
				return fmt.Sprintf("Ошбика сервера %d", count)
			}
			return fmt.Sprintf("Ошибки сервера %d", count)
		default:
			if count == 1 {
				return fmt.Sprintf("Server error %d", count)
			}
			return fmt.Sprintf("Server errors %d", count)
		}
	},
	MsgErrorInternalServer: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Ошибка на сервере")
		default:
			return fmt.Sprintf("Internal server error")
		}
	},
	MsgErrorGetUserNoInfo: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Недостаточно информации о пользователе")
		default:
			return fmt.Sprintf("Not enough info about user")
		}
	},
	MsgErrorGetUserNotFound: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Такой пользователь не был найден")
		default:
			return fmt.Sprintf("No such user found")
		}
	},
	MsgErrorCodeWrong: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Неправильный код подтверждения")
		default:
			return fmt.Sprintf("Wrong confirmation code")
		}
	},
	MsgErrorEmailWrong: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Данная электронная почта не существует")
		default:
			return fmt.Sprintf("Given email doesn't exist")
		}
	},
	MsgLoginInfoSent: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Данные для входа отправлены на почту %s", args[0])
		default:
			return fmt.Sprintf("Login information has been sent to %s", args[0])
		}
	},
	MsgEmailPlaceholder: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Почта")
		default:
			return fmt.Sprintf("Email")
		}
	},
	MsgCodePlaceholder: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Код")
		default:
			return fmt.Sprintf("Code")
		}
	},
	MsgSignIn: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Зарегестрироваться")
		default:
			return fmt.Sprintf("Sign in")
		}
	},
	MsgLogIn: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Войти")
		default:
			return fmt.Sprintf("Log in")
		}
	},
	MsgFindUser: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Найти пользователя")
		default:
			return fmt.Sprintf("Find User")
		}
	},
	MsgUserControl: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Панель управления пользователями")
		default:
			return fmt.Sprintf("User control")
		}
	},
	MsgUserList: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Список пользователей")
		default:
			return fmt.Sprintf("User list")
		}
	},
	MsgUserByID: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Найти пользователя по ID")
		default:
			return fmt.Sprintf("User by ID")
		}
	},
	MsgProfileInfo: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Профиль")
		default:
			return fmt.Sprintf("Profile")
		}
	},
}

func Localize(msg string, locale Locale) string {
	parts := strings.Split(msg, Divider)

	msgID64, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		logger.Error.Printf("No parsing message ID. msg: %s.", msg)
		return msg
	}
	msgID := Msg(uint16(msgID64))

	translate, found := translations[msgID]
	if !found {
		logger.Error.Printf("No such message ID found. msg: %s.", msg)
		return msg
	}
	return translate(locale, parts[1:])
}

func GetMessage(msgID Msg, args ...any) string {
	var out string = fmt.Sprint(msgID)
	for _, arg := range args {
		out += Divider + fmt.Sprint(arg)
	}
	return out
}

func GetError(msgID Msg, args ...any) error {
	return fmt.Errorf(GetMessage(msgID, args))
}

func GetLocalized(locale Locale, msgID Msg, args ...any) string {
	return Localize(GetMessage(msgID, args), locale)
}

type Localizer struct {
	locale Locale
}

func NewLocilizer(locale Locale) *Localizer {
	return &Localizer{
		locale: locale,
	}
}

func (l *Localizer) GetLocalized(msgID Msg, args ...any) string {
	return GetLocalized(l.locale, msgID, args...)
}

func (l *Localizer) Localize(msg string) string {
	return Localize(msg, l.locale)
}
