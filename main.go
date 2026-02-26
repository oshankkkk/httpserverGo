package main

import (
	"fmt"
	"net"
	"io"
)
func check(err error){
	if err!=nil{
		fmt.Println("error")
	panic(err)
	}
}
func main(){
	listner,err:=net.Listen("tcp",":8080")
	check(err)
	file,err:=listner.Accept()
	check(err)
	
stream:=make([]byte,8)

	var sentence string
	var mylist []string
	for{
	count,err:=file.Read(stream)
	if err==io.EOF{
		break
	}
	for _,letterrune:=range string(stream[:count]){
		letter:=string(letterrune)
		if letter=="\n"{
			mylist = append(mylist, sentence)
			sentence=""
		}else{
			sentence+=letter
		}


	}

}
for _,line:=range mylist{
	fmt.Println(line)
}

}
