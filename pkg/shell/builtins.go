package shell

import (
	"context"
	"io/ioutil"
	"path/filepath"
)

func Echo(ctx context.Context, args []string, sh *Shell) (int, error) {
	for _, arg := range args[1:] {
		sh.out.Write([]byte(arg))
		sh.out.Write([]byte(" "))
	}
	sh.out.Write([]byte("\n"))
	return 0, nil
}

// CD changes directories
func CD(ctx context.Context, args []string, sh *Shell) (int, error) {
	var dir string
	if len(args) < 2 {
		dir = "" // Empty string is this process's cwd
	} else if filepath.IsAbs(args[1]) {
		dir = filepath.Clean(args[1])
	} else {
		dir = filepath.Join(sh.Cwd(), args[1])
	}
	// todo check dir exists and is a directory
	sh.cwd = dir
	return 0, nil
}

func LS(ctx context.Context, args []string, sh *Shell) (int, error) {
	if len(args) < 2 {
		args = []string{"ls", sh.Cwd()}
	}
	for _, arg := range args[1:] {
		infos, err := ioutil.ReadDir(arg)
		if err != nil {
			return 1, err
		}
		for _, info := range infos {
			sh.out.Write([]byte(info.Name()))
			sh.out.Write([]byte("\n"))
		}
	}
	return 0, nil
}

func Pwd(ctx context.Context, args []string, sh *Shell) (int, error) {
	sh.out.Write([]byte(sh.Cwd()))
	sh.out.Write([]byte("\n"))
	return 0, nil
}

func Exit(ctx context.Context, args []string, sh *Shell) (int, error) {
	return 0, ExitError{}
}
