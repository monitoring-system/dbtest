package randgen

import (
	"bytes"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SQLGenerator struct {
	TestName string
	// 存放yyzz文件的path
	ConfPath string
	//  存放结构的path
	ResultPath string
	//  randgen主目录
	RmPath string

	CachedQueries []string
}

/*
container env

ConfPath = "/root/conf"
RmPath = "/root/randgenx"
ResultPath = "/root/result"
*/
const CONFPATH = "CONFPATH"
const RMPATH = "RMPATH"
const RESULTPATH = "RESULTPATH"

var ConfPath = os.Getenv(CONFPATH)
var RmPath = os.Getenv(RMPATH)
var ResultPath = os.Getenv(RESULTPATH)

func (rl *SQLGenerator) Init(testName string) {
	rl.ConfPath = ConfPath
	rl.ResultPath = ResultPath
	rl.RmPath = RmPath
	rl.TestName = testName
}

func (rl *SQLGenerator) LoadData(zzContent string, yyContent string, dbname string, queries int) (sqls []string, err error) {
	zzPath := fmt.Sprintf(filepath.Join(rl.ConfPath, "%s.zz"), rl.TestName)
	yyPath := fmt.Sprintf(filepath.Join(rl.ConfPath, "%s.yy"), rl.TestName)
	ioutil.WriteFile(zzPath, []byte(zzContent), os.ModePerm)
	ioutil.WriteFile(yyPath, []byte(yyContent), os.ModePerm)

	rPath := filepath.Join(rl.ResultPath, rl.TestName)
	_, err = execShell(rl.RmPath, "perl", "gentest.pl",
		fmt.Sprintf("--dsn=dummy:file:%s", rPath),
		fmt.Sprintf("--gendata=%s", zzPath),
		fmt.Sprintf("--grammar=%s", yyPath),
		fmt.Sprintf("--queries=%d", queries))

	if err != nil {
		return nil, err
	}

	f, err := os.Open(rPath)
	if err != nil {
		return nil, err
	}
	sqlBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	data, grammar := splitToDataAndGrammar(sqlBytes)
	if dbname != "test" {
		for i, d := range data {
			if d == "USE test" {
				data[i] = fmt.Sprintf("USE %s", dbname)
			}
			if d == "CREATE SCHEMA /*!IF NOT EXISTS*/ test" {
				data[i] = fmt.Sprintf("CREATE SCHEMA /*!IF NOT EXISTS*/ %s", dbname)
			}
		}
	}

	rl.CachedQueries = grammar

	return data, nil
}

func splitToDataAndGrammar(totalContent []byte) (data []string, grammar []string) {
	content := string(totalContent)

	gendataAndGrammar := strings.Split(content, "/* follow is grammar sql */;\n")

	return strings.Split(gendataAndGrammar[0], ";\n"), strings.Split(gendataAndGrammar[1], ";\n")
}

func (rl *SQLGenerator) Query() (sqls []string) {
	return rl.CachedQueries
}

// r1为Mysql输出  r2为TiDB输出
func (rl *SQLGenerator) Compare(r1 string, r2 string) (comment string, consistent bool) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(r1),
		B:        difflib.SplitLines(r2),
		FromFile: "Mysql",
		ToFile:   "Tidb",
	}
	text, _ := difflib.GetUnifiedDiffString(diff)

	return text, text == ""
}

//执行shell命令
func execShell(dir string, s string, args ...string) (string, error) {
	cmd := exec.Command(s, args...)

	var out bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Run()

	return out.String(), err
}
