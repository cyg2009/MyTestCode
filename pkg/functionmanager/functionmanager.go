package functiomanager

import (
    "bufio"
    "bytes"
    "io"      
    "io/ioutil"
    "syscall"
    "os/exec"
    "errors"
    "sync"
    "os"
    "fmt"
    "strings"
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
        // return fmt.Sprintf("get a line of data back: %d",len(line)) 
        newLine := string(line)

        fmt.Println(newLine)
        if strings.HasPrefix(newLine, `{{{}}}`) {
            break;
        }
       
        lines.WriteString(newLine)
        lines.WriteString("\n\r")

          //will break at an empty line 
        if err != nil {
            fmt.Println("Got error message:")
               fmt.Println(err)
            break
        }
    }

    return lines.String()
}

type ServerlessFunctionManager struct{
    functionStore map[string]*ServerlessFunction
}

func (mgr *ServerlessFunctionManager) CreateFunction (functionId string, data []byte, sfcommand string, sfarguments []string)(*ServerlessFunction, bool) {

    if len(functionId) == 0 || len(data) == 0 {
        return nil, false
    }

    // default value
    if len(sfcommand) == 0 {
        sfcommand = "node"
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
   
    if nil == sfarguments ||  len(sfarguments) == 0 {
        sfarguments = []string{"lambda-run", dest}
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

// Singleton 
var instance *ServerlessFunctionManager
var once sync.Once

func GetFunctionManager() (*ServerlessFunctionManager) {

    once.Do( func() {
            instance = &ServerlessFunctionManager {
            functionStore: make(map[string]*ServerlessFunction),
        }
    })

    return instance
}




