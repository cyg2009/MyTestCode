package functiomanager

import (
    "bufio"
    "bytes"
    "strings"
    "io"      
    "io/ioutil"
    "syscall"
    "os/exec"
    "errors"
    "sync"
    "os"
    // "os/signal"
    "net/http"
    "encoding/json"
    "path/filepath"
)

type ServerlessFunction struct{
    id string
    input io.Writer
    outputReader *bufio.Reader
    cmd  *exec.Cmd
    command string
    args []string
    started bool
}

func (sf *ServerlessFunction) Start(){

   if !sf.started {

        sf.started = true
        sf.cmd = exec.Command(sf.command, sf.args...)
        sf.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
        var pr io.Reader 
        pr, sf.input = io.Pipe()
        sf.cmd.Stdin = pr
        cmdReader, _  :=  sf.cmd.StdoutPipe()   
        // sf.outputReader = bufio.NewScanner(cmdReader) 
        sf.outputReader = bufio.NewReader(cmdReader) 

        go sf.cmd.Run()
   } 
}

func (sf *ServerlessFunction) Stop(){

   if sf.started {
        sf.started = false
        syscall.Kill(-1*sf.cmd.Process.Pid, syscall.SIGKILL)
   } 
}

func (sf *ServerlessFunction) Trigger (event []byte) (string) {
    
    if !sf.started {
        sf.Start()
    }
    
    // send a line of data
    sf.input.Write(event)
    
    // read multiple lines of data from the function process
    var lines bytes.Buffer
    for {

        line, _, err := sf.outputReader.ReadLine()

        //will break at an empty line 
        if err != nil || len(line) == 0 {
            break
        }
   
        lines.WriteString(string(line))
        lines.WriteString("\n\r")
    }

    return lines.String()
}

type FetchFunctionPackage interface{
   FetchFunction (functionId string) (*ServerlessFunction, bool) 
} 


type ServerlessFunctionManager struct{
    functionStore map[string]*ServerlessFunction
    fetchFunction FetchFunctionPackage
}

func (mgr *ServerlessFunctionManager) SetFunctionPackageManager(gf FetchFunctionPackage){
    mgr.fetchFunction = gf
} 

func (mgr *ServerlessFunctionManager) AddFunction (sf *ServerlessFunction) {
    // sf.Start()
    if mgr.functionStore[sf.id] == nil || mgr.functionStore[sf.id] != sf {
        mgr.functionStore[sf.id] = sf
    }
}

func (mgr *ServerlessFunctionManager) CreateFunction (functionId string, data []byte, sfcommand string, sfarguments []string)(*ServerlessFunction, bool) {

    if len(functionId) == 0 || len(data) == 0 {
        return nil, false
    }

    // default value
    if len(sfcommand) == 0 {
        sfcommand = "node"
    } 

    if nil == sfarguments ||  len(sfarguments) == 0 {
        sfarguments = []string{"index.js"}
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

    err := ioutil.WriteFile(dest + "/index.js", data, 0644)
    if err != nil {
        return nil, false
    }
 
    sf := &ServerlessFunction{
        id: functionId,
        input: nil,
        outputReader: nil,
        cmd: nil,
        started: false,
        command: sfcommand,
        args: sfarguments,
    }
 
    mgr.functionStore[sf.id] = sf
    return sf, true
}
func (mgr *ServerlessFunctionManager) GetFunction (functionId string) (*ServerlessFunction, bool){
    sf, ok := mgr.functionStore[functionId]
    return sf, ok
}

func (mgr *ServerlessFunctionManager) RemoveFunction (functionId string) (bool){
    sf, ok := mgr.functionStore[functionId]

    if ok {
       sf.Stop()
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

func getManifest(baseUrl string, funcName string, funcTag string) ([]byte, error) {
    var manifestUrl = baseUrl + funcName + "/manifests/" + funcTag

	resp, err := http.Get(manifestUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
    }
    
    return body, err
}

func getLayer(baseUrl string, funcName string, blobSum string) ([]byte, error) {
	layerUrl := baseUrl + funcName + "/blobs/" + blobSum
	resp, err := http.Get(layerUrl)
	if err != nil {
		return nil, err
    }
    
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//This is just an implementation of getting function package. we need abstract this out into an interface 
func  retrieveFunction(baseUrl string, functionId string) (string, error) {
    temp := strings.Split(functionId, ":")
    funcName := ""
    funcTag := "latest"
    if len(temp) > 1 {
        l := len(temp)
        funcTag = temp[l - 1]
        funcName = temp[l - 2]
    }

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

    body, err := getManifest(baseUrl, funcName, funcTag)
    if err != nil {
		return "", err
    }
    
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return "",  err
	}
	m := f.(map[string]interface{})

	fsLayers := m["fsLayers"].([]interface{})
    um := fsLayers[0].(map[string]interface{})
    blobSum := um["blobSum"].(string)
    layer, err := getLayer(baseUrl,funcName, blobSum)
    if err != nil {
        return "", err
    }

    if len(layer) == 0 {
        return "", errors.New("Zero size function fetched.")
    }
    target := "layer.tar"
    err = ioutil.WriteFile(target, layer, 0644)
    if err != nil {
        return "", err
    }
    cmd := exec.Command("tar", "-xvf", target, "-C", dest)
    err = cmd.Run()
    if err != nil {
        return "", err
    }

    exec.Command("rm", target).Run()

    
    return  dest, nil		
}

// Default implementation of the GetFunctionPackage interface
// At this moment, we use the tenantId to hold baseUrl for testing purpose. It can be configured in the future
func (mgr *ServerlessFunctionManager) FetchFunction (functionId string) (*ServerlessFunction, bool) {

    if mgr.fetchFunction == nil {
        return nil, false
    }

    return mgr.fetchFunction.FetchFunction(functionId)
}


// default implementation of interface GetFunctionPackage
type ServerlessFunctionPackageManager struct {
    functionStoreUrl string
}

// inject the store url 
func (pm *ServerlessFunctionPackageManager) SetFunctionStoreUrl (baseUrl string){
    pm.functionStoreUrl = baseUrl
} 



func (pm *ServerlessFunctionPackageManager) FetchFunction (functionId string)(*ServerlessFunction, bool) {

    funcpath, err := retrieveFunction(pm.functionStoreUrl, functionId)
    if err != nil {
        return nil, false
    }

    var dest = os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }

    // dest = dest + "/rtsp/nodejs/bin/lambda-run"

    sf := &ServerlessFunction{
        id: functionId,
        input: nil,
        outputReader: nil,
        cmd: nil,
        started: false,
        command: "node",
        args: []string{funcpath + "/index.js"},
    }

    return sf, true
} 

// Singleton 
var instance *ServerlessFunctionManager
var once sync.Once

func GetFunctionManager() (*ServerlessFunctionManager) {

    once.Do( func() {
            instance = &ServerlessFunctionManager {
            fetchFunction: nil, //need be injected later
            functionStore: make(map[string]*ServerlessFunction),
        }
    })

    return instance
}


