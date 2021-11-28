package sys

import (
	"fmt"
	"github.com/hacash/core/sys/inicnf"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 全局开发测试标记
var TestDebugLocalDevelopmentMark bool = false

type Inicnf struct {
	inicnf.File

	// cnf cache
	mustDataDirVersion string
	mustDataDir        string

	// 最低可被当前兼容的区块链数据库（仅blockdata）版本号
	blockChainStateDatabaseLowestCompatibleVersion int
	// 当前使用的区块链数据库版本号
	blockChainStateDatabaseCurrentUseVersion int

	mux sync.Mutex
}

// 必须调用
func (i *Inicnf) SetDatabaseVersion(curversion, compatible int) {
	i.blockChainStateDatabaseCurrentUseVersion = curversion
	i.blockChainStateDatabaseLowestCompatibleVersion = compatible
}
func (i *Inicnf) GetDatabaseVersion() (int, int) {
	return i.blockChainStateDatabaseCurrentUseVersion, i.blockChainStateDatabaseLowestCompatibleVersion
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

func AbsDir(dir string) string {
	if path.IsAbs(dir) == false {
		ppp, err := filepath.Abs(os.Args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		dir = path.Join(path.Dir(ppp), dir)
	}
	return dir
}

// data dir version
func (i *Inicnf) MustDataDirWithVersion() string {
	i.mux.Lock()
	defer i.mux.Unlock()
	if i.mustDataDirVersion != "" {
		return i.mustDataDirVersion
	}
	dir := i.mustDataDirUnsafe()
	dir = path.Join(dir, fmt.Sprintf("v%d", i.blockChainStateDatabaseCurrentUseVersion))
	i.mustDataDirVersion = dir
	return dir
}

// data dir
func (i *Inicnf) MustDataDir() string {
	i.mux.Lock()
	defer i.mux.Unlock()

	if len(i.mustDataDir) > 0 {
		return i.mustDataDir
	}
	dir := i.mustDataDirUnsafe()
	i.mustDataDir = dir
	return dir
}

// data dir
func (i *Inicnf) mustDataDirUnsafe() string {
	dir := i.Section("").Key("data_dir").MustString("~/.hacash_mainnet")
	if strings.HasPrefix(dir, "~/") {
		dir = os.Getenv("HOME") + string([]byte(dir)[1:])
	}
	dir = AbsDir(dir)
	return dir
}

// data dir Check Version
func (i *Inicnf) MustDataDirCheckVersion(version int) (string, bool) {
	dir := i.Section("").Key("data_dir").MustString("~/.hacash_mainnet")
	if strings.HasPrefix(dir, "~/") {
		dir = os.Getenv("HOME") + string([]byte(dir)[1:])
	}
	dir = AbsDir(dir)
	dir = path.Join(dir, fmt.Sprintf("v%d", version))
	// 检查是否存在
	_, nte := os.Stat(dir)
	if nte != nil {
		return dir, false // 不存在
	}
	// 目录存在
	return dir, true
}

//////////////////////////////

func LoadInicnf(source_file string) (*Inicnf, error) {
	rand.Seed(time.Now().Unix())
	inifile, err := inicnf.LooseLoad(source_file)
	if err != nil {
		return nil, err
	}
	cnf := &Inicnf{}
	cnf.File = *inifile
	return cnf, nil
}
