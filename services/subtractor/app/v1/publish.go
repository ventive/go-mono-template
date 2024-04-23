package v1

import (
	"encoding/json"

	"github.com/ventive/go-mono-template/pkg/logger"
	"github.com/ventive/go-mono-template/pkg/nats"
)

func (a *App) publishError(requestHeaders map[string]string, errInput error) {
	if errInput == nil {
		return
	}

	log := logger.New(appID, "App.publishError")

	natsError := nats.Error{
		Message: errInput.Error(),
	}

	err := a.publishData(requestHeaders, a.config.App.Queues.Publish.Errors, natsError, nil)
	if err != nil {
		log.Error("error when publishing error msg", err)
	}
}

func (a *App) publishData(requestHeaders map[string]string, subject string, data interface{}, withErr error) error {
	log := logger.New(appID, "App.publishData")

	if withErr != nil {
		requestHeaders["X-Error"] = withErr.Error()
	}

	var err error
	natsMsg := nats.NewMsgWithHeaders(subject, requestHeaders)
	if natsMsg.Data, err = json.Marshal(data); err != nil {
		log.Error("error when encoding output message data", err)

		return err
	}

	if err = a.nats.PublishMsg(natsMsg); err != nil {
		log.Error("error when publishing output msg", err)

		return err
	}

	return nil
}

func (a *App) subHandlerReturn(log *logger.Logger, err error, msg *nats.Msg, response interface{}) {
	requestHeaders := make(map[string]string)
	for k := range msg.Header {
		requestHeaders[k] = msg.Header.Get(k)
	}

	replySubject := a.config.App.Queues.Publish.Default
	if msg.Reply != "" {
		replySubject = msg.Reply
	}

	if err = a.publishData(requestHeaders, replySubject, response, err); err != nil {
		log.Error("error when publishing output msg", err)
	}

	a.publishError(requestHeaders, err)
}
