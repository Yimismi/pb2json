package main

import (
	"flag"
	"fmt"
	"os"
	pb2json "pb2json/proc"
)

var c = pb2json.Config{
	DoubleV:   1.0,
	FloatV:    1.0,
	Int32V:    1,
	Int64V:    1,
	Uint32V:   1,
	Uint64V:   1,
	Sint32V:   1,
	Sint64V:   1,
	Fixed32V:  1,
	Fixed64V:  1,
	Sfixed32V: 1,
	Sfixed64V: 1,
	BoolV:     true,
	StringV:   "a_string",
	BytesV:    []byte("a_bytes"),
}

func init() {
	flag.StringVar(&c.ProtoFile, "proto_file", "./test.proto", "pb文件")
	flag.StringVar(&c.ProtoPath, "proto_path", "./", "import路径，多个路径用';'间隔")
	flag.StringVar(&c.MessageName, "message_name", ".*", "message名，支持正则表达式")
	flag.StringVar(&c.OutputPath, "output_path", "stdout", "输出路径，默认为标准输出流")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: pb2json [-proto_file]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if len(os.Args) <= 1 {
		usage()
		return
	}
	ctx := new(pb2json.Context)
	ctx.SetConfig(c)
	cvt := new(pb2json.Converter)
	cvt.Convert(ctx)
}
