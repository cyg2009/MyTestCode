package prservice

import (
    "net/http"  
    "io/ioutil"
    "sync"
    "time"
    fm "../functionmanager"
)

func makeOKResponse(w http.ResponseWriter, body string){
    
    //w.Header().Set("Content-Type", "application/json")   
    w.Write([]byte(body))
}
func makeFailedResponse(w http.ResponseWriter, statusCode int, message string) {

   // bodyContent := slscommon.NewErrorMessage(statusCode, message).ToJsonString()    
    //w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write([]byte(message))
}

func ServeHTTPFetch(w http.ResponseWriter, req *http.Request) {

    functionId := req.URL.Query().Get("function")
     
    if  len(functionId) == 0 {
        makeFailedResponse(w, http.StatusBadRequest, "No function specified.")
        return
    }
    sf, ok := fm.GetFunctionManager().FetchFunction(functionId)
    if  ok == false {
        makeFailedResponse(w, http.StatusInternalServerError, "Failed to get the function package.")
        return
    }

    fm.GetFunctionManager().AddFunction(sf)
    body := "Fetch " + functionId + " successully!"
    makeOKResponse(w, body)
}

//This will receive a function package tar file and store it
func ServeHTTPAddFunction(w http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        makeFailedResponse(w, http.StatusInternalServerError, "Please POST this request.")
        return 
    }

    functionId := req.Header.Get("function")    
  
    if  len(functionId) == 0 {
        makeFailedResponse(w, http.StatusBadRequest, "No function specified.")
        return
    }
    
    data, err := ioutil.ReadAll(req.Body)
    if err != nil {
        makeFailedResponse(w, http.StatusInternalServerError, err.Error())
        return 
    }

    fm.GetFunctionManager().CreateFunction(functionId, data, "", nil)
    body := "Add " + functionId + " successully!"
    makeOKResponse(w, body)
}

func ServeHTTPInvoke(w http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        makeFailedResponse(w, http.StatusInternalServerError, "Please POST this request.")
        return 
    }
       
    functionId := req.Header.Get("function")
    if  len(functionId) == 0 {
        makeFailedResponse(w, http.StatusBadRequest, "No function specified.")
        return
    }

    evt, err := ioutil.ReadAll(req.Body)
    if err != nil {
        makeFailedResponse(w, http.StatusInternalServerError, err.Error())
        return 
    }
    //to simplify it , we treat the data as an js file index.js . TBD to replace it with a tar file




    // hack to return immediately
    //makeOKResponse(w, "Function " + functionId + " invoked successfully:" + string(evt[:]))
    //return
    
    if _, ok := fm.GetFunctionManager().GetFunction(functionId); ok == false {

            sf, ok2 := fm.GetFunctionManager().FetchFunction(functionId)
            if  ok2 == false {
                makeFailedResponse(w, http.StatusInternalServerError, "Failed to get the function package.")
                return
            }

            fm.GetFunctionManager().AddFunction(sf)

        // makeFailedResponse(w, http.StatusBadRequest, "Function " + functionId + "  not exists!")
        // return 
    }
            
    respData, _ := fm.GetFunctionManager().ExecuteFunction(functionId, evt)     
    makeOKResponse(w, respData)
    
}

func ServeHTTPConfig(w http.ResponseWriter, req *http.Request) {
    makeOKResponse(w, "TBD:Configuration")
}

func ServeHTTPRemove(w http.ResponseWriter, req *http.Request) {
    functionId := ""
    
    if req.Method == "POST" {
       functionId = req.Header.Get("function")
    }
    
    if req.Method == "GET" {
       functionId = req.URL.Query().Get("function")
    }
  
    if  len(functionId) == 0 {
        makeFailedResponse(w, http.StatusBadRequest, "No function specified.")
        return
    }
    
    if  ok := fm.GetFunctionManager().RemoveFunction(functionId); ok == false {
        makeFailedResponse(w, http.StatusBadRequest, "Fail to remove " + functionId)
        return 
    }

    makeOKResponse(w, functionId + " removed.")
}
func ServeHTTPInfo(w http.ResponseWriter, req *http.Request) {
    info := fm.GetFunctionManager().GetAllFunctionsJSON()
    makeOKResponse(w, info)
}
func ServeHTTPHealthCheck(w http.ResponseWriter, req *http.Request) {
    t := time.Now()
    body := "OK " + t.Format("20060102150405")
    makeOKResponse(w, body)
}

// Singleton 
var instance *http.ServeMux
var once sync.Once

func GetPrserviceHttpHandler() (*http.ServeMux) {

    once.Do( func() {
        instance = http.NewServeMux()   
        instance.HandleFunc("/health", ServeHTTPHealthCheck)
        instance.HandleFunc("/config", ServeHTTPConfig)
        instance.HandleFunc("/info", ServeHTTPInfo)
        instance.HandleFunc("/invoke", ServeHTTPInvoke)
        instance.HandleFunc("/fetch", ServeHTTPFetch)
        instance.HandleFunc("/add", ServeHTTPAddFunction)
        instance.HandleFunc("/remove", ServeHTTPRemove)
    })

    return instance
}
