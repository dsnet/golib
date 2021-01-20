// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsoncs

import (
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type containsError string

func (e containsError) Error() string {
	return string(e)
}
func (e containsError) Is(e2 error) bool {
	if e2 == nil {
		return false
	}
	return strings.Contains(e2.Error(), string(e))
}

func TestFormat(t *testing.T) {
	tests := []struct {
		inJSON   string
		wantJSON string
		wantErr  error
	}{{
		inJSON:  "",
		wantErr: io.ErrUnexpectedEOF,
	}, {
		inJSON:  "?",
		wantErr: containsError("expected next JSON value"),
	}, {
		inJSON:   "null",
		wantJSON: "null",
	}, {
		inJSON:   " \n\r\tnull \n\r\t",
		wantJSON: "null",
	}, {
		inJSON:  "nullnull",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:   "true",
		wantJSON: "true",
	}, {
		inJSON:   " \n\r\ttrue \n\r\t",
		wantJSON: "true",
	}, {
		inJSON:  "truetrue",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:   "false",
		wantJSON: "false",
	}, {
		inJSON:   " \n\r\tfalse \n\r\t",
		wantJSON: "false",
	}, {
		inJSON:  "falsefalse",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:   "0",
		wantJSON: "0",
	}, {
		inJSON:   "-0",
		wantJSON: "0",
	}, {
		inJSON:   "1",
		wantJSON: "1",
	}, {
		inJSON:   "-1",
		wantJSON: "-1",
	}, {
		inJSON:   "0.1",
		wantJSON: "0.1",
	}, {
		inJSON:   "-0.1",
		wantJSON: "-0.1",
	}, {
		inJSON:   "1234",
		wantJSON: "1234",
	}, {
		inJSON:   "-1234",
		wantJSON: "-1234",
	}, {
		inJSON:   "12.34",
		wantJSON: "12.34",
	}, {
		inJSON:   "-12.34",
		wantJSON: "-12.34",
	}, {
		inJSON:   "12E0",
		wantJSON: "12",
	}, {
		inJSON:   "12E1",
		wantJSON: "120",
	}, {
		inJSON:   "12e34",
		wantJSON: "1.2e+35",
	}, {
		inJSON:   "12E-0",
		wantJSON: "12",
	}, {
		inJSON:   "12e+1",
		wantJSON: "120",
	}, {
		inJSON:   "12e-34",
		wantJSON: "1.2e-33",
	}, {
		inJSON:   "-12E0",
		wantJSON: "-12",
	}, {
		inJSON:   "-12E1",
		wantJSON: "-120",
	}, {
		inJSON:   "-12e34",
		wantJSON: "-1.2e+35",
	}, {
		inJSON:   "-12E-0",
		wantJSON: "-12",
	}, {
		inJSON:   "-12e+1",
		wantJSON: "-120",
	}, {
		inJSON:   "-12e-34",
		wantJSON: "-1.2e-33",
	}, {
		inJSON:   "1.2E0",
		wantJSON: "1.2",
	}, {
		inJSON:   "1.2E1",
		wantJSON: "12",
	}, {
		inJSON:   "1.2e34",
		wantJSON: "1.2e+34",
	}, {
		inJSON:   "1.2E-0",
		wantJSON: "1.2",
	}, {
		inJSON:   "1.2e+1",
		wantJSON: "12",
	}, {
		inJSON:   "1.2e-34",
		wantJSON: "1.2e-34",
	}, {
		inJSON:   "-1.2E0",
		wantJSON: "-1.2",
	}, {
		inJSON:   "-1.2E1",
		wantJSON: "-12",
	}, {
		inJSON:   "-1.2e34",
		wantJSON: "-1.2e+34",
	}, {
		inJSON:   "-1.2E-0",
		wantJSON: "-1.2",
	}, {
		inJSON:   "-1.2e+1",
		wantJSON: "-12",
	}, {
		inJSON:   "-1.2e-34",
		wantJSON: "-1.2e-34",
	}, {
		inJSON:   "0E0",
		wantJSON: "0",
	}, {
		inJSON:   "0E1",
		wantJSON: "0",
	}, {
		inJSON:   "0e34",
		wantJSON: "0",
	}, {
		inJSON:   "0E-0",
		wantJSON: "0",
	}, {
		inJSON:   "0e+1",
		wantJSON: "0",
	}, {
		inJSON:   "0e-34",
		wantJSON: "0",
	}, {
		inJSON:   "-0E0",
		wantJSON: "0",
	}, {
		inJSON:   "-0E1",
		wantJSON: "0",
	}, {
		inJSON:   "-0e34",
		wantJSON: "0",
	}, {
		inJSON:   "-0E-0",
		wantJSON: "0",
	}, {
		inJSON:   "-0e+1",
		wantJSON: "0",
	}, {
		inJSON:   "-0e-34",
		wantJSON: "0",
	}, {
		inJSON:   "12345678901234567890",
		wantJSON: "12345678901234567000",
	}, {
		inJSON:   "-123456789.0123456789e+0123",
		wantJSON: "-1.2345678901234568e+131",
	}, {
		inJSON:   " \n\r\t-123456789.0123456789e+0123 \n\r\t",
		wantJSON: "-1.2345678901234568e+131",
	}, {
		inJSON:  "- 123456789.0123456789e+0123",
		wantErr: containsError("expected digit character in JSON number"),
	}, {
		inJSON:  "-123456789 .0123456789e+0123",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:  "-123456789. 0123456789e+0123",
		wantErr: containsError("expected digit character in JSON number"),
	}, {
		inJSON:  "-123456789.0123456789 e+0123",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:  "-123456789.0123456789e +0123",
		wantErr: containsError("expected digit character in JSON number"),
	}, {
		inJSON:  "-123456789.0123456789e+ 0123",
		wantErr: containsError("expected digit character in JSON number"),
	}, {
		inJSON:  "+0",
		wantErr: containsError("expected next JSON value"),
	}, {
		inJSON:  "00",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:  "01",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:  "0f",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:  "0.e0",
		wantErr: containsError("expected digit character in JSON number"),
	}, {
		inJSON:   "0.0E0",
		wantJSON: "0",
	}, {
		inJSON:   "0.0E+0",
		wantJSON: "0",
	}, {
		inJSON:   "0.0E-0",
		wantJSON: "0",
	}, {
		inJSON:   "0.0E+0",
		wantJSON: "0",
	}, {
		inJSON:   "0.0E-0",
		wantJSON: "0",
	}, {
		inJSON:  "0.0E0.",
		wantErr: containsError("invalid data"),
	}, {
		inJSON:   `""`,
		wantJSON: `""`,
	}, {
		inJSON:   " \n\r\t\"\" \n\r\t",
		wantJSON: `""`,
	}, {
		inJSON:   " \n\r\t\"hello\" \n\r\t",
		wantJSON: `"hello"`,
	}, {
		inJSON:  "\"\x00\"",
		wantErr: containsError("invalid character"),
	}, {
		inJSON:  "\"\xff\"",
		wantErr: containsError("invalid UTF-8"),
	}, {
		inJSON:   `"` + string(utf8.RuneError) + `"`,
		wantJSON: `"` + string(utf8.RuneError) + `"`,
	}, {
		inJSON:   `"\uFFFD"`,
		wantJSON: "\"\uFFFD\"",
	}, {
		inJSON:  `"\x"`,
		wantErr: containsError("invalid escape code"),
	}, {
		inJSON:  `"\uXXXX"`,
		wantErr: containsError("invalid escape code"),
	}, {
		inJSON:  `"\uDEAD"`, // unmatched surrogate pair
		wantErr: io.ErrUnexpectedEOF,
	}, {
		inJSON:  `"\uDEAD______"`, // unmatched surrogate pair
		wantErr: containsError("invalid escape code"),
	}, {
		inJSON:  `"\uDEAD\uBEEF"`, // invalid surrogate half
		wantErr: containsError("invalid escape code"),
	}, {
		inJSON:   `"\uD800\udead"`, // valid surrogate pair
		wantJSON: `"𐊭"`,
	}, {
		inJSON:   `"\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008\u0009\u000a\u000b\u000c\u000d\u000e\u000f"`,
		wantJSON: `"\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\b\t\n\u000b\f\r\u000e\u000f"`,
	}, {
		inJSON:   `"\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f"`,
		wantJSON: `"\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001a\u001b\u001c\u001d\u001e\u001f"`,
	}, {
		inJSON:   `"\u0020\u0021\u0022\u0023\u0024\u0025\u0026\u0027\u0028\u0029\u002a\u002b\u002c\u002d\u002e\u002f"`,
		wantJSON: `" !\"#$%&'()*+,-./"`,
	}, {
		inJSON:   `"\u0030\u0031\u0032\u0033\u0034\u0035\u0036\u0037\u0038\u0039\u003a\u003b\u003c\u003d\u003e\u003f"`,
		wantJSON: `"0123456789:;<=>?"`,
	}, {
		inJSON:   `"\u0040\u0041\u0042\u0043\u0044\u0045\u0046\u0047\u0048\u0049\u004a\u004b\u004c\u004d\u004e\u004f"`,
		wantJSON: `"@ABCDEFGHIJKLMNO"`,
	}, {
		inJSON:   `"\u0050\u0051\u0052\u0053\u0054\u0055\u0056\u0057\u0058\u0059\u005a\u005b\u005c\u005d\u005e\u005f"`,
		wantJSON: `"PQRSTUVWXYZ[\\]^_"`,
	}, {
		inJSON:   `"\u0060\u0061\u0062\u0063\u0064\u0065\u0066\u0067\u0068\u0069\u006a\u006b\u006c\u006d\u006e\u006f"`,
		wantJSON: "\"`abcdefghijklmno\"",
	}, {
		inJSON:   `"\u0070\u0071\u0072\u0073\u0074\u0075\u0076\u0077\u0078\u0079\u007a\u007b\u007c\u007d\u007e\u007f"`,
		wantJSON: "\"pqrstuvwxyz{|}~\u007f\"",
	}, {
		inJSON:   `"\u0080\u0081\u0082\u0083\u0084\u0085\u0086\u0087\u0088\u0089\u008a\u008b\u008c\u008d\u008e\u008f"`,
		wantJSON: "\"\u0080\u0081\u0082\u0083\u0084\u0085\u0086\u0087\u0088\u0089\u008a\u008b\u008c\u008d\u008e\u008f\"",
	}, {
		inJSON:   `"\u0090\u0091\u0092\u0093\u0094\u0095\u0096\u0097\u0098\u0099\u009a\u009b\u009c\u009d\u009e\u009f"`,
		wantJSON: "\"\u0090\u0091\u0092\u0093\u0094\u0095\u0096\u0097\u0098\u0099\u009a\u009b\u009c\u009d\u009e\u009f\"",
	}, {
		inJSON:   `"\u00a0\u00a1\u00a2\u00a3\u00a4\u00a5\u00a6\u00a7\u00a8\u00a9\u00aa\u00ab\u00ac\u00ad\u00ae\u00af"`,
		wantJSON: "\"\u00a0¡¢£¤¥¦§¨©ª«¬\u00ad®¯\"",
	}, {
		inJSON:   `"\u00b0\u00b1\u00b2\u00b3\u00b4\u00b5\u00b6\u00b7\u00b8\u00b9\u00ba\u00bb\u00bc\u00bd\u00be\u00bf"`,
		wantJSON: `"°±²³´µ¶·¸¹º»¼½¾¿"`,
	}, {
		inJSON:   `"\u00c0\u00c1\u00c2\u00c3\u00c4\u00c5\u00c6\u00c7\u00c8\u00c9\u00ca\u00cb\u00cc\u00cd\u00ce\u00cf"`,
		wantJSON: `"ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏ"`,
	}, {
		inJSON:   `"\u00d0\u00d1\u00d2\u00d3\u00d4\u00d5\u00d6\u00d7\u00d8\u00d9\u00da\u00db\u00dc\u00dd\u00de\u00df"`,
		wantJSON: `"ÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞß"`,
	}, {
		inJSON:   `"\u00e0\u00e1\u00e2\u00e3\u00e4\u00e5\u00e6\u00e7\u00e8\u00e9\u00ea\u00eb\u00ec\u00ed\u00ee\u00ef"`,
		wantJSON: `"àáâãäåæçèéêëìíîï"`,
	}, {
		inJSON:   `"\u00f0\u00f1\u00f2\u00f3\u00f4\u00f5\u00f6\u00f7\u00f8\u00f9\u00fa\u00fb\u00fc\u00fd\u00fe\u00ff"`,
		wantJSON: `"ðñòóôõö÷øùúûüýþÿ"`,
	}, {
		inJSON:   `"בְּרֵאשִׁ֖ית בָּרָ֣א אֱלֹהִ֑ים אֵ֥ת הַשָּׁמַ֖יִם וְאֵ֥ת הָאָֽרֶץ׃"`,
		wantJSON: `"בְּרֵאשִׁ֖ית בָּרָ֣א אֱלֹהִ֑ים אֵ֥ת הַשָּׁמַ֖יִם וְאֵ֥ת הָאָֽרֶץ׃"`,
	}, {
		inJSON:   `"ἐν ἀρχῇ ἐποίησεν ὁ θεὸς τὸν οὐρανὸν καὶ τὴν γῆν"`,
		wantJSON: `"ἐν ἀρχῇ ἐποίησεν ὁ θεὸς τὸν οὐρανὸν καὶ τὴν γῆν"`,
	}, {
		inJSON:   `"في البدء خلق الله السموات والارض."`,
		wantJSON: `"في البدء خلق الله السموات والارض."`,
	}, {
		inJSON:   `"起 初 ， 　 神 创 造 天 地 。"`,
		wantJSON: `"起 初 ， 　 神 创 造 天 地 。"`,
	}, {
		inJSON:   `"ในเริ่มแรกนั้นพระเจ้าทรงเนรมิตสร้างฟ้าและแผ่นดินโลก"`,
		wantJSON: `"ในเริ่มแรกนั้นพระเจ้าทรงเนรมิตสร้างฟ้าและแผ่นดินโลก"`,
	}, {
		inJSON:   "[]",
		wantJSON: "[]",
	}, {
		inJSON:   " \n\r\t[ \n\r\t] \n\r\t",
		wantJSON: "[]",
	}, {
		inJSON:   " \n\r\t[ \n\r\t0 \n\r\t] \n\r\t",
		wantJSON: "[0]",
	}, {
		inJSON:   " \n\r\t[ \n\r\t0 \n\r\t, \n\r\t1 \n\r\t] \n\r\t",
		wantJSON: "[0,1]",
	}, {
		inJSON:  "[0{}]",
		wantErr: containsError("expected ',' character in JSON array"),
	}, {
		inJSON:  "[][]",
		wantErr: containsError("invalid data after top-level value"),
	}, {
		inJSON:   "{}",
		wantJSON: "{}",
	}, {
		inJSON:   " \n\r\t{ \n\r\t} \n\r\t",
		wantJSON: "{}",
	}, {
		inJSON:   " \n\r\t{ \n\r\t\"0\" \n\r\t: \n\r\t0 \n\r\t} \n\r\t",
		wantJSON: `{"0":0}`,
	}, {
		inJSON:   " \n\r\t{ \n\r\t\"0\" \n\r\t: \n\r\t0 \n\r\t, \n\r\t\"1\" \n\r\t: \n\r\t1 \n\r\t} \n\r\t",
		wantJSON: `{"0":0,"1":1}`,
	}, {
		inJSON:  `{"k"}`,
		wantErr: containsError("expected ':' character in JSON object"),
	}, {
		inJSON:  `{"k":0,}`,
		wantErr: containsError("expected '\"' character in JSON string"),
	}, {
		inJSON:  `{"k":0[]}`,
		wantErr: containsError("expected ',' character in JSON object"),
	}, {
		inJSON:  `{"k":0,"\u006b":0}`,
		wantErr: containsError("duplicate key \"k\" in JSON object"),
	}, {
		inJSON:  "{}{}",
		wantErr: containsError("invalid data after top-level value"),
	}, {
		inJSON: `[
			56,
			{
				"d": true,
				"10": null,
				"1": [ ]
			}
		]`,
		wantJSON: `[56,{"1":[],"10":null,"d":true}]`,
	}, {
		inJSON: `{
			"peach": "This sorting order",
			"péché": "is wrong according to French",
			"pêche": "but canonicalization MUST",
			"sin":   "ignore locale"
		}`,
		wantJSON: `{"peach":"This sorting order","péché":"is wrong according to French","pêche":"but canonicalization MUST","sin":"ignore locale"}`,
	}, {
		inJSON: `{
			"1": {"f": {"f": "hi","F": 5} ,"\n": 56.0},
			"10": { },
			"": "empty",
			"a": { },
			"111": [ {"e": "yes","E": "no" } ],
			"A": { }
		}`,
		wantJSON: `{"":"empty","1":{"\n":56,"f":{"F":5,"f":"hi"}},"10":{},"111":[{"E":"no","e":"yes"}],"A":{},"a":{}}`,
	}, {
		inJSON: `{
			"Unnormalized Unicode":"A\u030a"
		}`,
		wantJSON: `{"Unnormalized Unicode":"Å"}`,
	}, {
		inJSON: `{
			"numbers": [333333333.33333329, 1E30, 4.50, 2e-3, 0.000000000000000000000000001],
			"string": "\u20ac$\u000F\u000aA'\u0042\u0022\u005c\\\"\/",
			"literals": [null, true, false]
		}`,
		wantJSON: `{"literals":[null,true,false],"numbers":[333333333.3333333,1e+30,4.5,0.002,1e-27],"string":"€$\u000f\nA'B\"\\\\\"/"}`,
	}, {
		inJSON: `{
			"\u20ac": "Euro Sign",
			"\r": "Carriage Return",
			"\u000a": "Newline",
			"1": "One",
			"\u0080": "Control\u007f",
			"\ud83d\ude02": "Smiley",
			"\u00f6": "Latin Small Letter O With Diaeresis",
			"\ufb33": "Hebrew Letter Dalet With Dagesh",
			"</script>": "Browser Challenge"
		}`,
		wantJSON: `{"\n":"Newline","\r":"Carriage Return","1":"One","</script>":"Browser Challenge","` + "\u0080" + `":"Control` + "\u007f" + `","ö":"Latin Small Letter O With Diaeresis","€":"Euro Sign","😂":"Smiley","דּ":"Hebrew Letter Dalet With Dagesh"}`,
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Logf("Input Data: %q", tt.inJSON)
			gotJSON, gotErr := Format([]byte(tt.inJSON))
			if diff := cmp.Diff(string(tt.wantJSON), string(gotJSON)); diff != "" {
				t.Errorf("Format output mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Format output mismatch (-want +got):\n%s", diff)
			}

			gotValid := Valid([]byte(tt.inJSON))
			wantValid := string(tt.inJSON) == string(gotJSON) && gotErr == nil
			if gotValid != wantValid {
				t.Errorf("Valid() = %v, want %v", gotValid, wantValid)
			}
		})
	}
}
