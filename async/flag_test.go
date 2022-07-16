package async

import "testing"

func TestFlag__should_signal_and_reset_flag(t *testing.T) {
	f := NewFlag()
	select {
	case <-f.Wait():
		t.Fatal("flag should not be set")
	default:
	}

	f.Signal()
	select {
	case <-f.Wait():
	default:
		t.Fatal("flag should be set")
	}

	f.Reset()
	select {
	case <-f.Wait():
		t.Fatal("flag should not be set")
	default:
	}

	f.Signal()
	select {
	case <-f.Wait():
	default:
		t.Fatal("flag should be set")
	}
}
