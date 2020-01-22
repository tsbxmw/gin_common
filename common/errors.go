package common

type (
    HttpAuthError struct {
        Code    int
        Message string
    }

    MysqlCreateError struct {
        Code    int
        Message string
    }
)

func NewHttpAuthError() error {
    return HttpAuthError{Code: HTTP_AUTH_ERROR, Message: HTTP_AUTH_ERROR_MSG}
}

func (hae HttpAuthError) Error() string {
    return hae.Message
}

func NewMySqlCreateError() error {
    return MysqlCreateError{Code: MYSQL_CREATE_ERROR, Message: MYSQL_CREATE_ERROR_MSG}
}

func (mce MysqlCreateError) Error() string {
    return mce.Message
}