package metadata

import (
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-json"
	"github.com/reiver/go-opt"
)

const (
	errNilReceiver = erorr.Error("nil receiver")
)

type MetaData struct {
	Muted         map[string]bool
	RecordersList []string
	Styles        opt.Optional[string]
}

var _ json.Marshaler = MetaData{}
var _ json.Unmarshaler = &MetaData{}

func (receiver MetaData) MarshalJSON() ([]byte, error) {
	var buffer [512]byte
	var p []byte = buffer[0:0]

	p = append(p, '{')

	{
		const prefix string = `"muted":`

		p = append(p, prefix...)

		marshaled, err := json.Marshal(receiver.Muted)
		if nil != err {
			return nil, erorr.Errorf("problem json-marshaling MetaData.Muted: %w", err)
		}
		p = append(p, marshaled...)
	}

	if 0 < len(receiver.RecordersList) {
		p = append(p, ',')

		{
			const prefix string = `"recordersList":[`

			p = append(p, prefix...)
			for index, recorder := range receiver.RecordersList {
				if 0 < index {
					p = append(p, ',')
				}
				p = append(p, json.MarshalString(recorder)...)
			}
			p = append(p, ']')
		}
	}

	p = append(p, ',')

	{
		const prefix string = `"recordingStarted":`

		p = append(p, prefix...)

		var recordingStarted bool = (0 < len(receiver.RecordersList))
		p = append(p, json.MarshalBool(recordingStarted)...)
	}

	p = append(p, ',')

	{
		const prefix string = `"styles":`

		p = append(p, prefix...)

		value, found := receiver.Styles.Get()
		switch {
		case found:
			p = append(p, json.MarshalString(value)...)
		default:
			p = append(p, "null"...)
		}
	}

	p = append(p, '}')

	return p, nil
}

func (receiver *MetaData) UnmarshalJSON(bytes []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	var data map[string]any = map[string]any{}

	err := json.Unmarshal(bytes, &data)
	if nil != err {
		return erorr.Errorf("problem json-unmarshaling into a %T (which would eventually be used to load a metadata.MetaData): %w", data, err)
	}

	if mutedAny, found := data["muted"]; found {
		if muted, casted := mutedAny.(map[string]any); casted {

			var values map[string]bool = map[string]bool{}

			for streamID, isMutedAny := range muted {
				if isMuted, isBool := isMutedAny.(bool); isBool {
					values[streamID] = isMuted
				} else {
					return erorr.Errorf("problem unmarshaling metadata — muted.%q is NOT a bool but is instead a %T (%#v)", streamID, isMutedAny, isMutedAny)
				}
			}

			if 0 < len(values) {
				receiver.Muted = values
			}
		}
	}

	if recordersListAny, found := data["recordersList"]; found {
		if recordersList, casted := recordersListAny.([]any); casted {
			var strings []string

			for _, recorder := range recordersList {
				if str, isString := recorder.(string); isString {
					strings = append(strings, str)
				}
			}

			receiver.RecordersList = strings
		}
	}

	if stylesAny, found := data["styles"]; found {
		if styles, casted := stylesAny.(string); casted {
			receiver.Styles = opt.Something(styles)
		}
	}

	return nil
}
