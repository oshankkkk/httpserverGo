package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
	"strconv"
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

func readConnection(file net.Conn) {
	stream:=make([]byte,1024)
	buff := []byte{}
	var contentlength int
	contentlength=0
	for{
		count,err:=file.Read(stream) 
		check(err)		
		buff = append(buff, stream[:count]...)	
		index:= bytes.Index(buff,[]byte("\r\n\r\n"))
		//Do not stop the function that reads the bytes instead read until header end, parse the then if the body is there start reading else stop
		if index!=-1{
		stringrequest:=formatBytes(buff[:index])
		writeTofile(stringrequest)
		startline,err,a:=headerParser(stringrequest)
		check(err)
		fmt.Println(startline.method)
		fmt.Println(startline.path)
		fmt.Println(startline.version)
		contentlength=a	
		fmt.Println("end of the headerParser", contentlength)

		}
	//parse the body if it is here by reading the test of the remainig bytes from the connection
		//body is basically buff[index:]
		if contentlength!=0{
			// we are going to check the content legth, if the reminder of the bytes in buf[index:] is smaller than the content length this mean there is more bytes to be read
			// so we go again through the loop and read the byte 
			fmt.Println("comes to cl check", contentlength)
			if len(buff[index+4:])<contentlength{
				fmt.Println("comes to buf len check",len(buff[index+4:]))
			continue
			}
			body:=bodyParser(buff[index+4:])
			fmt.Println("this is the body")
			fmt.Println(body)
			break
		}else{
		break
		}
	}
} 
func bodyParser(buff []byte) string{
	fullbodystring :=string(buff)
	return fullbodystring
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

func headerParser(buff []string) (StartLine,error,int){
	startlinestring:=buff[0]
	templist:=strings.Split(startlinestring," ")	
	var startline StartLine
	startline.method=templist[0]
	startline.path=templist[1]
	startline.version=templist[2]
	fmt.Println("these are the statss",startline.version,startline.method)
	for i:=0;i<len(startline.method);i++{
		if !unicode.IsUpper(rune(startline.method[i])){
			return StartLine{},errors.New("incorrect http method"),0
		} 
	}
	fmt.Println(startline.version, "this is the path" )
	if startline.version!="HTTP/1.1"{
		return StartLine{},errors.New("wrong http version"),0
	}
	// buf[1:] are headers
	headermap,contentlength:=headerfieldParser(buff[1:])

	for key,value:=range headermap{
		fmt.Println(key,":",value)
	}
	return startline,nil,contentlength
}

func fieldNameValidation(fieldName string) error{
	for _,value:=range fieldName{
		if !(value>='a' && value<='z'|| value>='A' && value<='Z'||strings.ContainsRune("!#$%&'*+-.^_`|~",value)||value>=0 && value<=9){
			fmt.Println(fieldName)
			return errors.New("invalid fieldName")
		}
	}

	return nil
}
func headerfieldParser(buff []string)(map[string]string,int){
	headermap:=make(map[string]string)
	seperator:=":"
	var contentlength int
	for _,value:=range buff{
		fmt.Println("value",value)
		value=string(value)
		fmt.Println("value string",value)
		index:=strings.Index(value,seperator)	
		fieldName:=value[:index]
		err:=fieldNameValidation(fieldName)	
		check(err)
		
		fieldValue:=value[index+1:]

		if strings.ToLower(fieldName)=="content-length"{

			contentlength,err=strconv.Atoi(fieldValue)
			check(err)

		}
		_,ok:=headermap[fieldName]
		if ok==false{
			headermap[fieldName]=strings.TrimSpace(fieldValue)
		}else{

			headermap[fieldName]+=","+strings.TrimSpace(fieldValue)
		}
	}
	return headermap,contentlength

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
	fmt.Println("passing done")
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
func writingStatusCode(code int){
	switch code {
		case 200:
		fmt.Println("god request")
		case 400:
		fmt.Println("bad one")
		case 500:
		fmt.Println("server err")
	}

		
}




func sendResponse(file net.Conn){
	fmt.Println("came to the respond e base")
	//msg2:="HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	msg:="HTTP/1.1 200 OK\r\n\r\nHello\r\n"
	n,err:=file.Write([]byte(msg))
	check(err)
	fmt.Println(n)
	fmt.Println("http response send")
}



func main(){
	listner,err:=net.Listen("tcp",":8080")
	check(err)
	for{
		file,err:=listner.Accept()
		check(err)
		readConnection(file)
		sendResponse(file)

	}
}

