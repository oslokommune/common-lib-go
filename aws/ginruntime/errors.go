package ginruntime

type UnauthorizedError struct{}

func (unauthorized UnauthorizedError) Error() string {
	return "Unauthorized"
}

type ApiError struct {
	Details error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (appError ApiError) Error() string {
	return appError.Message
}

type DbError struct {
	Err  error `json:"message"`
	Code int   `json:"code"`
}

func (dbError DbError) Error() string {
	return dbError.Err.Error()
}
