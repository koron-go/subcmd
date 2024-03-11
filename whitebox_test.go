package subcmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStripExeExt(t *testing.T) {
	for i, c := range []struct{ in, want string }{
		{"foo.exe", "foo"},
		{"foo", "foo"},
		{"bar.txt", "bar.txt"},
		{"bar.txt.exe", "bar.txt"},
		{"subcmd.test.exe", "subcmd.test"},
		{"subcmd.test", "subcmd.test"},
	} {
		got := stripExeExt(c.in)
		if d := cmp.Diff(c.want, got); d != "" {
			t.Errorf("unexpected result at #%d: -want +got\n%s", i, d)
		}
	}
}
