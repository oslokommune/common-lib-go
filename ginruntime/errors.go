package ginruntime

type UnauthorizedError struct {
}

func (unauthorized UnauthorizedError) Error() string {
	return "Unauthorized"
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details error  `json:"-"`
}

func (appError ApiError) Error() string {
	return appError.Message
}

type DbError struct {
	Code int   `json:"code"`
	Err  error `json:"message"`
}

func (dbError DbError) Error() string {
	return dbError.Err.Error()
}
