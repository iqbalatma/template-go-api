package errors

import (
	"template-go-api/enums"
	"template-go-api/utils"
)

func newError(code enums.ResponseCode, defaultMessage string, messages ...string) *utils.HTTPError {
	message := defaultMessage
	if len(messages) > 0 && messages[0] != "" {
		message = messages[0]
	}

	return utils.NewHttpError(message, code, nil)
}

func InvalidCredential(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_AUTHENTICATION, "Invalid credential", messages...)
}

func InvalidAction(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_INVALID_ACTION, "Invalid action", messages...)
}

func QueryParameterInvalid(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_INVALID_QUERY_PARAMETER, "Invalid query parameter", messages...)
}

func InvalidTokenTypeException(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_AUTHENTICATION, "Invalid token type", messages...)
}

func InternalServerError(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_INTERNAL_SERVER_ERROR, "Internal server error", messages...)
}

func UnauthorizedException(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_ACTION_UNAUTHORIZED, "Unauthorized", messages...)
}

func DataNotFoundException(messages ...string) *utils.HTTPError {
	return newError(enums.ERR_NOT_FOUND, "Data not found", messages...)
}
