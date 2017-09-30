package serverlessfunction

import (
    "bufio"
    "bytes"
    "io"     
    "io/ioutil"
    "syscall"
    "os/exec"
    "strings"
    "os"
    "errors"
)

const EndOfFunction = "--{{{|}}}--"
type ServerlessFunction struct{
    id string
    input io.Writer
    outputReader *bufio.Reader
    cmd  *exec.Cmd
    command string
    args []string
    started bool
}

func LoadFunction(functionId string) *ServerlessFunction {

    dest := os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }

    sfCommand := "node"
    lambda := os.Getenv("RUNTIME_LAMBDA")
    if len(lambda) == 0 {
        lambda = "/var/runtime/bin/lambda-run"
    }
   
    dest += "/func/" + functionId
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        return nil
    } 
    dest += "/index.js"
    if _, err := os.Stat(dest); os.IsNotExist(err) {
        return nil
    } 

    sfArguments := []string{lambda, dest}
    ff := &ServerlessFunction{
        id: functionId,
        input: nil,
        outputReader: nil,
        cmd: nil,
        started: false,
        command: sfCommand,
        args: sfArguments,
    }

    return ff
}

// Use default value for sfCommand and sfArguments
func CreateServerlessFunction (functionId string, data []byte, sfCommand string, sfArguments []string) (*ServerlessFunction, error){

    if len(functionId) == 0 {
        return nil, errors.New("Invalid function id.")
    }

    if  len(data) == 0 {
        return nil, errors.New("Invalid function data.")
    }
    // default value
    if len(sfCommand) == 0 {
        sfCommand = "node"
    } 
 
    // Save the function js file
    dest := os.Getenv("RUNTIME_ROOT")
    if len(dest) == 0 {
        dest = "/var/runtime"
    }

    dest = dest + "/func/" + functionId
    os.RemoveAll(dest)
    os.Mkdir(dest, os.ModeDir)

    dest = dest + "/index.js"
    err := ioutil.WriteFile(dest, data, 0644)
    if err != nil {
        return nil, err
    }

    if nil == sfArguments ||  len(sfArguments) == 0 {
        lambda := os.Getenv("RUNTIME_LAMBDA")
        if len(lambda) == 0 {
            lambda = "/var/runtime/bin/lambda-run"
        }
        sfArguments = []string{lambda, dest}
    }

    ff := &ServerlessFunction{
        id: functionId,
        input: nil,
        outputReader: nil,
        cmd: nil,
        started: false,
        command: sfCommand,
        args: sfArguments,
    }

    return ff, nil
}

func (sf *ServerlessFunction) Name() string{
   return sf.id
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
   // A lock is need to protect sf.started operation
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
        if err != nil {
            lines.WriteString(err.Error())
            lines.WriteString("\n\r")
            break
        }

        newLine := string(line)

        if strings.HasPrefix(newLine, EndOfFunction) {
            break;
        }
       
        lines.WriteString(newLine)
        lines.WriteString("\n\r")
      
    }

    return lines.String()
}
