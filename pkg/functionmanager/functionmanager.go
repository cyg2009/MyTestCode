package functiomanager

import ( 
    "os"
    "io/ioutil"
    "errors"
    "sync"
    "encoding/json"
    sf "github.com/cyg2009/MyTestCode/pkg/serverlessfunction"
)

type ServerlessFunctionManager struct{
    functionStore map[string]*sf.ServerlessFunction
}

func (mgr *ServerlessFunctionManager) CreateFunction (functionId string, data []byte, sfCommand string, sfArguments []string)(*sf.ServerlessFunction, error) {
    ff, err := sf.CreateServerlessFunction(functionId, data, sfCommand, sfArguments) 
    if nil != err {
        return nil, err
    }
    
    mgr.functionStore[functionId] = ff        
    return ff, nil
}
func (mgr *ServerlessFunctionManager) GetFunction (functionId string) (*sf.ServerlessFunction, bool){

    ff, ok := mgr.functionStore[functionId]   
    if ok {
        return ff, ok
    }

    ff = sf.LoadFunction(functionId)
    if ff == nil {
        return nil, false
    }
    
    return ff, true
}

func (mgr *ServerlessFunctionManager) RemoveFunction (functionId string) (bool){
    ff, ok := mgr.functionStore[functionId]

    if ok {
       ff.Stop()
       delete(mgr.functionStore, functionId) 
    }

    return true
}

func (mgr *ServerlessFunctionManager) ExecuteFunction (functionId string, event []byte) (string, error) {
    
    ff, ok := mgr.functionStore[functionId]
    if ok {
        return ff.Trigger(event), nil
    }

    return  "", errors.New(functionId + " not exists!")
}

func (mgr *ServerlessFunctionManager) GetAllFunctionsJSON () (string){
    mgr.LoadAllFunctions()
    ret, err := json.Marshal(mgr.functionStore)

    if err != nil {
        return err.Error()
    }

    return string(ret)
}

func (mgr *ServerlessFunctionManager) LoadAllFunctions(){
    dest := os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }

    dest += "/func"
    files, err := ioutil.ReadDir(dest)
    if err != nil {
        return 
    }
    
    for _, file := range files {
        
        functionId := file.Name()   

        if file.IsDir() == false {
            continue
        }
        
        _, ok := mgr.functionStore[functionId]
        
        if ok == false {                  
            ff := sf.LoadFunction(functionId)
            mgr.functionStore[functionId] = ff     
        }
    }
}
// Singleton 
var instance *ServerlessFunctionManager
var once sync.Once

func GetFunctionManager() (*ServerlessFunctionManager) {

    once.Do( func() {
        instance = &ServerlessFunctionManager {
            functionStore: make(map[string]*sf.ServerlessFunction),
        }
        instance.LoadAllFunctions()
    })

    // load existing functions that not in the store
    instance.LoadAllFunctions()
    return instance
}




