// Code generated by go-bindata.
// sources:
// translate/templates/defs.gotmpl
// translate/templates/main.gotmpl
// DO NOT EDIT!

package translate

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _translateTemplatesDefsGotmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x50\xc1\x4e\xc3\x30\x0c\xbd\xe7\x2b\x7c\x04\x09\x75\xfc\x02\x02\xd6\x1b\x9b\xd8\xc4\x05\x71\xb0\x3a\xb7\xaa\x68\xdc\x2a\x4d\x0e\x95\xe5\x7f\x47\x4d\x02\xcd\xa4\x71\x8a\xed\xbc\xe7\xf7\x9e\x27\x6c\xbe\xb1\x23\xb0\xd8\xb3\x31\x22\x0e\xb9\x23\xa8\xce\xcb\x44\xb3\xaa\xf1\xcb\x44\x20\x52\xd5\xe3\xc9\xbb\xd0\xf8\x37\xb4\xa4\x0a\x73\x6c\x40\x36\x42\x4d\x7e\xdf\xd3\x70\x99\x55\x13\x3e\x23\x45\xfa\x16\x98\xe0\x11\xaa\x27\xe7\x70\x79\x1e\x03\x7b\xd5\x4f\x91\xab\xfe\x4b\x84\xf8\xa2\x5a\xc0\x8f\x63\xcf\x9e\x5c\x06\x88\x44\x4f\x55\x3d\x46\xda\x5a\x27\x05\x11\x1a\x66\x2a\x01\xa5\xd3\xbc\xd6\xfc\xbe\x5b\xb5\x59\xdf\x07\x6e\xd6\xac\xa6\x0d\xdc\x5c\x99\xbf\x2b\xe3\x1d\xd1\xa1\xcd\xf1\x4e\xd8\xd2\x5f\xc0\xe8\xa9\x5b\xc5\x55\x1f\xf2\xfa\xfb\x14\x7c\xe5\xbd\x93\x0f\x8e\xd3\x77\x44\x97\xa3\x5b\x9e\x21\xef\x00\x31\x00\x00\xbb\x1d\x9c\x0f\x2f\x07\xe8\xed\x34\x90\x25\xf6\xa5\xc5\x9b\x2a\x91\xe6\xe2\xe0\x1f\xc1\x57\x3b\xf9\xe5\x03\x87\x50\x9c\x28\x9e\x24\x95\x3f\x01\x00\x00\xff\xff\x75\x30\xb0\x84\x15\x02\x00\x00")

func translateTemplatesDefsGotmplBytes() ([]byte, error) {
	return bindataRead(
		_translateTemplatesDefsGotmpl,
		"translate/templates/defs.gotmpl",
	)
}

