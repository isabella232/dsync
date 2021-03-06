/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log"
	"strconv"

	"github.com/minio/dsync/v3"
)

const basePort = 50001

func main() {
	// Start four lock servers.
	for i := 0; i < 4; i++ {
		go StartLockServer(basePort + i)
	}

	// Create four lock clients.
	lockClients := make([]dsync.NetLocker, 4)
	for i := 0; i < 4; i++ {
		lockClients[i] = newReconnectRPCClient("localhost", serviceEndpointPrefix+strconv.Itoa(basePort+i))
	}

	// Initialize dsync and treat 0th index on lockClients as self node.
	ds, err := dsync.New(lockClients)
	if err != nil {
		log.Fatal("Fail to initialize dsync.", err)
	}

	// Get new distributed RWMutex on resource "Music"
	drwMutex := dsync.NewDRWMutex(context.Background(), "Music", ds)

	// Lock "music" resource.
	drwMutex.Lock("Music", "main.go:50:main()")

	// As we got writable lock on Music, do some crazy things.

	// Unlock "Music" resource.
	drwMutex.Unlock()
}
