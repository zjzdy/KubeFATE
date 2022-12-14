/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Flags: []cli.Flag{},
		Usage: "Show kubefate command line and service version",
		Action: func(c *cli.Context) error {
			fmt.Printf("* kubefate commandLine version=%s\r\n", api.ServiceVersion)
			serviceVersion, err := GetServiceVersion()
			if err != nil {
				return fmt.Errorf("* kubefate service connection error, %s", err.Error())
			}
			fmt.Printf("* kubefate service version=%s\r\n", serviceVersion)
			return nil
		},
	}
}

func GetServiceVersion() (string, error) {
	r := &Request{
		Type: "GET",
		Path: "version",
		Body: nil,
	}

	serviceUrl := viper.GetString("serviceurl")
	apiVersion := api.APIVersion + "/"
	if serviceUrl == "" {
		serviceUrl = "localhost:8080/"
	}
	safeconnect := viper.GetString("safeconnect")
	scheme := "http://"
	if safeconnect == "true" {
		scheme = "https://"
	}
	Url := scheme + serviceUrl + "/" + apiVersion + r.Path
	body := bytes.NewReader(r.Body)
	log.Debug().Str("Type", r.Type).Str("url", Url).Msg("Request")
	request, err := http.NewRequest(r.Type, Url, body)
	if err != nil {
		return "", err
	}

	token, err := getToken()
	if err != nil {
		return "", err
	}
	Authorization := fmt.Sprintf("Bearer %s", token)

	request.Header.Add("Authorization", Authorization)
	request.Header.Add("user-agent", "kubefate")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp.StatusCode=%d, error: %s", resp.StatusCode, respBody)
	}

	type VersionResultMsg struct {
		Msg     string
		Version string
	}

	VersionResult := new(VersionResultMsg)

	err = json.Unmarshal(respBody, &VersionResult)
	if err != nil {
		return "", err
	}

	log.Debug().Int("Code", resp.StatusCode).Bytes("Body", respBody).Msg("ok")
	return VersionResult.Version, err
}
