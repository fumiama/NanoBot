package nano

import "testing"

func TestWriteHTTPQueryIfNotNil(t *testing.T) {
	expstr := "https://api.sgroup.qq.com/testapi?b=1&d=0.5"
	str := WriteHTTPQueryIfNotNil(StandardAPI+"/testapi", "a", 0, "b", 1, "c", "", "d", 0.5)
	if str != expstr {
		t.Fatal("expected", expstr, "but got", str)
	}
}
