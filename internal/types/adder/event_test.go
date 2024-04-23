package adder

import (
	"testing"

	"github.com/ventive/go-mono-template/pkg/decoder"
)

func TestDecode(t *testing.T) {
	type args struct {
		data  map[string]interface{}
		event interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"decode add event", args{data: map[string]interface{}{
			"a": int(10),
			"b": float32(11.2),
		}, event: &AddEvent{},
		}, false},
		{"decode add event should fail", args{data: map[string]interface{}{
			"a": "test",
			"b": int(11),
		}, event: &AddEvent{},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := decoder.Decode(tt.args.data, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("Decode returns error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
