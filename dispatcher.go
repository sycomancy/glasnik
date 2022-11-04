package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/types"
)

type Dispatcher struct{}

func (d *Dispatcher) NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(data types.RequestResult) error {
	json_data, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("unable to marshal data into JSON")
	}

	res, error := http.Post(data.CallbackURL, "application/json", bytes.NewBuffer(json_data))

	logrus.WithFields(logrus.Fields{
		"webhook":   data.CallbackURL,
		"requestID": data.RequestID,
		"status":    res.Status,
	}).Info("sending data to webhook")

	if error != nil {
		return error
	}

	return nil
}
