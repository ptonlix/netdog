package dogmail

import (
	"testing"
)

func TestSend(t *testing.T) {
	dogmail := NewDogmail()
	err := dogmail.Send("test")
	if err != nil {
		t.Error("Mail Send error", err)
		return
	}
	t.Log("success")
}
