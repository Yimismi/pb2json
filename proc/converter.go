package proc

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	_ "github.com/jhump/protoreflect/dynamic"
	"os"
	"reflect"
)

type Converter struct {
	ctx    *Context
	msgMap map[string]*dynamic.Message
}

func (c *Converter) Convert(ctx *Context) {
	c.ctx = ctx
	c.msgMap = make(map[string]*dynamic.Message)
	p := protoparse.Parser{ImportPaths: ctx.GetImportPath()}
	fds, e := p.ParseFiles(ctx.GetPbFileName())
	if e != nil {
		fmt.Fprintf(os.Stderr, "解析pb文件[%s]错误. error:%+v", ctx.GetPbFileName(), e)
		return
	}
	fd := fds[0]
	for _, md := range fd.GetMessageTypes() {
		c.convertMsg(md, make(map[string]bool))
	}
	for name, msg := range c.msgMap {
		bs, e := msg.MarshalJSONPB(&jsonpb.Marshaler{Indent: "  ", OrigName: true})
		if e != nil {
			fmt.Fprintf(os.Stderr, "message[%s]转json错误. error:%+v", name, e)
			continue
		}
		ctx.WriteJson(name, string(bs))
	}
}
func (c *Converter) convertMsg(md *desc.MessageDescriptor, path map[string]bool) *dynamic.Message {
	if m, ok := c.msgMap[md.GetName()]; ok {
		return m
	}
	if _, ok := path[md.GetName()]; ok {
		return nil
	}
	path[md.GetName()] = true
	dmsg := dynamic.NewMessage(md)
	for _, fid := range md.GetFields() {
		if fid.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			msg := c.convertMsg(fid.GetMessageType(), path)
			if msg != nil {
				setMsgField(dmsg, fid, msg)
			}
		} else {
			val, e := c.mapVal(fid)
			if e != nil {
				continue
			}
			setMsgField(dmsg, fid, val)
		}
	}
	c.msgMap[md.GetName()] = dmsg
	delete(path, md.GetName())
	return dmsg
}
func setMsgField(dmsg *dynamic.Message, fd *desc.FieldDescriptor, val interface{}) {
	if isNil(val) {
		return
	}
	if fd.IsRepeated() {
		dmsg.AddRepeatedField(fd, val)
	} else {
		dmsg.SetField(fd, val)
	}
}
func (c *Converter) mapVal(fd *desc.FieldDescriptor) (interface{}, error) {
	return c.ctx.MapVal(fd)
}
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
