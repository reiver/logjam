package msgs_test

import (
	"testing"

	"bytes"

	"github.com/reiver/go-json"

	"github.com/reiver/logjam/lib/msgs"
)

func TestLeave_MarshalJSON(t *testing.T) {

	tests := []struct{
		Object msgs.Leave
		Expected []byte
	}{
		{
			Expected: []byte(`{"type":"Leave"}`),
		},
		{
			Object: msgs.Leave{},
			Expected: []byte(`{"type":"Leave"}`),
		},



		{
			Object: msgs.Leave{
				Actor:"apple",
			},
			Expected: []byte(`{"type":"Leave","actor":"apple"}`),
		},
		{
			Object: msgs.Leave{
				Actor:"BANANA",
			},
			Expected: []byte(`{"type":"Leave","actor":"BANANA"}`),
		},
		{
			Object: msgs.Leave{
				Actor:"Cherry",
			},
			Expected: []byte(`{"type":"Leave","actor":"Cherry"}`),
		},
		{
			Object: msgs.Leave{
				Actor:"dAtE",
			},
			Expected: []byte(`{"type":"Leave","actor":"dAtE"}`),
		},



		{
			Object: msgs.Leave{
				Actor:"acct:reiver@mastodon.social",
			},
			Expected: []byte(`{"type":"Leave","actor":"acct:reiver@mastodon.social"}`),
		},
		{
			Object: msgs.Leave{
				Actor:"http://mastodon.social/@reiver",
			},
			Expected: []byte(`{"type":"Leave","actor":"http://mastodon.social/@reiver"}`),
		},
		{
			Object: msgs.Leave{
				Actor:"http://mastodon.social/users/reiver",
			},
			Expected: []byte(`{"type":"Leave","actor":"http://mastodon.social/users/reiver"}`),
		},
	}

	for testNumber, test := range tests {

		actual, err := json.Marshal(test.Object)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
			t.Logf("ERROR: %s", err)
			t.Logf("OBJECT: (%T) %#v", test.Object, test.Object)
			continue
		}

		expected := test.Expected

		if !bytes.Equal(expected, actual) {
			t.Errorf("For test #%d, the actual 'object' (of type %T) is not what was expected.", testNumber, test.Object)
			t.Logf("EXPECTED: (%d)\n%s", len(expected), expected)
			t.Logf("ACTUAL:   (%d)\n%s", len(actual), actual)
			t.Logf("OBJECT: (%T) %#v", test.Object, test.Object)
			continue
		}
	}
}
