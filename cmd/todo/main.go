package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"protoPrac2/todo"
	"strings"
)

func main(){

	flag.Parse()
	if flag.NArg()<1{
		fmt.Fprintln(os.Stderr,"missing subcommand: list or add")
		os.Exit(1)
	}

	conn,err:=grpc.Dial(":8888",grpc.WithInsecure())
	if err!=nil{
		log.Fatalf("could not connect to backend: %v",err)
	}
	client:=todo.NewTasksClient(conn)
	switch cmd:=flag.Arg(0);cmd {
	case "list":
		err=list(context.Background(),client)
	case "add":
		err=add(strings.Join(flag.Args()[1:]," "))
	default:
		err=fmt.Errorf("unknown subcommand %s",cmd)
	}
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}
type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)
func add(text string) error{

	return fmt.Errorf("add not implemented")
}
func list(ctx context.Context,client todo.TasksClient)error{
	l,err:=client.List(ctx,&todo.Void{})
	if err!=nil{
		return fmt.Errorf("could not fetch tasks: %v",err)
	}
	for _,t:=range l.Tasks{
		if t.Done {
			fmt.Printf("done")
		}else{
			fmt.Printf("not done ")
		}
		fmt.Printf("%s\n",t.Text)
	}
	return nil
}