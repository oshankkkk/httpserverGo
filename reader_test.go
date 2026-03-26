package main

import (
	"strings"
	"testing"
	"fmt"
)



func TestFormatBytes(t *testing.T) {
a:=`POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: */ /*
Content-Length: 21

{"flavor":"dark mode"}
`

	buff := formatBytes([]byte(a))
	result:=strings.Join(buff,"\n")
	if strings.TrimSpace(a)!=result{
		t.Errorf("expected %s, got %s", a,result)
	}
}

func TestHeaderParser(t *testing.T){

a:=`POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: */ /*
Content-Length: 21`

b:=strings.Split(a, "\n")
sline:=StartLine{
	method: "POST",
	path: "/coffee",
	version: "HTTP/1.1",
}

startline,err,contlen:=HeaderParser(b)
if err!=nil{
	fmt.Println("yo")
}
if startline!=sline{
		t.Errorf("expected POST, got %s",startline.method )
}
if contlen!=21{
		t.Errorf("expected a 21, got %s",startline.version)
}
}


