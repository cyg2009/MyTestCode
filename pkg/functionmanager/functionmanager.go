package functiomanager

import ( 

    "errors"
    "sync"
    "encoding/json"
    sf "MyTestCode/pkg/serverlessfunction"
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
    return ff, ok
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
    ret, err := json.Marshal(mgr.functionStore)

    if err != nil {
        return err.Error()
    }

    return string(ret)
}

// Singleton 
var instance *ServerlessFunctionManager
var once sync.Once

func GetFunctionManager() (*ServerlessFunctionManager) {

    once.Do( func() {
            instance = &ServerlessFunctionManager {
            functionStore: make(map[string]*sf.ServerlessFunction),
        }
    })

    return instance
}




