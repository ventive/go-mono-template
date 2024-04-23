package v1

import (
	"encoding/json"

	"github.com/ventive/go-mono-template/internal/types"
	"github.com/ventive/go-mono-template/internal/types/subtractor"
	"github.com/ventive/go-mono-template/pkg/decoder"
	"github.com/ventive/go-mono-template/pkg/logger"
	"github.com/ventive/go-mono-template/pkg/nats"
)

func (a *App) subtractHandler(msg *nats.Msg) {
	log := logger.New(appID, "App.subtractHandler")
	log.Info("New Event")
	var event types.InputEvent

	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		log.Error("Could not decode event", err)
	}
	log.DebugWithExtra("Decoded event", map[string]interface{}{
		"Data": event.Data,
	})
	result, err := a.processSubtractEvent(event)

	a.subHandlerReturn(log, err, msg, result)
}

func (a *App) processSubtractEvent(input types.InputEvent) (float64, error) {
	log := logger.New(appID, "App.processSubtractEvent")
	log.DebugWithExtra("Processing event", map[string]interface{}{
		"Data": input.Data,
	})

	var event subtractor.SubtractEvent
	err := decoder.Decode(input.Data, &event)
	if err != nil {
		log.ErrorWithExtra("could not decode event", map[string]interface{}{"event": input}, err)
		return 0, err
	}

	return event.A - event.B, nil
}
