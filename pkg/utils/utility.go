package utils

import (
	"strconv"

	"message-service-kata/pkg/cerror"
)

// ValidateUserID is use for convert user id from string to integer
func ValidateUserID(userID string) (int64, error) {
	if userID == "" {
		return int64(0), cerror.ErrUserIDNotPresentInHeader
	}

	formattedUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return int64(0), cerror.ErrUserIDNotPresentInHeader
	}

	return formattedUserID, nil
}
