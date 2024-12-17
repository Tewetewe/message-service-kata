package middleware

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	// RouteTypePublic public route forwarded from traefik
	RouteTypePublic = "public"
	// RouteTypePrivate define private forwarded from traefik
	RouteTypePrivate = "private"
	// RouteTypeProtect define protect route forwarded from traefik
	RouteTypeProtect = "protect"
	// RouteTypeStrict define strict route forwarded from traefik
	RouteTypeStrict = "strict"
	// RouteTypeShared define shared route forwarded from traefik
	RouteTypeShared = "shared"
	// RouteTypeExclusive define exclusive forwarded from traefik
	RouteTypeExclusive = "exclusive"

	// RestHeaderKeyRouteType define rest header for x-kata-route-type
	RestHeaderKeyRouteType = "x-kata-route-type"
	// RestHeaderKeyUserID define rest header for x-kata-auth-user-id
	//nolint:gosec
	RestHeaderKeyUserID = "x-kata-auth-user-id"
	// RestHeaderKeyUserEmail define rest header for x-kata-auth-user-email
	RestHeaderKeyUserEmail = "x-kata-auth-user-email"
	// RestHeaderKeyUserCode  define rest header for x-kata-auth-user-code
	RestHeaderKeyUserCode = "x-kata-auth-user-code"
	// RestHeaderKeyUserType  define rest header for x-kata-auth-user-type
	RestHeaderKeyUserType = "x-kata-auth-user-type"
	// RestHeaderKeyPhoneArea define rest header for x-kata-auth-user-phone-area
	RestHeaderKeyPhoneArea = "x-kata-auth-user-phone-area"
	// RestHeaderKeyPhoneNumber define rest header for x-kata-auth-user-phone-number
	RestHeaderKeyPhoneNumber = "x-kata-auth-user-phone-number"
	// RestHeaderKeyUserDivision define rest header for x-kata-auth-user-division
	RestHeaderKeyUserDivision = "x-kata-auth-user-division"
	// RestHeaderKeyVendorID  define rest header for x-kata-auth-vendor-id
	RestHeaderKeyVendorID = "x-kata-auth-vendor-id"
	// RestHeaderKeyVendorUUID define rest header for x-kata-auth-vendor-uuid
	RestHeaderKeyVendorUUID = "x-kata-auth-vendor-uuid"
	// RestHeaderKeyVendorCode define rest header for x-kata-auth-vendor-code
	RestHeaderKeyVendorCode = "x-kata-auth-vendor-code"
	// RestHeaderKeyVendorName  define rest header for x-kata-auth-vendor-name
	RestHeaderKeyVendorName = "x-kata-auth-vendor-name"
	// RestHeaderKeyContextType define rest header for x-kata-auth-context-type
	RestHeaderKeyContextType = "x-kata-auth-context-type"
	// RestHeaderKeyContextKey define rest header for x-kata-auth-context-key
	RestHeaderKeyContextKey = "x-kata-auth-context-key"
)

var headerSpecification = map[string][]struct {
	Key        string
	IsOptional bool
}{
	RouteTypeProtect: {
		{
			Key: RestHeaderKeyRouteType,
		},
		{
			Key: RestHeaderKeyUserID,
		},
		{
			Key: RestHeaderKeyUserEmail,
		},
		{
			Key: RestHeaderKeyUserCode,
		},
		{
			Key: RestHeaderKeyUserType,
		},
	},
	RouteTypeStrict: {
		{
			Key: RestHeaderKeyRouteType,
		},
		{
			Key: RestHeaderKeyUserID,
		},
		{
			Key: RestHeaderKeyUserEmail,
		},
		{
			Key: RestHeaderKeyUserType,
		},
		{
			Key: RestHeaderKeyUserDivision,
		},
	},
	RouteTypeShared: {
		{
			Key: RestHeaderKeyRouteType,
		},
		{
			Key: RestHeaderKeyVendorID,
		},
		{
			Key: RestHeaderKeyVendorUUID,
		},
		{
			Key: RestHeaderKeyVendorCode,
		},
		{
			Key: RestHeaderKeyVendorName,
		},
	},
	RouteTypeExclusive: {
		{
			Key: RestHeaderKeyRouteType,
		},
		{
			Key: RestHeaderKeyVendorID,
		},
		{
			Key: RestHeaderKeyVendorUUID,
		},
		{
			Key: RestHeaderKeyVendorCode,
		},
		{
			Key: RestHeaderKeyVendorName,
		},
		{
			Key: RestHeaderKeyContextType,
		},
		{
			Key: RestHeaderKeyContextKey,
		},
	},
	RouteTypePrivate: {
		{
			Key: RestHeaderKeyRouteType,
		},
		{
			Key:        RestHeaderKeyUserEmail,
			IsOptional: true,
		},
		{
			Key:        RestHeaderKeyUserCode,
			IsOptional: true,
		},
		{
			Key:        RestHeaderKeyPhoneArea,
			IsOptional: true,
		},
		{
			Key:        RestHeaderKeyPhoneNumber,
			IsOptional: true,
		},
	},
}

