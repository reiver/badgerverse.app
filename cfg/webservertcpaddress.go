package cfg

import (
	"fmt"

	"github.com/reiver/badgerverse.app/env"
)

func WebServerTCPAddress() string {
	return fmt.Sprintf(":%s", env.TcpPort)
}
