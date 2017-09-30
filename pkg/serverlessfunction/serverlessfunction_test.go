package serverlessfunction

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

func TestLoadFunction(t *testing.T){

    functionId := "f3"
    dest := os.Getenv("RUNTIME_ROOT") + "/func/" + functionId
	t.Log(dest)
 
    data :=  []byte(`
       exports.handler = function(event) {
		   console.log((new Date()).toString())
           console.log('f3 received an event:' + JSON.stringify(event))           
       }
	`)
	
	os.RemoveAll(dest)
	os.Mkdir(dest, os.ModeDir)
 
    dest = dest + "/index.js"
	err := ioutil.WriteFile(dest, data, 0644)
	
    if err != nil {
         t.Errorf("Fail  to create function " + functionId)
    }
	
    ff := LoadFunction(functionId)
	
	if ff == nil  {
		 t.Errorf("fail to load function " + functionId)
	}
}


func TestIngestAndExecuteFunction(t *testing.T){
    functionId := "f1"
	 // Save the function js file
    dest := os.Getenv("RUNTIME_ROOT") + "/func/" + functionId + "/index.js"
	t.Log(dest)
 
    data :=  []byte(`
       exports.handler = function(event) {
		   console.log((new Date()).toString())
           console.log('f1 received an event:' + JSON.stringify(event))           
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
