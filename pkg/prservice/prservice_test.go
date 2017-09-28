package prservice
import (
    "os"

	"net/http"
	"testing"
    "net/http/httptest"
   
)

func TestMain(m *testing.M) {    

    os.Setenv("RUNTIME_ROOT", "/work/src/processrouter/runtime")
    
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


/*
func TestFetchAndExecuteOneFunctionOnce(t *testing.T){
    
    fgr := fm.GetFunctionManager()
    pkgMgr := &fm.ServerlessFunctionPackageManager{}
    baseUrl := "http://10.21.119.117:5000/v2/serverless/"
    pkgMgr.SetFunctionStoreUrl(baseUrl)
    fgr.SetFunctionPackageManager(pkgMgr)
    
    functionId := "f3:1.0.0"

    //Fetch function
    req, err := http.NewRequest("POST", "/fetch", nil)
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

func TestFetchAndExecuteOneFunctionTwice(t *testing.T){
    
    fgr := fm.GetFunctionManager()
    pkgMgr := &fm.ServerlessFunctionPackageManager{}
    baseUrl := "http://10.21.119.117:5000/v2/serverless/"
    pkgMgr.SetFunctionStoreUrl(baseUrl)
    fgr.SetFunctionPackageManager(pkgMgr)
    
    functionId := "f3:1.0.0"
    req, err := http.NewRequest("POST", "/fetch", nil)
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

    //Execution f3
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
    
    time.Sleep(100*time.Millisecond)
    
    t.Log("Execution f3 again -------------------------")
    evt = []byte(`{"name":"cacia2", "age":"29"}`)
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

func TestFetchAndExecuteTwoFunctions(t *testing.T){
    
    fgr := fm.GetFunctionManager()
    pkgMgr := &fm.ServerlessFunctionPackageManager{}
    baseUrl := "http://10.21.119.117:5000/v2/serverless/"
    pkgMgr.SetFunctionStoreUrl(baseUrl)
    fgr.SetFunctionPackageManager(pkgMgr)
    
    functionId := "f3:1.0.0"
    //Fetch f3
    req, err := http.NewRequest("POST", "/fetch", nil)
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
	   t.Errorf("Fail to get function " + functionId)
   }

    //Execution f3
    evt := []byte(`{"name":"cacia", "age":"17"}`)
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

   functionId = "f4:1.0.0"
    
   //Fetch f4
    req, err = http.NewRequest("POST", "/fetch", nil)
    if err != nil {
        t.Fatal(err)
	}

    req.Header.Add("function",functionId)

	rr = httptest.NewRecorder()
	
	handler = GetPrserviceHttpHandler()
	
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("Function-Fetch handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
	}
        
    sf, ok = fgr.GetFunction(functionId)
    if ok == false || sf == nil  {
	   t.Errorf("Fail to get function " + functionId)
   }

    //Execution f4
    evt = []byte(`{"name":"caciasister", "age":"20"}`)
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
*/