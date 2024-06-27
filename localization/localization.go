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
	MsgErrorTest Msg = iota
	MsgErrorInternalServer
	MsgErrorGetUserNoInfo
	MsgErrorGetUserNotFound
	MsgErrorCodeWrong
	MsgErrorEmailWrong
	MsgLoginInfoSent
	MsgEmailPlaceholder
	MsgCodePlaceholder
	MsgSignIn
	MsgLogIn
	MsgFindUser
	MsgUserControl
	MsgUserList
	MsgUserByID
	MsgProfileInfo
	MsgLogOut
	MsgErrorGetSessionNotFound
	MsgErrorPasswordWrong
	MsgEnglish
	MsgRussian
	MsgEmailExists
	MsgUnhide
	MsgHide
	MsgPersons
	MsgUsername
	MsgAdd
	MsgErrorUsernameEmpty
	MsgErrorUsernameAlreadyExists
	MsgErrorProductTitle
	MsgErrorProductCalories
	MsgErrorProductNutrient
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
			return fmt.Sprintf("Зарегистрироваться")
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
	MsgLogOut: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Выйти")
		default:
			return fmt.Sprintf("Logout")
		}
	},
	MsgErrorGetSessionNotFound: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Пользовательская сессия устарела")
		default:
			return fmt.Sprintf("User session expired")
		}
	},
	MsgErrorPasswordWrong: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Неверный пароль")
		default:
			return fmt.Sprintf("Wrong password")
		}
	},
	MsgEnglish: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Английский")
		default:
			return fmt.Sprintf("English")
		}
	},
	MsgRussian: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Русский")
		default:
			return fmt.Sprintf("Russian")
		}
	},
	MsgEmailExists: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Такая почта уже сущствует")
		default:
			return fmt.Sprintf("This email already exists")
		}
	},
	MsgUnhide: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Показать")
		default:
			return fmt.Sprintf("Unhide")
		}
	},
	MsgHide: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Скрыть")
		default:
			return fmt.Sprintf("Hide")
		}
	},
	MsgPersons: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Люди")
		default:
			return fmt.Sprintf("Persons")
		}
	},
	MsgUsername: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Никнейм")
		default:
			return fmt.Sprintf("Username")
		}
	},
	MsgAdd: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Добавить")
		default:
			return fmt.Sprintf("Add")
		}
	},
	MsgErrorUsernameEmpty: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Введите имя")
		default:
			return fmt.Sprintf("Enter a name")
		}
	},
	MsgErrorUsernameAlreadyExists: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Данное имя уже занято")
		default:
			return fmt.Sprintf("This name is already taken")
		}
	},
	MsgErrorProductTitle: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Название должно содержать больше 4 и меньше 128 символов")
		default:
			return fmt.Sprintf("Title should contain more than 4 and less than 128 characters")
		}
	},
	MsgErrorProductCalories: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Количество калорий не может быть больше 1000")
		default:
			return fmt.Sprintf("Amount of calories can't be more than 1000")
		}
	},
	MsgErrorProductNutrient: func(locale Locale, args []string) string {
		switch locale {
		case LocaleRuRU:
			return fmt.Sprintf("Значение не может быть больше 100г")
		default:
			return fmt.Sprintf("Value can't be more than 100g")
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

func LocaleFromString(localeStr string) Locale {
	switch localeStr {
	case "ru-RU":
		return LocaleRuRU
	default:
		return LocaleEnUS
	}
}
