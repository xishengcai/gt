package gt

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"unsafe"
)

// Decoder is the decoding interface
type Decoder interface {
	Decode(v interface{}) error
	Value() interface{}
}

// BodyDecode body decoder structure
type BodyDecode struct {
	r   io.Reader
	obj interface{}
}

// NewBodyDecode create a new body decoder
func NewBodyDecode(r io.Reader) Decoder {
	if r == nil {
		return nil
	}

	return &BodyDecode{r: r}
}

var convertBodyFunc = map[reflect.Kind]convert{
	reflect.Uint:    {bitSize: 0, cb: setUintField},
	reflect.Uint8:   {bitSize: 8, cb: setUintField},
	reflect.Uint16:  {bitSize: 16, cb: setUintField},
	reflect.Uint32:  {bitSize: 32, cb: setUintField},
	reflect.Uint64:  {bitSize: 64, cb: setUintField},
	reflect.Int:     {bitSize: 0, cb: setIntField},
	reflect.Int8:    {bitSize: 8, cb: setIntField},
	reflect.Int16:   {bitSize: 16, cb: setIntField},
	reflect.Int32:   {bitSize: 32, cb: setIntField},
	reflect.Int64:   {bitSize: 64, cb: setIntDurationField},
	reflect.Float32: {bitSize: 32, cb: setFloatField},
	reflect.Float64: {bitSize: 64, cb: setFloatField},
}

// Decode body decoder
func (b *BodyDecode) Decode(v interface{}) error {
	return Body(b.r, v)
}

// Value ...
func (b *BodyDecode) Value() interface{} {
	return b.obj
}

// Body body decoder
func Body(r io.Reader, obj interface{}) error {
	if w, ok := obj.(io.Writer); ok {
		_, err := io.Copy(w, r)
		return err
	}

	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	value := LoopElem(reflect.ValueOf(obj))

	if value.Kind() == reflect.String {
		value.SetString(BytesToString(all))
		return nil
	}

	if _, ok := value.Interface().([]byte); ok {
		value.SetBytes(all)
		return nil
	}

	fn, ok := convertBodyFunc[value.Kind()]
	if ok {
		return fn.cb(BytesToString(all), fn.bitSize, emptyField, value)
	}

	return fmt.Errorf("type (%T) %s", value, ErrUnknownType)
}

// LoopElem 不停地对指针解引用
func LoopElem(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}

	return v
}

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
