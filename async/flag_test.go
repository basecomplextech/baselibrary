// Copyright 2022 Ivan Korobkov. All rights reserved.

package async

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
