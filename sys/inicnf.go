package sys

import (
	"github.com/hacash/core/sys/inicnf"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type Inicnf struct {
	inicnf.File

	// cnf cache
	mustDataDir string
}

// val list
func (i *Inicnf) StringValueList(section string, name string) []string {
	valstr := i.Section(section).Key(name).MustString("")
	valstr = regexp.MustCompile(`[,，\s]+`).ReplaceAllString(valstr, ",")
	valstr = strings.Trim(valstr, ",")
	if valstr == "" {
		return []string{}
	}
	return strings.Split(valstr, ",")
}

// data dir
func (i *Inicnf) MustDataDir() string {
	if i.mustDataDir != "" {
		return i.mustDataDir
	}
	dir := i.Section("").Key("data_dir").MustString("~/.hacash_mainnet")
	if strings.HasPrefix(dir, "~/") {
		dir = os.Getenv("HOME") + string([]byte(dir)[1:])
	}
	i.mustDataDir = dir
	return dir
}

//////////////////////////////

func LoadInicnf(source_file string) (*Inicnf, error) {
	rand.Seed(time.Now().Unix())
	inifile, err := inicnf.Load(source_file)
	if err != nil {
		return nil, err
	}
	cnf := &Inicnf{}
	cnf.File = *inifile
	return cnf, nil
}
