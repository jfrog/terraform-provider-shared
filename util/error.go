package util

import (
	"fmt"

	"github.com/samber/lo"
)

type JFrogErrors struct {
	Errors []JFrogError `json:"errors"`
}

func (e JFrogErrors) String() string {
	return lo.Reduce(
		e.Errors,
		func(agg string, err JFrogError, _ int) string {
			if agg == "" {
				return err.Message
			}

			return fmt.Sprintf("%s %s.", agg, err.Message)
		},
		"",
	)
}

type JFrogError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