// ProtectMiddleware to verify Protect token from request
func ProtectMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !verifyRestHeader(RouteTypeProtect, c.Request().Header) {
			err := errors.New("rest header is not valid")
			c.Error(err)
			log.Error().Any("Error", err.Error()).Msg("Verify protect rest header error")
			return nil
		}

		return next(c)
	}
}

// StrictMiddleware to verify strict token from request
func StrictMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !verifyRestHeader(RouteTypeStrict, c.Request().Header) {
			err := errors.New("rest header is not valid")
			c.Error(err)
			log.Error().Any("Error", err.Error()).Msg("Verify strict rest header error")
			return nil
		}

		return next(c)
	}
}

// PublicMiddleware to verify strict token from request
func PublicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !verifyRestHeader(RouteTypePublic, c.Request().Header) {
			err := errors.New("rest header is not valid")
			c.Error(err)
			log.Error().Any("Error", err.Error()).Msg("Verify public rest header error")
			return nil
		}

		return next(c)
	}
}

// verifyRestHeader to verify rest header
//
//nolint:gocyclo
func verifyRestHeader(routeType string, header http.Header) (isValid bool) {
	for _, spec := range headerSpecification[routeType] {
		var values []string
		if values, isValid = header[http.CanonicalHeaderKey(spec.Key)]; (len(values) == 0 || values[0] == "" ||
			!isValid) && !spec.IsOptional {
			isValid = false

			return isValid
		}

		if spec.IsOptional && (len(values) == 0 || !isValid) {
			isValid = true

			continue
		}

		value := values[0]

		switch spec.Key {
		case RestHeaderKeyRouteType:
			if isValid = IsExist(value, []string{
				RouteTypePrivate,
				RouteTypeProtect,
				RouteTypeStrict,
				RouteTypeShared,
				RouteTypeExclusive,
			}); !isValid {
				return isValid
			}

		case RestHeaderKeyUserID, RestHeaderKeyPhoneArea, RestHeaderKeyPhoneNumber, RestHeaderKeyVendorID:
			if isValid = govalidator.IsNumeric(value); !isValid {
				return isValid
			}

		case RestHeaderKeyUserEmail:
			if isValid = govalidator.IsEmail(value); !isValid {
				return isValid
			}

		case RestHeaderKeyUserType, RestHeaderKeyVendorUUID:
			if isValid = govalidator.IsUUID(value); !isValid {
				return isValid
			}
		}
	}

	switch routeType {
	case RouteTypePrivate:
		var (
			_, isEmailExist       = header[http.CanonicalHeaderKey(RestHeaderKeyUserEmail)]
			_, isPhoneAreaExist   = header[http.CanonicalHeaderKey(RestHeaderKeyPhoneArea)]
			_, isPhoneNumberExist = header[http.CanonicalHeaderKey(RestHeaderKeyPhoneArea)]
		)

		isValid = isEmailExist || (isPhoneAreaExist && isPhoneNumberExist)
	case RouteTypePublic:
		isValid = true
	case RouteTypeProtect:
		var (
			_, isUserIDExist   = header[http.CanonicalHeaderKey(RestHeaderKeyUserID)]
			_, isUserCodeExist = header[http.CanonicalHeaderKey(RestHeaderKeyUserCode)]
		)

		isValid = isUserIDExist || isUserCodeExist
	}

	return isValid
}

// IsExist to checking value from interface
func IsExist(value, array interface{}) (exist bool) {
	exist = false
	if reflect.TypeOf(array).Kind() == reflect.Slice {
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				exist = true
				return exist
			}
		}
	}

	return exist
}
