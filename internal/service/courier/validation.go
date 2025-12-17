package courier

import (
	"regexp"
	"service-order-avito/internal/domain/model"
	"unicode"
)

var validName = regexp.MustCompile(`^[A-Za-zА-Яа-яЁё]+$`)

func IsValidName(name string) bool {
	return validName.MatchString(name)
}

func IsValidPhone(phone string) bool {
	const correctPhoneNumberLen = 12
	if len(phone) != correctPhoneNumberLen {
		return false // длина должна быть ровно 12 символов
	}
	if phone[0] != '+' {
		return false // первый символ должен быть '+'
	}
	for _, r := range phone[1:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func IsValidStatus(status string) bool {
	if status == model.StatusAvailable || status == model.StatusBusy || status == model.StatusPaused {
		return true
	}
	return false
}

func IsValidTransportType(transportType string) bool {
	if transportType == model.TransportTypeFoot || transportType == model.TransportTypeScooter || transportType == model.TransportTypeCar {
		return true
	}
	return false
}
