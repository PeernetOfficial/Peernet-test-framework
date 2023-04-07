package testframework

import (
	"bufio"
	"fmt"
	"github.com/PeernetOfficial/core/webapi"
	"os"
)

func AddFileTest(web *webapi.WebapiInstance) {
	file, err := os.Open("./TestFile/2023-02-08-05-06-20.mp4")
	if err != nil {
		fmt.Println(err)
	}
	_, status, err := web.Backend.UserWarehouse.CreateFile(bufio.NewReader(file), 0)
	if err != nil {
		fmt.Println(err)
	}

	web.Backend.LogError("warehouse.CreateFile", "status %d error: %v", status, err)

}

func AddFilesInNodes(Nodes []PeernetNode) {
	// For nodes provided
	// Add test file
	for i := range Nodes {
		AddFileTest(Nodes[i].NodeConfig)
	}
}
