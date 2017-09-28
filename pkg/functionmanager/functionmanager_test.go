package functiomanager

import (
	"testing"
	"fmt"
	"os"
)

func setup(){
    fmt.Println("Set up testing for functionmanager...") 
    os.Setenv("RUNTIME_ROOT", "/work/src/processrouter/runtime")
}

func shutdown(){
	fmt.Println("Tear down testing for functionmanager...")
}

func TestMain(m *testing.M) {
    setup()
    code := m.Run() 
    shutdown()
    os.Exit(code)
}

type FetchFunctionPackageMock struct {

} 

func (gf *FetchFunctionPackageMock) FetchFunction (functionId string) (*ServerlessFunction, bool) {
    return nil, true
} 

// testing of one function but with multiple sequential invocations 
func TestFunctionManager_GetFunctionMock(t *testing.T){

   fmg := GetFunctionManager()
   fmg.SetFunctionPackageManager(&FetchFunctionPackageMock{})
   functionId := "Test3:1.0.0"
   
   _, ok := fmg.FetchFunction(functionId)

   if ok == false {
	   t.Errorf("Fail to get function.")
   }
}

func Test_getManifest(t *testing.T){

   baseUrl := "http://10.21.119.117:5000/v2/serverless/"
   body, err := getManifest(baseUrl, "f5", "1.0.0")
  
   if err != nil || len(body) < 1 {
	   t.Errorf("Fail to get function.")
   }
}

func TestFunctionManager_getFunction(t *testing.T){

    fgr := GetFunctionManager()
    pkgMgr := &ServerlessFunctionPackageManager{}
    baseUrl := "http://10.21.119.117:5000/v2/serverless/"
    pkgMgr.SetFunctionStoreUrl(baseUrl)
    fgr.SetFunctionPackageManager(pkgMgr)

    sf, ok := fgr.FetchFunction("f3:1.0.0")
  
    if ok == false || sf == nil {
	   t.Errorf("Fail to get function.")
    }
    fgr.AddFunction(sf)
    ret := fgr.GetAllFunctionsJSON()

    t.Log(ret)
}