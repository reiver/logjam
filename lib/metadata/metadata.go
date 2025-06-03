package metadata

import (
	"github.com/reiver/go-erorr"
	"github.com/reiver/go-json"
	"github.com/reiver/go-opt"
)

type MetaData struct {
	Muted         map[string]bool
	RecordersList []string
	Styles        opt.Optional[string]
}

var _ json.Marshaler = MetaData{}

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
