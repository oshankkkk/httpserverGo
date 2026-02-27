package main
import (
	"fmt"
	"net"
	"os"
)
func check(err error){
	if err!=nil{
		fmt.Println("error")
	panic(err)
	}
}
func handleConn(file net.Conn) []string{
	stream:=make([]byte,8)
	var sentence string
	var mylist []string
	for{
		count,err:=file.Read(stream)
		check(err)
		for _,letterrune:=range string(stream[:count]){
			letter:=string(letterrune)
			if letter=="\n"{
				mylist = append(mylist, sentence)
				sentence=""
			}else {
				sentence+=letter
			}
		}
		//break in end of request
		if count<8{
			break
		}

	}
return mylist
} 
func writeTofile(request []string){
	logfile,err:=os.OpenFile("serverlogs.txt",os.O_APPEND|os.O_CREATE|os.O_RDWR,0666)
	check(err)

	fmt.Println("file was made nicely")
	for _,line:=range request{
		//fmt.Println(line)
		_,err:=logfile.WriteString(line+"\n")
		check(err)

	}
}
func sendResponse(file net.Conn){
	msg:="HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
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
		request:=handleConn(file)
		writeTofile(request)
		sendResponse(file)

	}
}

