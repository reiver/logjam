package libmetadata_test

import (
	"testing"

	"reflect"

	"github.com/reiver/go-opt"

	"github.com/reiver/logjam/lib/metadata"
	"github.com/reiver/logjam/lib/reply"
)

func TestMetaData_CloneFrom(t *testing.T) {

	tests := []struct{
		Original libmetadata.MetaData
	}{
		{
			Original: libmetadata.MetaData{},
		},



		{
			Original: libmetadata.MetaData{},
		},



		{
			Original: libmetadata.MetaData{
				Messages: []libreply.UserMessageModel(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				Messages: []libreply.UserMessageModel{libreply.UserMessageModel{}},
			},
		},
		{
			Original: libmetadata.MetaData{
				Messages: []libreply.UserMessageModel{libreply.UserMessageModel{Message:"Hello world!", SenderId:2}},
			},
		},



		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool{
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
				},
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
				},
			},
		},
		{
			Original: libmetadata.MetaData{
				Muted: map[string]bool{
					"0ed1d891-629d-44d8-9274-b8b9d1d4e6f3":true,
					"30511117-4063-4b2f-8e30-99400eb5f719":false,
					"595661ff-862b-42e6-8989-9e48158e6cf5":false,
					"8a94aaec-db1a-48b7-a892-baf2cd17c51b":false,
				},
			},
		},
		{
			Original: libmetadata.MetaData{
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
			Original: libmetadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string(nil),
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string{"1"},
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string{"1","2"},
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string{"1","2","3"},
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string{"1","2","3","4"},
			},
		},
		{
			Original: libmetadata.MetaData{
				RecordersList: []string{"1","2","3","4","5"},
			},
		},



		{
			Original: libmetadata.MetaData{
				Styles: opt.Nothing[string](),
			},
		},
		{
			Original: libmetadata.MetaData{
				Styles: opt.Nothing[string](),
			},
		},
		{
			Original: libmetadata.MetaData{
				Styles: opt.Something(""),
			},
		},
		{
			Original: libmetadata.MetaData{
				Styles: opt.Something(".something {\n\tfont-weight:bold;\n}\n"),
			},
		},
	}

	for testNumber, test := range tests {

		var actual libmetadata.MetaData

		err := actual.CloneFrom(&(test.Original))
		if nil != err {
			t.Errorf("For test #%d, did not expect an error, but actually got one.", testNumber)
			t.Logf("ERROR: %s", err)
			t.Logf("ORIGINAL:\n%#v", test.Original)
			continue
		}

		expected := &(test.Original)

		if !reflect.DeepEqual(expected, &actual) {
			t.Errorf("For test #%d, the actual clone is not what was expected.", testNumber)
			t.Logf("EXPECTED:\n%#v", *expected)
			t.Logf("ACTUAL:\n%#v", actual)
			t.Logf("ORIGINAL:\n%#v", test.Original)
			continue
		}
	}
}
