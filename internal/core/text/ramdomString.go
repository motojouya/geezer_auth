package text

import (
	"errors"
	"strconv"
)

func GetStringTimes[T any](getString func() (T, error), checkString func(T) (bool, error), times int) (T, error) {
	var str T

	for i := 0; i < times; i++ {
		var err error
		str, err = getString()
		if err != nil {
			var zero T
			return zero, err
		}

		var isExist, checkErr = checkString(str)
		if checkErr != nil {
			var zero T
			return zero, checkErr
		}

		if !isExist {
			return str, nil
		}
	}

	var zero T
	return zero, NewStringIssueError(times, "Failed to get unique identifier after retries")
}

/*
 * StringIssueError
 */
type StringIssueError struct {
	Times int
	error
}

func NewStringIssueError(times int, message string) *StringIssueError {
	return &StringIssueError{
		Times:  times,
		error: errors.New(message),
	}
}

func (e StringIssueError) Error() string {
	return e.error.Error() + ", Times: " + strconv.Itoa(e.Times)
}

func (e StringIssueError) Unwrap() error {
	return e.error
}

func (e StringIssueError) HttpStatus() uint {
	return 500
}
