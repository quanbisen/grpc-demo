package handler

import "errors"

func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("User Service--" + err.Error())
		panic(err)
	}
}

func PanicIfTaskError(err error) {
	if err != nil {
		err = errors.New("Task Service--" + err.Error())
		panic(err)
	}
}
