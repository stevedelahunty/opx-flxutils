package keepalive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sysd"
	"time"
	"utils/ipcutils"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type SysdClient struct {
	ipcutils.IPCClientBase
	ClientHdl *sysd.SYSDServicesClient
}

type KeepAlive struct {
	sysdClient SysdClient
	name       string
	status     int32
}

const (
	KA_ACTIVE = 1
)

func InitKeepAlive(name string, paramsFile string) {
	var clientsList []ClientJson

	ka := new(KeepAlive)
	ka.name = name
	ka.status = KA_ACTIVE

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		fmt.Println("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		fmt.Println("Error in Unmarshalling Json")
		return
	}

	for _, client := range clientsList {
		if client.Name == "sysd" {
			fmt.Println("found sysd at port", client.Port)
			ka.sysdClient.Address = "localhost:" + strconv.Itoa(client.Port)
			ka.sysdClient.TTransport, ka.sysdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(ka.sysdClient.Address)
			if err != nil {
				fmt.Println("Failed to connect to sysd, retrying until connection is successful")
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					ka.sysdClient.TTransport, ka.sysdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(ka.sysdClient.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 10) == 0 {
						fmt.Println("Still can't connect to sysd, retrying...")
					}
				}
			}
		}
	}
	fmt.Println("Initialized KA for ", ka.name, " status ", ka.status)
	retryTimer := time.NewTicker(time.Second * 5)
	for t := range retryTimer.C {
		_ = t
		ka.sysdClient.ClientHdl.PeriodicKeepAlive(ka.name)
	}
	return
}

func (ka *KeepAlive) SendPeriodicKeepAlive() error {
	fmt.Println("SendPeriodicKeepAlive ", ka.name, ka.status)

	return nil
}

func (ka *KeepAlive) SetMyStatus(status int32) error {
	return nil
}
