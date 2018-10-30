// Copyright 2018 Thales UK Limited
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions
// of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// ParseConfig contains user-supplied configuration that affects the parsing process.
type ParseConfig struct {
	VoidType []VoidType
}

// A VoidType defines a mapping between a void pointer parameter and its underlying type.
type VoidType struct {
	Function  string
	Parameter string
	Type      string
	Struct    string
}

func (t VoidType) valid() bool {
	return t.Function != "" &&
		t.Parameter != "" &&
		((t.Type != "" && t.Struct == "") || (t.Type == "" && t.Struct != ""))
}

// FromString parses TOML configuration data and returns a new config.
func FromString(config string) (cfg ParseConfig, err error) {
	_, err = toml.Decode(config, &cfg)

	if err != nil {
		err = errors.Wrap(err, "failed to parse config file")
		return
	}

	// Check it's valid
	for _, v := range cfg.VoidType {
		if !v.valid() {
			err = errors.New("invalid void type mapping in config file")
			return
		}
	}
	return
}

// FromFile parses TOML configuration file and returns a new config.
func FromFile(configFile string) (ParseConfig, error) {

	// #nosec G304 We don't care which file the user wants to try and read; if they have the permission, it will work.
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return ParseConfig{}, err
	}

	return FromString(string(configBytes))
}
