package functiomanager

import (
	"testing"
	"os"
	"io/ioutil"
)

func TestMain(m *testing.M) {    
	wd, _ := os.Getwd()
	os.Setenv("RUNTIME_ROOT",  wd)	
	os.Setenv("RUNTIME_LAMBDA", wd + "/../../runtime/bin/lambda-run")
	dest := wd + "/func"
	
	_, err := os.Stat(dest)	
	if err == nil {
		os.RemoveAll(dest)		
	}

    os.Mkdir(dest, os.ModeDir)
    code := m.Run() 
    
	_, err = os.Stat(dest)	
	if err == nil {
		os.RemoveAll(dest)
	}

    os.Exit(code)
}


func TestGetExistingFunctions(t *testing.T){
	functionId := "f2"
	mgr := GetFunctionManager()
	_, ok := mgr.functionStore[functionId]
	if ok {
		 t.Errorf("Function already exists: " + functionId)
		  
	}

    dest := os.Getenv("RUNTIME_ROOT") + "/func/" + functionId
 
    data :=  []byte(`
       exports.handler = function(event) {
		   console.log((new Date()).toString())
           console.log('f2 received an event:' + JSON.stringify(event))           
       }
	`)

	os.RemoveAll(dest)
	os.Mkdir(dest, os.ModeDir)

	dest = dest + "/index.js"
	err := ioutil.WriteFile(dest, data, 0644)
	
    if err != nil {
         t.Errorf("Fail  to create function " + functionId)
	}
		
	mgr = GetFunctionManager()

	_, ok = mgr.GetFunction(functionId)

	if ok == false {
		 t.Errorf("Fail  to get existing function " + functionId)
	}
 
}
