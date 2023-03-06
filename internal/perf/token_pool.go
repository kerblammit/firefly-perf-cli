// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package perf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/firefly/pkg/fftypes"
	log "github.com/sirupsen/logrus"
)

func (pr *perfRunner) CreateTokenPool() error {
	log.Infof("Creating Token Pool: %s", pr.poolName)
	var config fftypes.JSONObject = make(map[string]interface{})
	body := fftypes.TokenPool{
		Connector: pr.cfg.TokenOptions.TokenPoolConnectorName,
		Name:      pr.poolName,
		Type:      getTokenTypeEnum(pr.cfg.TokenOptions.TokenType),
		Config:    config,
	}

	if pr.cfg.TokenOptions.Config.PoolAddress != "" {
		config["address"] = pr.cfg.TokenOptions.Config.PoolAddress
	}

	if pr.cfg.TokenOptions.Config.PoolBlockNumber != "" {
		config["blockNumber"] = pr.cfg.TokenOptions.Config.PoolBlockNumber
	}

	payload, err := json.Marshal(body)
	fmt.Println("Create token body:")
	fmt.Println(string(payload))

	fmt.Println("Creating token pool at")
	fmt.Println(fmt.Sprintf("/%s/api/v1/namespaces/%s/tokens/pools?confirm=true", pr.cfg.APIPrefix, pr.cfg.FFNamespace))
	fmt.Println(body)
	res, err := pr.client.R().
		SetBody(&body).
		Post(fmt.Sprintf("/%s/api/v1/namespaces/%s/tokens/pools?confirm=true", pr.cfg.APIPrefix, pr.cfg.FFNamespace))

	if err != nil || !res.IsSuccess() {
		fmt.Println(res.Body())
		fmt.Println(res.String())
		fmt.Println(res.Status())
		fmt.Println(res.RawResponse)
		fmt.Printf("Pool err %v", res.StatusCode())
		fmt.Println(err)
		return errors.New("Failed to create token pool")
	}
	fmt.Println("Finished making token pool")
	return err
}

func getTokenTypeEnum(tokenType string) fftypes.FFEnum {
	if tokenType == "nonfungible" {
		return fftypes.TokenTypeNonFungible
	}
	return fftypes.TokenTypeFungible
}
