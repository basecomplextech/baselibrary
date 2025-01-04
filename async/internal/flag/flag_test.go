// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package flag

import "testing"

func TestFlag__should_set_and_reset_flag(t *testing.T) {
	f := UnsetFlag()
	select {
	case <-f.Wait():
		t.Fatal("flag should not be set")
	default:
	}

	f.Set()
	select {
	case <-f.Wait():
	default:
		t.Fatal("flag should be set")
	}

	f.Unset()
	select {
	case <-f.Wait():
		t.Fatal("flag should not be set")
	default:
	}

	f.Set()
	select {
	case <-f.Wait():
	default:
		t.Fatal("flag should be set")
	}
}

// Reverse

func TestReverseFlag__should_reverse_source_flag(t *testing.T) {
	f := UnsetFlag()
	r := ReverseFlag(f)

	select {
	case <-r.Wait():
	default:
		t.Fatal("reverse flag should be set")
	}

	f.Set()
	select {
	case <-r.Wait():
		t.Fatal("reverse flag should be unset")
	default:
	}

	f.Unset()
	select {
	case <-r.Wait():
	default:
		t.Fatal("reverse flag should be set")
	}

	f.Set()
	select {
	case <-r.Wait():
		t.Fatal("reverse flag should be unset")
	default:
	}
}
