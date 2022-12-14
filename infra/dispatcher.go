package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/types"
)

func Dispatch(data types.RequestResult) error {
	logCtx := logrus.WithFields(logrus.Fields{
		"callbackUrl": data.CallbackURL,
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

	logCtx.Info("sending results")

	return nil
}