func translateTemplatesDefsGotmpl() (*asset, error) {
	bytes, err := translateTemplatesDefsGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "translate/templates/defs.gotmpl", size: 533, mode: os.FileMode(436), modTime: time.Unix(1541497838, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _translateTemplatesMainGotmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x56\xcb\x6e\x22\x39\x14\x5d\xe3\xaf\xb8\xc3\x62\x54\x45\x92\x22\xb3\x25\xc3\x82\x41\x0a\xbb\x28\x0a\xd1\xcc\x22\xca\xc2\x2a\x5c\xc8\x0a\xb8\x18\x97\x49\x3a\x6d\xf9\xdf\x5b\xbe\x7e\xd4\x03\x13\xa4\x2c\xba\x37\x4d\xec\xfb\xf0\x39\xe7\x3e\xea\x40\xcb\x37\xba\x65\xb0\xa7\x5c\x10\x32\x9d\x10\xad\x25\x15\x5b\x06\xc5\xf3\xe7\x81\x35\xc6\x10\xad\x79\x05\x82\xc1\x78\xec\xce\x1e\xe8\x9e\xc1\x8d\x31\xa0\x3e\x0f\x6c\xc3\x2a\xd0\x9a\x89\x0d\x9e\x34\x4a\x1e\x4b\x05\x5a\x17\xd6\xc8\x18\xd0\x6d\xb8\x15\x53\xf7\x9c\xed\x36\x8d\x31\x60\x2d\x6c\xa8\x62\x49\x45\x2d\x78\x49\x77\xc1\x5e\x17\x8f\x35\x17\x8a\xc9\xc6\x98\x18\x26\xbc\xe0\x16\x8a\x85\x94\xf4\x73\x59\x1f\x85\x32\xe6\x45\xeb\xde\xdf\xaf\xf8\x10\x63\xee\x88\xff\x61\x42\x1e\x17\x26\x9e\x93\xc9\x94\xf0\xfd\xa1\x96\x0a\xc6\xcb\x31\xd1\xfa\x06\x78\x05\x21\x73\xc0\x1d\x2c\x8e\xa2\xa1\x15\x1b\x47\xe7\x16\xd2\xc0\xc1\xf2\x61\x33\xae\x6a\x7c\x55\x9b\x38\xf2\x42\x46\xe5\x86\x2a\x0a\x30\xd1\xba\x58\x6e\xeb\x35\x9e\x3b\x23\x32\xda\xd6\x78\xf9\xf2\x8a\x31\x7a\x77\x86\x90\xea\x28\x4a\xc8\x28\xba\x9e\x64\xc8\x61\xbd\xe3\x25\xcb\x76\x4c\x00\x17\x2a\x4f\x05\xb1\xe9\x79\x05\xb4\xf0\x79\xfe\x98\x83\xe0\x3b\x7b\x3a\x92\x4c\x1d\xa5\x88\x57\x64\x64\x08\x19\x45\xc3\x39\xec\xe9\x1b\xcb\x12\x21\xaf\x61\xc7\x44\x4e\x46\x55\x2d\x81\xc3\x6c\x0e\xb7\x77\xc0\xe1\x6f\x7b\x7a\x07\xfc\xea\x0a\x83\x87\x38\x2f\xfc\x15\xe6\x96\x9f\x65\x2d\xde\x99\x6c\x78\x2d\xee\x8f\xa2\x54\xbc\x16\xc6\x64\x93\x2c\xc1\x49\x9e\x39\xf2\x03\xd3\xd9\x91\x0b\x75\x50\x72\x78\x4c\x0b\xe4\x35\xcf\xe1\x0a\x82\x09\xcf\x27\xde\x6a\xcd\x7f\xb2\xba\xca\x26\xd1\x2a\xcf\x1d\xc2\x21\x6c\x43\x2e\xd2\xfc\x21\xb9\x62\xff\xd0\xf2\x2d\xcb\x2d\xb8\x08\xdc\x55\x44\xa4\xcc\xe2\xfe\xad\x88\x1c\xb3\x4f\x49\x6a\x3b\x02\x20\x70\x93\x28\xe5\x50\xc3\x0e\x7f\x89\x61\xd4\x73\x7d\xa2\x78\xd6\x40\x02\x14\xa4\x8a\x0d\x00\xc0\x33\x7c\xee\xda\xfe\x4b\x0d\x08\xb4\x77\x96\x33\xc0\xee\x67\xff\xc3\x5f\xb1\xe1\x7c\xbf\x87\x21\x72\x22\x93\x46\x56\x66\xd0\x14\x5a\x17\x6b\x5a\xf9\x63\xa3\x35\xdb\x35\x38\x4f\xdc\xdc\x58\x7c\x50\xc9\x5a\xca\x2c\xba\xbe\x47\xee\x89\xba\xee\xbc\xd6\x31\x67\x7f\x7f\xc1\xe4\x03\x63\x1b\x2e\xb6\xdd\xd8\x7d\x6a\xef\x65\xbd\x3f\x43\xee\xea\x94\xdb\x01\xdf\x9e\xbd\x34\x73\x71\x54\xf6\xc9\x42\x64\x81\xd4\xa2\x53\xc6\xa4\x03\x2b\x4e\xf2\x81\x78\xe9\xfc\x5f\xa8\xd7\x52\xe8\xf5\xbb\xfc\x24\xd4\xac\xa3\x10\x4a\x9b\x2e\xe8\x9e\xe3\x77\x34\x8a\xf5\x92\x2e\x79\x5f\x55\x5d\xc8\xe8\x31\xd8\x37\x76\x4c\x24\x76\x90\xdf\x6c\xc3\x16\xc9\x24\x6b\xce\x9b\x9f\x88\x8e\xaf\x1f\x8e\x97\x0e\xf1\x92\x35\x71\xa0\xba\x84\xc9\xd6\xc7\x9e\x77\x44\xb4\x92\x26\x39\xb1\x5e\x48\xc7\x74\xca\x7e\xe0\xe6\x8b\xfb\xd7\x31\xd4\x6e\xf5\xac\xab\xfa\x23\x95\x74\xdf\xb8\x5d\xdb\xca\xee\x54\x2f\xfe\xa3\xcd\xbf\x35\xdf\x78\xdd\x8d\xe9\x4f\xb9\x8e\xda\xbd\xad\x9f\x22\x30\xa8\xec\xff\xcf\x7d\x82\x15\x53\x4f\x88\xc9\xfa\xb8\x47\xf4\x8e\x92\xb1\x42\x85\xc4\x3e\xba\x81\x53\x3c\x67\xc7\x0e\xfa\xbc\x53\x89\x1f\x2b\xd6\x78\xd1\x2c\xe4\xd6\x18\x5c\x19\xe9\x81\x84\x2e\xbc\x1a\x50\x14\x17\x70\x5b\xb8\xfd\x80\x73\xf8\xf3\xe2\x8c\x4b\xf3\x9c\x4d\xce\xd5\x61\x16\x5b\xbd\xf3\x94\x74\x90\xd0\x59\xa1\x93\x4e\x26\x85\xbf\x09\xdf\x4f\x03\x2d\x6c\xbd\xcf\xe6\x2d\xd7\x6d\xcb\x9e\xad\x9f\x2e\xfa\xa8\x74\x52\xe8\x6f\x0b\x37\xc8\xd2\x9d\x84\x09\x68\x9d\xfd\x75\xb1\xaa\x6c\x83\xe7\x71\xf8\x60\x6b\xf9\x1e\xc3\xf6\xb1\x1f\xd8\xf6\xa3\xc1\xfc\x0a\x00\x00\xff\xff\x24\x08\x64\xcd\x72\x0b\x00\x00")

func translateTemplatesMainGotmplBytes() ([]byte, error) {
	return bindataRead(
		_translateTemplatesMainGotmpl,
		"translate/templates/main.gotmpl",
	)
}

func translateTemplatesMainGotmpl() (*asset, error) {
	bytes, err := translateTemplatesMainGotmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "translate/templates/main.gotmpl", size: 2930, mode: os.FileMode(436), modTime: time.Unix(1541503049, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"translate/templates/defs.gotmpl": translateTemplatesDefsGotmpl,
	"translate/templates/main.gotmpl": translateTemplatesMainGotmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"translate": &bintree{nil, map[string]*bintree{
		"templates": &bintree{nil, map[string]*bintree{
			"defs.gotmpl": &bintree{translateTemplatesDefsGotmpl, map[string]*bintree{}},
			"main.gotmpl": &bintree{translateTemplatesMainGotmpl, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

