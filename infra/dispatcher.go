package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/types"
)

type Dispatcher struct{}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Dispatch(data types.RequestResult) error {
	logCtx := logrus.WithFields(logrus.Fields{
		"callbackURL": data.CallbackURL,
		"requestID":   data.RequestID,
	})

	json_data, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("unable to marshal data into JSON")
	}

	_, error := http.Post(data.CallbackURL, "application/json", bytes.NewBuffer(json_data))

	if error != nil {
		logCtx.Error("dispatcher: unable to dispatch results")
		return error
	}

	return nil
}
