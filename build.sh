p=$s0
go build -buildmode=plugin -ldflags "-pluginpath=flex/plugins/check-host"  -o check-host1.so ./plugins/check-host
go build -buildmode=plugin -ldflags "-pluginpath=flex/plugins/content-type" -o content-type1.so ./plugins/content-type
go build -buildmode=plugin -ldflags "-pluginpath=flex/plugins/path" -o path2.so ./plugins/path


# go build -buildmode=plugin  ./plugins/check-host
# go build -buildmode=plugin  ./plugins/content-type
# go build -buildmode=plugin  ./plugins/path

mv *.so ./sos