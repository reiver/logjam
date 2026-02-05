package libmetadata_test

import (
	"testing"

	"bytes"

	"github.com/reiver/go-opt"

	"codeberg.org/greatape/logjam/lib/metadata"
)

func TestMetaData_MarshalJSON(t *testing.T) {

	tests := []struct{
		MetaData libmetadata.MetaData
		Expected []byte
	}{
		{
			Expected: []byte(`{"muted":{},"recordersList":[],"recordingStarted":false,"styles":null}`),
		},



		{
			MetaData: libmetadata.MetaData{
				Muted: nil,
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool(nil),
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
			Expected: []byte(
				`{`+
					`"muted":{`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
					`}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
			Expected: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
					`}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
				},
			},
			Expected: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
						`,`+
						`"595661ff-862b-42e6-8989-9e48158e6cf5":false`+
					`}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
					"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false,
				},
			},
			Expected: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
						`,`+
						`"595661ff-862b-42e6-8989-9e48158e6cf5":false`+
						`,`+
						`"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false`+
					`}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"39fa5b20-6a19-4530-bc02-4d04ffce314c":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
					"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false,
				},
			},
			Expected: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
						`,`+
						`"39fa5b20-6a19-4530-bc02-4d04ffce314c":true`+
						`,`+
						`"595661ff-862b-42e6-8989-9e48158e6cf5":false`+
						`,`+
						`"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false`+
					`}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},



		{
			MetaData: libmetadata.MetaData{
				RecordersList: nil,
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string(nil),
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{"1"},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":["1"]`+
					`,`+
					`"recordingStarted":true`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{"1","2"},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":["1","2"]`+
					`,`+
					`"recordingStarted":true`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{"1","2","3"},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":["1","2","3"]`+
					`,`+
					`"recordingStarted":true`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{"1","2","3","4"},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":["1","2","3","4"]`+
					`,`+
					`"recordingStarted":true`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				RecordersList: []string{"1","2","3","4","5"},
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":["1","2","3","4","5"]`+
					`,`+
					`"recordingStarted":true`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},



		{
			MetaData: libmetadata.MetaData{
				Styles: opt.Nothing[string](),
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Styles: opt.Something(""),
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":""`+
				`}`,
			),
		},
		{
			MetaData: libmetadata.MetaData{
				Styles: opt.Something(".something {\n\tfont-weight:bold;\n}\n"),
			},
			Expected: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":[]`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":".something {\n\tfont-weight:bold;\n}\n"`+
				`}`,
			),
		},
	}

	for testNumber, test := range tests {

		actual, err := test.MetaData.MarshalJSON()
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one.", testNumber)
			t.Logf("ERROR: %s", err)
			t.Logf("META-DATA: %#v", test.MetaData)
			continue
		}

		expected := test.Expected

		if !bytes.Equal(expected, actual) {
			t.Errorf("For test #%d, the actual marshaled-json is not what was expected.", testNumber)
			t.Logf("EXPECTED: (%d)\n%s", len(expected), expected)
			t.Logf("ACTUAL:   (%d)\n%s", len(actual), actual)
			t.Logf("META-DATA: %#v", test.MetaData)
			continue
		}
	}
}
