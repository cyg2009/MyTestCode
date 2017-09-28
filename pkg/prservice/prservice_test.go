package prservice

import (
    "os"
    "bytes"
	"net/http"
	"testing"
    "net/http/httptest"
    fm "MyTestCode/pkg/functionmanager"
)

func TestMain(m *testing.M) {    

    os.Setenv("RUNTIME_ROOT", "/work/src/MyTestCode/runtime")
    
    code := m.Run() 

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
   
    functionId := "f1"

    //ingest function

    body :=  bytes.NewBufferString(`
       exports.handler = function(event) {
           console.log('got event:' + event)
           console.log('bye!')
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
        t.Errorf("Function-Fetch handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
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
