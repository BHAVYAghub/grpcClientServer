package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"protoPrac2/todo"
)

type taskServer struct {
	todo.UnimplementedTasksServer
}
func (t taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error){
	return nil,fmt.Errorf("not implemented till now")
}

func main(){
	srv:=grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(srv,tasks)
	l,err:=net.Listen("tcp",":8888")
	if err!=nil{
		log.Fatalf("could not listen to :8888: %v",err)
	}
	log.Fatal(srv.Serve(l))
}
type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)
func add(text string) error{
	task:=&todo.Task{
		Text: text,
		Done: false,
	}
	b,err:=proto.Marshal(task)
	if err!=nil{
		return fmt.Errorf("could not encode task: %v",task)
	}
	f,err:=os.OpenFile(dbPath,os.O_WRONLY | os.O_CREATE| os.O_APPEND,0666)
	if err!=nil{
		return fmt.Errorf("could not open file %s : %v",dbPath,err)
	}
	if err:=gob.NewEncoder(f).Encode(int64(len(b)));err!=nil{
		return fmt.Errorf("could not encode length of the message: %v",err)
	}
	_,err=f.Write(b)
	if err!=nil{
		return fmt.Errorf("could not write task to file: %v",err)
	}
	if err:=f.Close();err!=nil{
		return fmt.Errorf("could not close file %s: %v",dbPath,err)
	}
	return nil
}
func list()error{

	b,err := ioutil.ReadFile(dbPath)
	if err!=nil{
		return fmt.Errorf("could not read %s: %v",dbPath,err)
	}
	var task todo.Task

	for{
		if len(b)==0{
			return nil
		}else if len(b)<4{
			return fmt.Errorf("remaining odd %d bytes",len(b))
		}
		var length int64
		if err:=gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length);err!=nil{
			return fmt.Errorf("could not decode message length: %v",err)
		}
		b=b[4:]
		if err:=proto.Unmarshal(b[:length],&task);err!=nil{
			return fmt.Errorf("could not raed task: %v",err)
		}
		b=b[length:]
		if task.Done {
			fmt.Printf("done")
		}else{
			fmt.Printf("not done ")
		}
		fmt.Printf("%s\n",task.Text)
	}
}