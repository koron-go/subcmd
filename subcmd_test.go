package subcmd_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron-go/subcmd"
)

func TestCommand(t *testing.T) {
	var called bool
	cmd := subcmd.DefineCommand("foo", t.Name(), func(context.Context, []string) error {
		called = true
		return nil
	})
	err := subcmd.Run(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("command func is not called")
	}
}

func TestCommandNil(t *testing.T) {
	cmd := subcmd.DefineCommand("foo", t.Name(), nil)
	err := subcmd.Run(cmd)
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	if d := cmp.Diff("no function declared for command: foo", err.Error()); d != "" {
		t.Errorf("error unmatch: -want +got\n%s", d)
	}
}

func TestSet(t *testing.T) {
	var (
		gotNames []string
		gotArgs  []string
	)
	record := func(ctx context.Context, args []string) error {
		gotNames = subcmd.Names(ctx)
		gotArgs = args
		return nil
	}

	set := subcmd.DefineSet("set", "",
		subcmd.DefineSet("user", "",
			subcmd.DefineCommand("list", "", record),
			subcmd.DefineCommand("add", "", record),
			subcmd.DefineCommand("delete", "", record),
		),
		subcmd.DefineSet("post", "",
			subcmd.DefineCommand("list", "", record),
			subcmd.DefineCommand("add", "", record),
			subcmd.DefineCommand("delete", "", record),
		),
	)

	for i, c := range []struct {
		args      []string
		wantNames []string
		wantArgs  []string
	}{
		{
			[]string{"user", "list"},
			[]string{"set", "user", "list"},
			[]string{},
		},
		{
			[]string{"user", "add", "-email", "foobar@example.com"},
			[]string{"set", "user", "add"},
			[]string{"-email", "foobar@example.com"},
		},
		{
			[]string{"user", "delete", "-id", "123"},
			[]string{"set", "user", "delete"},
			[]string{"-id", "123"},
		},
		{
			[]string{"post", "list"},
			[]string{"set", "post", "list"},
			[]string{},
		},
		{
			[]string{"post", "add", "-title", "Brown fox..."},
			[]string{"set", "post", "add"},
			[]string{"-title", "Brown fox..."},
		},
		{
			[]string{"post", "delete", "-id", "ABC"},
			[]string{"set", "post", "delete"},
			[]string{"-id", "ABC"},
		},
	} {
		err := subcmd.Run(set, c.args...)
		if err != nil {
			t.Fatalf("failed for case#%d (%+v): %s", i, c, err)
			continue
		}
		if d := cmp.Diff(c.wantNames, gotNames); d != "" {
			t.Errorf("unexpected names on #%d: -want +got\n%s", i, d)
		}
		if d := cmp.Diff(c.wantArgs, gotArgs); d != "" {
			t.Errorf("unexpected args on #%d: -want +got\n%s", i, d)
		}
	}
}

func TestSetFails(t *testing.T) {
	rootSet := subcmd.DefineSet("fail", "",
		subcmd.DefineCommand("list", "list all entries", nil),
		subcmd.DefineCommand("add", "add a new entry", nil),
		subcmd.DefineCommand("delete", "delete an entry", nil),
		subcmd.DefineSet("item", "operate items"),
	)
	for i, c := range []struct {
		args []string
		want string
	}{
		{[]string{}, `no commands selected.

Available sub-commands are:

	list        list all entries
	add         add a new entry
	delete      delete an entry
	item        operate items`},
		{[]string{"foo"}, `command not found.

Available sub-commands are:

	list        list all entries
	add         add a new entry
	delete      delete an entry
	item        operate items`},
	} {
		err := subcmd.Run(rootSet, c.args...)
		if err == nil {
			t.Fatalf("unexpected succeed at #%d %+v", i, c)
		}
		got := err.Error()
		if d := cmp.Diff(c.want, got); d != "" {
			t.Errorf("unexpected error at #%d: -want +got\n%s", i, d)
		}
	}
}

func TestAutoWidth(t *testing.T) {
	rootSet := subcmd.DefineSet("faillong", "",
		subcmd.DefineCommand("verylongname", "long name command", nil),
		subcmd.DefineCommand("short", "short name command", nil),
	)
	err := subcmd.Run(rootSet)
	if err == nil {
		t.Fatal("unexpected succeed")
	}
	want := `no commands selected.

Available sub-commands are:

	verylongname    long name command
	short           short name command`
	got := err.Error()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("unexpected error: -want +got\n%s", d)
	}
}

func TestRootSet(t *testing.T) {
	rootSet := subcmd.DefineRootSet()
	if d := cmp.Diff("subcmd.test", rootSet.Name()); d != "" {
		t.Errorf("unexpected name: -want +got\n%s", d)
	}
}

func TestFlagSet(t *testing.T) {
	t.Run("no name", func(t *testing.T) {
		ctx := context.Background()
		fs := subcmd.FlagSet(ctx)
		got := fs.Name()
		if got != "" {
			t.Errorf("wrong name of FlagSet: got=%q", got)
		}
	})

	t.Run("with names", func(t *testing.T) {
		var gotName string
		set := subcmd.DefineSet("first", "",
			subcmd.DefineSet("second", "",
				subcmd.DefineCommand("third", "", func(ctx context.Context, args []string) error {
					fs := subcmd.FlagSet(ctx)
					gotName = fs.Name()
					return nil
				}),
			),
		)
		err := subcmd.Run(set, "second", "third")
		if err != nil {
			t.Fatalf("failed: %s", err)
		}
		if gotName != "first second third" {
			t.Errorf("wrong name of FlagSet: got=%q", gotName)
		}
	})
}
