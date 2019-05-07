package proc

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"os"
	"path"
	"regexp"
	"strings"
)

type Config struct {
	ProtoFile   string
	ProtoPath   string
	MessageName string
	OutputPath  string
	DoubleV     float64
	FloatV      float32
	Int32V      int32
	Int64V      int64
	Uint32V     uint32
	Uint64V     uint64
	Sint32V     int32
	Sint64V     int64
	Fixed32V    uint32
	Fixed64V    uint64
	Sfixed32V   int32
	Sfixed64V   int64
	BoolV       bool
	StringV     string
	BytesV      []byte
}
type Context struct {
	Config
	msgMatcher *regexp.Regexp
}

func (c *Context) SetConfig(cfg Config) {
	c.Config = cfg
	c.msgMatcher = regexp.MustCompile(cfg.MessageName)
}
func (c *Context) MatchMsgName(name string) bool {
	return c.msgMatcher.MatchString(name)
}
func (c *Context) GetPbFileName() string {
	return c.ProtoFile
}
func (c *Context) GetImportPath() []string {
	return strings.Split(c.ProtoPath, ";")
}
func (c *Context) MapVal(fd *desc.FieldDescriptor) (interface{}, error) {
	switch fd.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return c.BoolV, nil
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return c.BytesV, nil
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return c.DoubleV, nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return c.Fixed32V, nil
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return c.Fixed64V, nil
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return c.FloatV, nil
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return c.Int32V, nil
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return c.Int64V, nil
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return c.Uint32V, nil
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return c.Uint64V, nil
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return c.Sint32V, nil
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return c.Sint64V, nil
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return c.Sfixed32V, nil
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return c.Sfixed64V, nil
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return c.StringV, nil
	default:
		return nil, errors.New("unsupported type: " + fd.GetType().String())
	}
}
func (c *Context) WriteJson(name string, jsonStr string) {
	if !c.MatchMsgName(name) {
		return
	}
	var output *os.File
	if strings.EqualFold(c.OutputPath, "stdout") {
		output = os.Stdout
		fmt.Fprintf(output, "------------------%s.json------------------\n", name)
		defer output.WriteString("\n")
	} else {
		fm := path.Join(c.OutputPath, name)
		if isExist(fm) {
			fmt.Fprintf(os.Stderr, "文件[%s]已经存在.", fm)
			return
		} else {
			var e error
			output, e = os.Create(fm)
			defer output.Close()
			if e != nil {
				fmt.Fprintf(os.Stderr, "创建文件[%s]失败. error:%+v", fm, e)
				return
			}
		}
	}
	output.WriteString(jsonStr)
}
func isExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
