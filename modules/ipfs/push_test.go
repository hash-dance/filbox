package ipfs

import (
	"regexp"
	"testing"
	"time"
)

func TestSring(t *testing.T) {
	str := "adminer_4-fastcgi.tar"
	t.Log(regexp.MustCompile(`^[a-zA-Z]+`).FindString(str))
	t.Logf("%s", str)

	t1 := time.Now()
	t2 := t1.Add(-time.Hour * 1)
	t.Log(t1)
	t.Log(t2)
	t.Log(t1.After(t2))
}
