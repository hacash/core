package test

import (
	"fmt"
	"github.com/hacash/core/inicnf"
	"os"
	"testing"
)

func Test_t1(t *testing.T) {

	testcnffn := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/core/test/config.ini"

	cnf, e := inicnf.Load(testcnffn)
	if e != nil {
		fmt.Println(e)
		return
	}

	data_dir := cnf.Section("").Key("data_dir").MustString("~/.hacash_mainnet")

	fmt.Println(data_dir)

	fmt.Println(os.Getenv("HOME"))

}
