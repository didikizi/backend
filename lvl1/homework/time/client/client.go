package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//1
	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Print(string(buf))
	}
	//1
}

/*
1
Или можно сделать вывод без обратки пустого значения
for {
	io.Copy(os.Stdout, conn)
}
*/
