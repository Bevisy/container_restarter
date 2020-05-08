package main

import (
    "flag"
    "log"
    "os"
    "time"

    "github.com/bevisy/container_restart/util"
)

var (
    container  string
    id         string
    domainName string
    interval   time.Duration
    host       string
    version    string
)

func init() {
    flag.StringVar(&container, "c", "", "Name of the container to restart")
    flag.StringVar(&id, "id", "", "ID of the container to restart")
    flag.StringVar(&domainName, "d", "", "Domain name to watch IP change")
    flag.DurationVar(&interval, "t", time.Duration(10*time.Second),
        "Time interval to check IP change on domain name")
    flag.StringVar(&host, "h", "unix:///var/run/docker.sock",
        "Docker server host")
    flag.StringVar(&version, "v", "1.40", "Docker API version")
}

func main() {
    //命令行参数初始化
    flag.Parse()

    //参数判断: container和id 二选一, domainName 必填
    if (container == "" && id == "") || domainName == "" {
        flag.Usage()
        os.Exit(1)
    }

    //根据传入参数，新建docker客户端
    cli, err := util.MyDockerClient(host, version)
    if err != nil {
        log.Fatal(err)
    }

    //根据传入的container参数，修改执行函数，用于调试程序
    var act func() error
    if container == "debug" {
        act = func() error { return nil }
    } else {
        act = util.fRestartC(cli, container)
    }

    //Type: channel
    //用于存取解析的IP
    newIP := util.resolver(domainName, interval)

    for {
        // <-newIP() is blocked till *domainName resolved IP changes
        log.Printf("IP address changed to %s", <-newIP())
        act()
    }

}
