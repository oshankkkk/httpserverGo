package main

import (
	"fmt"
	"io"
	"os"
)
func check(err error){
	if err!=nil{
		fmt.Println("error")
	panic(err)
	}
}
func main(){
	fmt.Println("this is my http server")
	file,err:=os.Open("message.txt")	
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
