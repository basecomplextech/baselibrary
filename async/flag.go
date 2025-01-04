// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal/flag"

// Flag is a read-only boolean flag that can be waited on until set.
//
// Example:
//
//	serving := async.UnsetFlag()
//
//	func handle(ctx Context, req *request) {
//		if !serving.Get() { // just to show Get in example
//			select {
//			case <-ctx.Wait():
//				return ctx.Status()
//			case <-serving.Wait():
//			}
//		}
//
//		// ... handle request
//	}
type Flag = flag.Flag

// MutFlag is a routine-safe boolean flag that can be set, reset, and waited on until set.
//
// Example:
//
//	serving := async.UnsetFlag()
//
//	func serve() {
//		s.serving.Set()
//		defer s.serving.Unset()
//
//		// ... start server ...
//	}
//
//	func handle(ctx Context, req *request) {
//		select {
//		case <-ctx.Wait():
//			return ctx.Status()
//		case <-serving.Wait():
//		}
//
//		// ... handle request
//	}
type MutFlag = flag.MutFlag

// SetFlag returns a new set flag.
func SetFlag() MutFlag {
	return flag.SetFlag()
}

// UnsetFlag returns a new unset flag.
func UnsetFlag() MutFlag {
	return flag.UnsetFlag()
}

// ReverseFlag returns a new flag which reverses the original one.
func ReverseFlag(f Flag) Flag {
	return flag.ReverseFlag(f)
}
