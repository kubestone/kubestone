/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"bytes"
	"log"
	"os/exec"

	"github.com/firepear/qsplit"
)

func run(command string) (stdout, stderr string, err error) {
	commandArray := qsplit.ToStrings([]byte(command))
	cmd := exec.Command(commandArray[0], commandArray[1:]...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdOut, &stdErr
	if err = cmd.Run(); err != nil {
		log.Printf("Error during execution of `%v`\nerr: %v\nstdout: %v\nstderr: %v\n",
			command, err, stdOut.String(), stdErr.String())
		return "", "", err
	}

	return stdOut.String(), stdErr.String(), nil
}
