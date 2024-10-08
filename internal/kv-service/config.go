package kvservice

import (
	"time"
)

var clientId int

func StartServers(peerAddresses []string, demoStoragePath string) *Clerk {
	kvServers := []*Server{}

	for i := range peerAddresses {
		kvServer := StartServer(peerAddresses, i, demoStoragePath)
		kvServers = append(kvServers, kvServer)
	}

	time.Sleep(2 * time.Second)

	clerk := MakeClerk(kvServers, int32(clientId))
	clientId++

	return clerk
}

func init() {
	clientId++
}
