package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

type StartLine struct{
	method string
	path string
	version string
}

func check(err error){
	if err!=nil{
		fmt.Println("error")
	panic(err)
	}
}
/*
POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: *//*
Content-Length: 21

{"flavor":"dark mode"}
*/

func readConnection(file net.Conn) []byte{
	stream:=make([]byte,1024)
	buff := []byte{}
	for{
		count,err:=file.Read(stream) 
		check(err)		
		buff = append(buff, stream[:count]...)	
		index:= bytes.Index(buff,[]byte("\r\n\r\n"))
		if index!=-1{
			return buff[:index]
		 }
	}
} 

func sendResponse(file net.Conn){
	msg:="HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	n,err:=file.Write([]byte(msg))
	check(err)
	fmt.Println(n)
	fmt.Println("http response send")


}

func formatBytes(header []byte) []string {
buff:=[]string{}
	var temp string
	for _,value:=range string(header){ 
		if string(value)=="\n"{
			temp=strings.TrimSpace(temp)
			buff = append(buff,temp) 
			temp=""
		}else{
		temp+=string(value)
		}

	}
temp = strings.TrimSpace(temp)
	if temp != "" {
		buff = append(buff, temp)
	}
	return buff
}

func headerParser(buff []string) (StartLine,error){
	startlinestring:=buff[0]
	templist:=strings.Split(startlinestring," ")	
	var startline StartLine
	startline.method=templist[0]
	startline.path=templist[1]
	startline.version=templist[2]
	fmt.Println("these are the statss",startline.version,startline.method)
	for i:=0;i<len(startline.method);i++{
		if !unicode.IsUpper(rune(startline.method[i])){
			return StartLine{},errors.New("incorrect http method")

	} 
}
	fmt.Println(startline.version, "this is the path" )
	if startline.version!="HTTP/1.1"{
			return StartLine{},errors.New("wrong http version")
	}
	// buf[1:] are headers
	headermap,hasmessagebody:=headerfieldParser(buff[1:])
	fmt.Println("this has a message body: ",hasmessagebody)
	for key,value:=range headermap{
		fmt.Println(key,":",value)

	}
	return startline,nil
	}
	

func headerfieldParser(buff []string)(map[string]string,bool){

	headermap:=make(map[string]string)
	seperator:=":"
	var hasmessegebody bool
	for _,value:=range buff{
		value=string(value)
		index:=strings.Index(value,seperator)	
		fieldName:=value[:index]
		if strings.ToLower(fieldName)=="content-length"{
			hasmessegebody=true
		}
		fieldValue:=value[index+1:]
		headermap[fieldName]=strings.TrimSpace(fieldValue)

	}
	return headermap,hasmessegebody

}

/*
The normal procedure for parsing an HTTP message is to read the start-line into a structure, 
read each header field line into a hash table by field name until the empty line, and then use the parsed data to determine if a message body is expected. 
If a message body has been indicated, then it is read as a stream until an amount of octets equal to the message body length is read or the connection is closed.
*/


func writeTofile(request []string){
	logfile,err:=os.OpenFile("serverlogs.txt",os.O_APPEND|os.O_CREATE|os.O_RDWR,0666)
	check(err)
	fmt.Println("passing bout to happen")
	//parse(request[0])
	fmt.Println("passing has happen")
	fmt.Println("file was made nicely")
	for _,line:=range request{
		//fmt.Println(line)
		_,err:=logfile.WriteString(line+"\n")
		check(err)
	}
}

func parse(startLine string){
	headermap:=make(map[string]string)
	complist:=strings.Split(startLine," ")	
	fmt.Println("method",complist[0],"path",complist[1])
	headermap["method"]=complist[0]
	headermap["path"]=complist[1]
	headermap["version"]=complist[2]
}



func main(){
	listner,err:=net.Listen("tcp",":8080")
	check(err)
	for{
		file,err:=listner.Accept()
		check(err)
		request:=readConnection(file)
		stringrequest:=formatBytes(request)
		writeTofile(stringrequest)
		startline,err:=headerParser(stringrequest)
		check(err)
		fmt.Println(startline.method)
		fmt.Println(startline.path)
		fmt.Println(startline.version)


		sendResponse(file)
	}
}

