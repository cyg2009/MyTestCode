package serverlessfunction

import (
    "bufio"
    "bytes"
    "io"      
    "syscall"
    "os/exec"
    "fmt"
    "strings"
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

func NewServerlessFunction(functionId string, sfcommand string, sfarguments []string) *ServerlessFunction{
     return &ServerlessFunction{
        id: functionId,
        input: nil,
        outputReader: nil,
        cmd: nil,
        started: false,
        command: sfcommand,
        args: sfarguments,
    }
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
