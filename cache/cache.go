package cache

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Cache struct {
	opts options
}

type options struct {
	storePath string
}

type Option interface {
	apply(*options)
}

type storePathOption string

func (h storePathOption) apply(opts *options) {
	opts.storePath = string(opts.storePath)
}
func WithStorePath(path string) Option {
	return storePathOption(path)
}

func New(opts ...Option) *Cache {

	options := options{
		storePath: "cache",
	}

	for _, o := range opts {
		o.apply(&options)
	}

	return &Cache{opts: options}
}

func (c *Cache) Store(key string, value []byte) error {

	filePath := c.getFilePathFromKey(key)
	return ioutil.WriteFile(filePath, value, 0644)
}

func (c *Cache) Read(key string) ([]byte, error) {

	filePath := c.getFilePathFromKey(key)
	return ioutil.ReadFile(filePath)
}

func (c *Cache) Remove(key string) error {

	filePath := c.getFilePathFromKey(key)
	return os.Remove(filePath)
}

func (c *Cache) getFilePathFromKey(key string) string {

	hash := md5.Sum([]byte(key))
	return filepath.Join(c.opts.storePath, hex.EncodeToString(hash[:]))
}

func (c *Cache) RemoveAll() error {

	files, err := filepath.Glob(filepath.Join(c.opts.storePath, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) StoreUint64(key string, value uint64) error {
	return c.Store(key, []byte(fmt.Sprintf("%d", value)))
}

func (c *Cache) ReadUint64(key string) (uint64, error) {

	b, err := c.Read(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(b), 10, 64)
}

// Uses `\n` as delimiter
func (c *Cache) StoreStrings(key string, value []string) error {
	return c.Store(key, []byte(strings.Join(value, "\n")))
}

// Uses `\n` as delimiter
func (c *Cache) ReadStrings(key string) ([]string, error) {
	b, err := c.Read(key)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

// Only exported fields of an struct are stored
func (c *Cache) StoreAny(key string, value any) error {
	var bData bytes.Buffer
	enc := gob.NewEncoder(&bData)

	err := enc.Encode(value)
	if err != nil {
		return err
	}
	return c.Store(key, bData.Bytes())
}

// The second argument must be a reference to the a variable of the stored type
// Ex: 	output := []Mystruct{}
//
//	err := c.ReadAny("TEST_KEY", &output)
func (c *Cache) ReadAny(key string, output any) error {

	b, err := c.Read(key)
	if err != nil {
		return err
	}

	bData := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bData)

	return dec.Decode(output)
}
