package serverlessfunction

import (
	"testing"
	"os"
)

func TestMain(m *testing.M) {    
	wd, _ := os.Getwd()
	os.Setenv("RUNTIME_ROOT",  wd)	
	os.Setenv("RUNTIME_LAMBDA", "/work/src/MyTestCode/runtime/bin/lambda-run")
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

func TestIngestAndExecuteFunction(t *testing.T){
    functionId := "f1"
	 // Save the function js file
    dest := os.Getenv("RUNTIME_ROOT") + "/func/" + functionId + "/index.js"
	t.Log(dest)
 
    data :=  []byte(`
       exports.handler = function(event) {
           console.log('got event:' + event)
           console.log('bye!')
       }
	`)
	
    lambda := os.Getenv("RUNTIME_LAMBDA") 
    sf, err := CreateServerlessFunction(functionId, data, "", []string{lambda, dest})
    if  err != nil  {
	   t.Errorf("Fail to create function:" + err.Error())
    } 
	
	//Execution of the function
	sf.Start()	
	t.Log(sf)

	resp := sf.Trigger([]byte(`{"name":"hayoung"}`))
	t.Log(resp)
	
	sf.Stop()
	t.Log(sf)
}
