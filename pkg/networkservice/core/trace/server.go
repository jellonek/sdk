// Copyright (c) 2020 Cisco Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/networkservicemesh/api/pkg/api/connection"
	"github.com/networkservicemesh/api/pkg/api/networkservice"

	"github.com/networkservicemesh/sdk/pkg/tools/spanhelper"
	"github.com/networkservicemesh/sdk/pkg/tools/typeutils"
)

type traceServer struct {
	traced networkservice.NetworkServiceServer
}

// NewNetworkServiceServer - wraps tracing around the supplied traced
func NewNetworkServiceServer(traced networkservice.NetworkServiceServer) networkservice.NetworkServiceServer {
	return &traceServer{traced: traced}
}

func (t *traceServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*connection.Connection, error) {
	// Create a new span
	span := spanhelper.FromContext(ctx, fmt.Sprintf("%s.Request", typeutils.GetTypeName(t.traced)))
	defer span.Finish()

	// Make sure we log to span

	ctx = withLog(span.Context(), span.Logger())

	span.LogObject("request", request)

	// Actually call the next
	rv, err := t.traced.Request(ctx, request)

	if err != nil {
		span.LogError(err)
		return nil, err
	}
	span.LogObject("response", rv)
	return rv, err
}

func (t *traceServer) Close(ctx context.Context, conn *connection.Connection) (*empty.Empty, error) {
	// Create a new span
	span := spanhelper.FromContext(ctx, fmt.Sprintf("%s.Close", typeutils.GetTypeName(t.traced)))
	defer span.Finish()
	// Make sure we log to span
	ctx = withLog(span.Context(), span.Logger())

	span.LogObject("request", conn)
	rv, err := t.traced.Close(ctx, conn)

	if err != nil {
		span.LogError(err)
		return nil, err
	}
	span.LogObject("response", rv)
	return rv, err
}
