// Copyright 2017, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package jsonfmt

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		in   string
		out  string
		err  error
		opts []Option
	}{{
		in:  "",
		out: "",
		err: jsonError{line: 1, column: 1, message: `unable to parse value: unexpected EOF`},
	}, {
		in:  `["]`,
		out: `["]`,
		err: jsonError{line: 1, column: 2, message: `unable to parse string: "\"]"`},
	}, {
		in:  "[\n\n\n]",
		out: `[]`,
	}, {
		in:  "{\n\n\n}",
		out: `{}`,
	}, {
		in: `{"firstName":"John","lastName":"Smith","isAlive":true,"age":25,"address":{"streetAddress":"21 2nd Street","city":"New York","state":"NY","postalCode":"10021-3100"},"phoneNumbers":[{"type":"home","number":"212 555-1234"},{"type":"office","number":"646 555-4567"},{"type":"mobile","number":"123 456-7890"}],"children":[],"spouse":null}`,
		out: `
			{
				"firstName": "John",
				"lastName": "Smith",
				"isAlive": true,
				"age": 25,
				"address": {
					"streetAddress": "21 2nd Street",
					"city": "New York",
					"state": "NY",
					"postalCode": "10021-3100"
				},
				"phoneNumbers": [
					{"type": "home", "number": "212 555-1234"},
					{"type": "office", "number": "646 555-4567"},
					{"type": "mobile", "number": "123 456-7890"}
				],
				"children": [],
				"spouse": null
			}`,
	}, {
		in: `[[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]]`,
		out: `
			[
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]
			]`,
	}, {
		in: `[[{"0123456789012345678901234567890123456789": "0123456789"}, {"0123456789": "0123456789012345678901234567890123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}], [{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]]`,
		out: `
			[
				[
					{"0123456789012345678901234567890123456789": "0123456789"},
					{"0123456789": "0123456789012345678901234567890123456789"}
				],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}],
				[{"0123456789": "0123456789"}, {"0123456789": "0123456789"}]
			]`,
	}, {
		in: `
			{
				"Management": {
					"ServeAddress": "localhost:8080", "PasswordSalt": "", "PasswordHash": "",
					"SMTP": {"RelayServer": "mail.example.com:587", "Password":"abcdefghijklmnopqrstuvwxyz", "From":"noreply@example.com", "To":"noreply@example.com"},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles":       ["key.priv"], // SSH key file
					"KnownHostFiles": [], // SSH known hosts file
				},

				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":  "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
		out: `
			{
				"Management": {
					"ServeAddress": "localhost:8080",
					"PasswordSalt": "",
					"PasswordHash": "",
					"SMTP": {
						"RelayServer": "mail.example.com:587",
						"Password": "abcdefghijklmnopqrstuvwxyz",
						"From": "noreply@example.com",
						"To": "noreply@example.com",
					},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles": ["key.priv"], // SSH key file
					"KnownHostFiles": [], // SSH known hosts file
				},

				"RateLimit": "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags": ["-w"],
				"RecvFlags": ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source": "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source": "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source": "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
	}, {
		in: `
			{
				"Management": {
					"ServeAddress": "localhost:8080", "PasswordSalt": "", "PasswordHash": "",
					"SMTP": {"RelayServer": "mail.example.com:587", "Password":"abcdefghijklmnopqrstuvwxyz", "From":"noreply@example.com", "To":"noreply@example.com"},
				},

				// SSH comment.
				"SSH": {
					"KeyFiles":       ["key.priv"], // SSH key file
					"KnownHostFiles": [], // SSH known hosts file
				},

				"RateLimit":    "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags":    ["-w"],
				"RecvFlags":    ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source":  "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"],
				}, {
					"Source":  "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"],
				}, {
					// Seperate dataset so it has its own readonly setting
					"Source":  "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"],
				}],
			}`,
		out:  `{"Management":{"ServeAddress":"localhost:8080","PasswordSalt":"","PasswordHash":"","SMTP":{"RelayServer":"mail.example.com:587","Password":"abcdefghijklmnopqrstuvwxyz","From":"noreply@example.com","To":"noreply@example.com"}},"SSH":{"KeyFiles":["key.priv"],"KnownHostFiles":[]},"RateLimit":"10Mi","AutoSnapshot":{"Cron":"* * * * *","Count":3,"TimeZone":"Local"},"SendFlags":["-w"],"RecvFlags":["-s"],"Datasets":[{"AutoSnapshot":{"Cron":"0 6 * * *","TimeZone":"Local","Count":30},"Source":"//example.com/tank/fizz","Mirrors":["//foo.example.com/tank/replicas/fizz-drive"]},{"Source":"//example.com/tank/buzz","Mirrors":["//foo.example.com/tank/replicas/buzz-drive"]},{"Source":"//example.com/tank/users","Mirrors":["//foo.example.com/tank/replicas/users"]}]}`,
		opts: []Option{Minify()},
	}, {
		in: `{"Management":{"ServeAddress":"localhost:8080","PasswordSalt":"","PasswordHash":"","SMTP":{"RelayServer":"mail.example.com:587","Password":"abcdefghijklmnopqrstuvwxyz","From":"noreply@example.com","To":"noreply@example.com"}},"SSH":{"KeyFiles":["key.priv"],"KnownHostFiles":[]},"RateLimit":"10Mi","AutoSnapshot":{"Cron":"* * * * *","Count":3,"TimeZone":"Local"},"SendFlags":["-w"],"RecvFlags":["-s"],"Datasets":[{"AutoSnapshot":{"Cron":"0 6 * * *","TimeZone":"Local","Count":30},"Source":"//example.com/tank/fizz","Mirrors":["//foo.example.com/tank/replicas/fizz-drive"]},{"Source":"//example.com/tank/buzz","Mirrors":["//foo.example.com/tank/replicas/buzz-drive"]},{"Source":"//example.com/tank/users","Mirrors":["//foo.example.com/tank/replicas/users"]}]}`,
		out: `
			{
				"Management": {
					"ServeAddress": "localhost:8080",
					"PasswordSalt": "",
					"PasswordHash": "",
					"SMTP": {
						"RelayServer": "mail.example.com:587",
						"Password": "abcdefghijklmnopqrstuvwxyz",
						"From": "noreply@example.com",
						"To": "noreply@example.com"
					}
				},
				"SSH": {"KeyFiles": ["key.priv"], "KnownHostFiles": []},
				"RateLimit": "10Mi",
				"AutoSnapshot": {"Cron": "* * * * *", "Count": 3, "TimeZone": "Local"},
				"SendFlags": ["-w"],
				"RecvFlags": ["-s"],
				"Datasets": [{
					"AutoSnapshot": {"Cron": "0 6 * * *", "TimeZone": "Local", "Count": 30},
					"Source": "//example.com/tank/fizz",
					"Mirrors": ["//foo.example.com/tank/replicas/fizz-drive"]
				}, {
					"Source": "//example.com/tank/buzz",
					"Mirrors": ["//foo.example.com/tank/replicas/buzz-drive"]
				}, {
					"Source": "//example.com/tank/users",
					"Mirrors": ["//foo.example.com/tank/replicas/users"]
				}]
			}`,
	}, {
		in: "[\n123456789,\n123456789,\n123456789,\n]",
		out: `
			[
				123456789,
				123456789,
				123456789,
			]`,
	}, {
		in: "[\n123456789,\n123456789,\n123456789,\n]",
		out: `
			[
				123456789,
				123456789,
				123456789
			]`,
		opts: []Option{Standardize()},
	}, {
		in:   "[\n123456789,\n123456789,\n123456789,\n]",
		out:  "[123456789,123456789,123456789]",
		opts: []Option{Minify()},
	}, {
		in: `

					/*
					* Block comment.
					*/
					"Text"
		`,
		out: `
			/*
			 * Block comment.
			 */
			"Text"`,
	}, {
		in: `
				[
									{
										"fwafwa" /*ffawe*/:
								    		"fewafwaf",

						"fwafwafwae":




					                 		"fwafewa",},

				[/*comment*/
				{/*comment*/},
					{

					}




					],

						{"fwafwa":



								    		"fewafwaf",
						"fwafwafwae": "dwafewa",//fea
						"fwafwafwae"://fa
						"fwafewa",},

						{
							"fwafwa": 0.0000000000000000000033242000000,
						"fwafwafwae"


						:				"fwafewa",
						 },
						 ["fweafewa","faewfaew","afwefawe"/*
						 fewfaew
						 fewafewa*/]
				 				    ]`,
		out: `
			[
				{
					"fwafwa" /*ffawe*/ :
						"fewafwaf",

					"fwafwafwae":
						"fwafewa",
				},

				[ /*comment*/
					{ /*comment*/ },
					{},
				],

				{
					"fwafwa":
						"fewafwaf",
					"fwafwafwae": "dwafewa", //fea
					"fwafwafwae": //fa
						"fwafewa",
				},

				{
					"fwafwa": 3.3242e-21,
					"fwafwafwae": "fwafewa",
				},
				["fweafewa", "faewfaew", "afwefawe", /*
				fewfaew
				fewafewa*/ ],
			]`,
	}}

	for i, tt := range tests {
		// Adjust output for leading tabs and newlines.
		out := strings.Join(strings.Split(tt.out, "\n\t\t\t"), "\n")
		if strings.HasPrefix(out, "\n") {
			out = out[1:] + "\n"
		}

		got, err := Format([]byte(tt.in), tt.opts...)
		if string(got) != out || err != tt.err {
			t.Errorf("test %d, Format(`%s`):\ngot  (`%s`, %v)\nwant (`%s`, %v)", i, tt.in, got, err, out, tt.err)
		}
	}
}
