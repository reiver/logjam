package metadata_test

import (
	"testing"

	"reflect"

	"github.com/reiver/go-opt"

	"github.com/reiver/logjam/lib/metadata"
)

func TestMetaData_UnmarshalJSON(t *testing.T) {

	tests := []struct{
		Bytes []byte
		Expected metadata.MetaData
	}{
		{
			Bytes: []byte(`null`),
			Expected: metadata.MetaData{},
		},



		{
			Bytes: []byte(`{}`),
			Expected: metadata.MetaData{},
		},



		{
			Bytes: []byte(
				`{`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":null`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
					`}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool{
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
					`}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{`+
						`"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true`+
						`,`+
						`"30511117-4063-4b2f-8e30-99400eb5f719":false`+
						`,`+
						`"595661ff-862b-42e6-8989-9e48158e6cf5":false`+
					`}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
				},
			},
		},
		{
			Bytes: []byte(
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
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
					"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false,
				},
			},
		},
		{
			Bytes: []byte(
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
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"39fa5b20-6a19-4530-bc02-4d04ffce314c":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
					"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false,
				},
			},
		},



		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordersList":null`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string{"1"},
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string{"1","2"},
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string{"1","2","3"},
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string{"1","2","3","4"},
			},
		},
		{
			Bytes: []byte(
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
			Expected: metadata.MetaData{
				RecordersList: []string{"1","2","3","4","5"},
			},
		},



		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
				`}`,
			),
			Expected: metadata.MetaData{
				Styles: opt.Nothing[string](),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":null`+
				`}`,
			),
			Expected: metadata.MetaData{
				Styles: opt.Nothing[string](),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":""`+
				`}`,
			),
			Expected: metadata.MetaData{
				Styles: opt.Something(""),
			},
		},
		{
			Bytes: []byte(
				`{`+
					`"muted":{}`+
					`,`+
					`"recordingStarted":false`+
					`,`+
					`"styles":".something {\n\tfont-weight:bold;\n}\n"`+
				`}`,
			),
			Expected: metadata.MetaData{
				Styles: opt.Something(".something {\n\tfont-weight:bold;\n}\n"),
			},
		},
	}

	for testNumber, test := range tests {

		var actual metadata.MetaData

		err := actual.UnmarshalJSON(test.Bytes)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one.", testNumber)
			t.Logf("ERROR: %s", err)
			t.Logf("BYTES:\n%s", test.Bytes)
			continue
		}

		expected := test.Expected

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("For test #%d, the actual unmarshaled-json is not what was expected.", testNumber)
			t.Logf("EXPECTED:\n%#v", expected)
			t.Logf("ACTUAL:\n%#v", actual)
			t.Logf("BYTES:\n%s", test.Bytes)
			continue
		}
	}
}

func TestMetaData_UnmarshalJSON_fail(t *testing.T) {

	tests := []struct{
		Bytes []byte
		ExpectedError string
	}{
		{
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): unexpected end of JSON input`,
		},



		{
			Bytes: []byte(`false`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal bool into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`true`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal bool into Go value of type map[string]interface {}`,
		},



		{
			Bytes: []byte(`-1`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal number into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`0`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal number into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`1`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal number into Go value of type map[string]interface {}`,
		},



		{
			Bytes: []byte(`""`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal string into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`"once"`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal string into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`"twice"`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal string into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`"thrice"`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal string into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`"fource"`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal string into Go value of type map[string]interface {}`,
		},



		{
			Bytes: []byte(`[]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[false]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[true]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[false,true]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[-1]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[0]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[1]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`[-1,0,1]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
		{
			Bytes: []byte(`["once","twice","thrice", "fource"]`),
			ExpectedError: `problem json-unmarshaling into a map[string]interface {} (which would eventually be used to load a metadata.MetaData): json: cannot unmarshal array into Go value of type map[string]interface {}`,
		},
	}

	for testNumber, test := range tests {

		var metaData metadata.MetaData

		err := metaData.UnmarshalJSON(test.Bytes)
		if nil == err {
			t.Errorf("For test #%d, expected an error, but did not actually get one.", testNumber)
			t.Logf("BYTES:\n%s", test.Bytes)
			continue
		}

		actual := err.Error()

		expected := test.ExpectedError

		if expected != actual {
			t.Errorf("For test #%d, the actual error is not what was expected.", testNumber)
			t.Logf("EXPECTED: %s", expected)
			t.Logf("ACTUAL:   %s", actual)
			t.Logf("BYTES:\n%s", test.Bytes)
			continue
		}
	}
}
