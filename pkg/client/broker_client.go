/*
 * Copyright 2018 The Service Manager Authors
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
package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/Peripli/service-manager/pkg/types"
	"github.com/Peripli/service-manager/pkg/util"
	"net/http"
)

type BrokerClient struct {
	transport                *BrokerTransport
	broker                   *types.ServiceBroker
	requestWithClientHandler util.DoRequestWithClientFunc
	requestHandlerDecorated  util.DoRequestFunc
}

func NewBrokerClient(broker *types.ServiceBroker, requestHandler util.DoRequestWithClientFunc) (*BrokerClient, error) {
	tlsConfig, err := broker.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to get client for broker %s: %v", broker.Name, err)
	}

	if requestHandler == nil {
		return nil, errors.New("a request handler func is required")
	}

	bc := &BrokerClient{}
	bc.transport = NewBrokerTransport(tlsConfig)
	bc.broker = broker
	bc.requestWithClientHandler = requestHandler
	bc.requestHandlerDecorated = bc.authAndTlsDecorator()
	return bc, nil
}

func (bc *BrokerClient) addBasicAuth(req *http.Request) *BrokerClient {
	req.SetBasicAuth(bc.broker.Credentials.Basic.Username, bc.broker.Credentials.Basic.Password)
	return bc
}

func (bc *BrokerClient) authAndTlsDecorator() util.DoRequestFunc {
	return func(req *http.Request) (*http.Response, error) {
		client := http.DefaultClient

		if bc.broker.Credentials.Basic != nil && bc.broker.Credentials.Basic.Username != "" && bc.broker.Credentials.Basic.Password != "" {
			bc.addBasicAuth(req)
		}

		useDedicatedClient, transportWithTLS := bc.transport.GetTransportWithTLS()

		if useDedicatedClient {
			client = &http.Client{}
			client.Transport = transportWithTLS
			return bc.requestWithClientHandler(req, client)
		}

		return bc.requestWithClientHandler(req, client)
	}
}

func (bc *BrokerClient) SendRequest(ctx context.Context, method, url string, params map[string]string, body interface{}, headers map[string]string) (*http.Response, error) {
	return util.SendRequestWithHeaders(ctx, bc.requestHandlerDecorated, method, url, params, body, headers)
}
