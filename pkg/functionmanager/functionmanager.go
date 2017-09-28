package functiomanager

import ( 
    "io/ioutil"
    "errors"
    "sync"
    "os"
    "encoding/json"
    "path/filepath"
    sf "MyTestCode/pkg/serverlessfunction"
)


type ServerlessFunctionManager struct{
    functionStore map[string]*sf.ServerlessFunction
}

func (mgr *ServerlessFunctionManager) CreateFunction (functionId string, data []byte, sfCommand string, sfArguments []string)(*sf.ServerlessFunction, bool) {

    if len(functionId) == 0 || len(data) == 0 {
        return nil, false
    }

    // default value
    if len(sfCommand) == 0 {
        sfCommand = "node"
    } 
 
    // Save the function js file
    var dest = os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }

    dest = dest + "/func/" + functionId
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        os.Mkdir(dest, os.ModeDir)
    } else {
        RemoveContents(dest)
    }
 
    dest = dest + "/index.js"
    err := ioutil.WriteFile(dest, data, 0644)
    if err != nil {
        return nil, false
    }
   
    if nil == sfArguments ||  len(sfArguments) == 0 {
        sfArguments = []string{"lambda-run", dest}
    }

    ff := sf.NewServerlessFunction(functionId, sfCommand, sfArguments)
 
    mgr.functionStore[functionId] = ff
    return ff, true
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

func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
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




