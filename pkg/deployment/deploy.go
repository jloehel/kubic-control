// Copyright 2019 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deployment

import (
	"os"
	"io"
	"encoding/hex"
	"crypto/sha256"

        "gopkg.in/ini.v1"
        "github.com/thkukuk/kubic-control/pkg/tools"
)

func sha256sum(filePath string) (result string, err error) {
    file, err := os.Open(filePath)
    if err != nil {
        return
    }
    defer file.Close()

    hash := sha256.New()
    _, err = io.Copy(hash, file)
    if err != nil {
        return
    }

    result = hex.EncodeToString(hash.Sum(nil))
    return
}

func DeployFile(yamlName string) (bool, string) {

	success, message := tools.ExecuteCmd("kubectl", "--kubeconfig=/etc/kubernetes/admin.conf",
		"apply", "-f", yamlName)
	if success != true {
		return success, message
	}

	result, err := sha256sum(yamlName)

	cfg, err := ini.LooseLoad("/var/lib/kubic-control/k8s-yaml.conf")
	if err != nil {
		return false, "Cannot load k8s-yaml.conf: " + err.Error()
        }

	cfg.Section("").Key(yamlName).SetValue(result)
	err = cfg.SaveTo("/var/lib/kubic-control/k8s-yaml.conf")
        if err != nil {
		return false, "Cannot write k8s-yaml.conf: " + err.Error()
        }

	return true, ""
}
