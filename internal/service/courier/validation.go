package courier

import (
	"regexp"
	"service-order-avito/internal/domain"
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
	if status == domain.StatusAvailable || status == domain.StatusBusy || status == domain.StatusPaused {
		return true
	}
	return false
}

func IsValidTransportType(transportType string) bool {
	if transportType == domain.TransportTypeFoot || transportType == domain.TransportTypeScooter || transportType == domain.TransportTypeCar {
		return true
	}
	return false
}
