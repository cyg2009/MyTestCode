package prservice

import (
    "os"
    "bytes"
    "io/ioutil"
	"net/http"
	"testing"
    "net/http/httptest"
    fm "github.com/cyg2009/MyTestCode/pkg/functionmanager"
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

func TestHealthCheckHandler(t *testing.T){

    req, err := http.NewRequest("POST", "/health", nil)
    if err != nil {
        t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	handler := GetPrserviceHttpHandler()
	
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    t.Log(rr.Body.String())
}

func TestOtherHandler(t *testing.T){

    req, err := http.NewRequest("POST", "/other", nil)
    if err != nil {
        t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	handler := GetPrserviceHttpHandler()
	
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
	}
		
    t.Log(rr.Body.String())
}


func TestIngestAndExecuteFunction(t *testing.T){
    
    fgr := fm.GetFunctionManager()
   
    functionId := "f1:1.0"

    //ingest function
    body :=  bytes.NewBufferString(`
       exports.handler = function(event) {
           console.log('f1:1.0 got an event:' + JSON.stringify(event))
           console.log((new Date()).toString())
       }
    `)

    req, err := http.NewRequest("POST", "/add", body)
    if err != nil {
        t.Fatal(err)
	}
    req.Header.Add("function",functionId)

	rr := httptest.NewRecorder()
	
	handler := GetPrserviceHttpHandler()
	
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("Function-Add handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
        t.Log(rr.Body.String())
        return
	}
        
    sf, ok := fgr.GetFunction(functionId)
    if ok == false || sf == nil  {
	   t.Errorf("Fail to get function.")
   }

    //Execution function
    evt := []byte(`{"name":"cacia", "age":"19"}`)
    req, err = http.NewRequest("POST", "/invoke", bytes.NewBuffer(evt))
    if err != nil {
        t.Fatal(err)
	}

    req.Header.Add("function",functionId)
    req.Header.Set("Content-Type","applicaiton/json")

	rr = httptest.NewRecorder()	
	handler = GetPrserviceHttpHandler()
	
    handler.ServeHTTP(rr, req)
    
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("Function-Invoke handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}
  
    t.Log(rr.Body.String())   
}

func TestExecuteExistingFunction(t *testing.T){
    // Use a random id would be better
    functionId := "f2:1.0"

    dest := os.Getenv("RUNTIME_ROOT") + "/func/" + functionId
 
    data :=  []byte(`
       exports.handler = function(event) {
		   console.log((new Date()).toString())
           console.log('f2:1.0 received an event:' + JSON.stringify(event))           
       }
	`)

	os.RemoveAll(dest)
	os.Mkdir(dest, os.ModeDir)

	dest = dest + "/index.js"
	err := ioutil.WriteFile(dest, data, 0644)
	
    if err != nil {
         t.Errorf("Fail  to create function " + functionId)
	}		
    
    //Execution function
    evt := []byte(`{"name":"cacia", "age":"19"}`)
    req, err := http.NewRequest("POST", "/invoke", bytes.NewBuffer(evt))
    if err != nil {
        t.Fatal(err)
	}

    req.Header.Add("function",functionId)
    req.Header.Set("Content-Type","applicaiton/json")

	rr := httptest.NewRecorder()	
	handler := GetPrserviceHttpHandler()
	
    handler.ServeHTTP(rr, req)
    
	if status := rr.Code; status != http.StatusOK {
        t.Errorf("Function-Invoke handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}
  
    t.Log(rr.Body.String())   
}