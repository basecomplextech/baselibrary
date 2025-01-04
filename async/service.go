// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/service"
	"github.com/basecomplextech/baselibrary/status"
)

// Service is a service which can be started and stopped.
type Service = service.Service

// NewService returns a new stopped service.
func NewService(fn func(ctx Context) status.Status) Service {
	return service.New(fn)
}
