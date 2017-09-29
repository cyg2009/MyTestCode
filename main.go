package main

import (
    "os"
    "fmt"
    "net"
    "github.com/cyg2009/MyTestCode/pkg/prservice"
    "net/http"
)

const (
    DEFAULT_ETH_NAME = "eth0"
    DEFAULT_SERVER_ADDR = "127.0.0.1"
    DEFAULT_SERVER_PORT = 28903
)

func getLocalAddress() string {
    ethname := os.ExpandEnv("$SERVERLESS_ETH_NAME")
    if ethname == "" {
        ethname = DEFAULT_ETH_NAME
    }

    iface, err := net.InterfaceByName(ethname)
    if err != nil {
        //log.LOGGER.Errorf(err, "get interface '%s' failed.", ethname)
        return DEFAULT_SERVER_ADDR
    }

    addrs, err := iface.Addrs()
    if err != nil || len(addrs) == 0 {
        //log.LOGGER.Errorf(err, "get addrs of interface '%s' failed.", iface.Name)
        return DEFAULT_SERVER_ADDR
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
            return ipnet.IP.String()
        }
    }

    return DEFAULT_SERVER_ADDR
}

func getServerAddress() string {

    /*hostIP := appconfig.DefaultString("httpaddr", "")
    if hostIP == "" || hostIP == "127.0.0.1" {
        hostIP = getLocalAddress()
    }
    port := appconfig.DefaultInt("httpport", DEFAULT_SERVER_PORT)*/


    return fmt.Sprintf("%s:%d", getLocalAddress(), DEFAULT_SERVER_PORT)
}

func main() {
    fmt.Println("Starting process router...")

    var handler http.Handler = prservice.GetPrserviceHttpHandler()

    listenAddress := getServerAddress()
    fmt.Println("Listen at " + listenAddress)
    err := http.ListenAndServe(listenAddress, handler)
    if err != nil {
        fmt.Errorf("listen on address %s failed.", listenAddress)
        panic(err)
    }

    os.Exit(0)
}

