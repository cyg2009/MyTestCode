package mytest

import (
	"testing"
	"os"
	"fmt"
)

func TestMain(m *testing.M) {    

    fmt.Println("Start the test suite...")
    code := m.Run() 

    os.Exit(code)
}

func Test1(t *testing.T){
    t.Log("firt test is done.")	
}