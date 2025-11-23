package models

type ErrorCode string

const (
	TeamExists  ErrorCode = "TEAM_EXISTS"
	PRExists    ErrorCode = "PR_EXISTS"
	PRMerged    ErrorCode = "PR_MERGED"
	NotAssigned ErrorCode = "NOT_ASSIGNED"
	NoCandidate ErrorCode = "NO_CANDIDATE"
	NotFound    ErrorCode = "NOT_FOUND"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *AppError) Error() string {
	return string(e.Code) + ": " + e.Message
}

type ErrorResponse struct {
	Error *AppError `json:"error"`
}

func NewErrorResponse(code ErrorCode, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func ToErrorResponse(err error) *ErrorResponse {
	if appErr, ok := err.(*AppError); ok {
		return &ErrorResponse{Error: appErr}
	}
	return &ErrorResponse{
		Error: &AppError{
			Code:    NotFound,
			Message: err.Error(),
		},
	}
}
