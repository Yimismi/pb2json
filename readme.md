```
将pb文件中message转换成一个json case，直接输入命令显示帮助

Usage: pb2json [-proto_file] [-proto_path] [-message_name] [-output_path]

Options:
  -message_name string
        message名，支持正则表达式 (default ".*")
  -output_path string
        输出路径，默认为标准输出流 (default "stdout")
  -proto_file string
        pb文件 (default "./test.proto")
  -proto_path string
        import路径，多个路径用';'间隔 (default "./")
```