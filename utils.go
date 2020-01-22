package common

import (
    "github.com/gin-gonic/gin"
    "net"
    "strconv"
)

func LocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func InitKey(c *gin.Context) {
    if len(c.Keys) == 0 {
        c.Keys = make(map[string]interface{})
        c.Keys["auth"] = &AuthGlobal{}
        c.Keys["code"] = 0
    }
}

func GetDBIndex(taskId int) (index string) {
    indexInt := taskId % 11
    index = strconv.Itoa(indexInt)
    return index
}
