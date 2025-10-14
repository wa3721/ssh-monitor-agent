package scriptinit

import (
	"go.uber.org/zap"

	"sshmonitor/config"

	"os"
)

//1.检查是否存在dst目录 /etc/profile.d
//2.检查项目内是否存在脚本
//3.检查目录内是否已经存在脚本，存在则skip

type checkfunc func()

type Checklist struct {
	checkDirectory         checkfunc
	checkScript            checkfunc
	checkScriptInDirectory checkfunc
}

const (
	targetDir  = "/etc/profile.d"   // 目标目录
	scriptName = "command_audit.sh" // 脚本文件名
	projectDir = "./"               // 项目目录，假设脚本在当前目录下
)

func NewChecklist() *Checklist {
	return &Checklist{
		checkDirectory:         checkDirectory,
		checkScript:            checkScript,
		checkScriptInDirectory: checkScriptInDirectory,
	}
}

// checkDirectory 检查目标目录（如 /etc/profile.d）是否存在
func checkDirectory() {
	_, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			config.GlobalLogger.Fatal("❌ target directory  is not exist。\n", zap.String("", targetDir))
		} else {
			config.GlobalLogger.Fatal("⚠️  check directory  error: \n", zap.String("", targetDir), zap.Error(err))
		}
		return
	}
	config.GlobalLogger.Info("✅ target directory '%s' is exist。\n", zap.String("", targetDir))
}

// checkScript 检查项目内（如当前目录）是否存在指定的脚本文件
func checkScript() {
	scriptPath := projectDir + scriptName
	_, err := os.Stat(scriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			config.GlobalLogger.Fatal("❌ No script files found within the project.\n")
		} else {
			config.GlobalLogger.Fatal("⚠️  An error occurred while checking the test item script: \n", zap.Error(err))
		}
		return
	}
	config.GlobalLogger.Info("✅ Script files exist within the project.\n", zap.String("", scriptPath))
}

// checkScriptInDirectory 检查目标目录中是否已存在同名脚本，避免重复部署
func checkScriptInDirectory() {
	targetScriptPath := targetDir + "/" + scriptName
	_, err := os.Stat(targetScriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			config.GlobalLogger.Info("ℹ️  The target directory does not yet contain deployed scripts and is ready for installation.\n", zap.String("", targetScriptPath))
		} else {
			config.GlobalLogger.Fatal("⚠️  An error occurred while checking the target directory script: \n", zap.String("", targetScriptPath), zap.Error(err))
		}
		return
	}
	config.GlobalLogger.Info("⏩ The target directory already contains the script; deployment skipped.\n", zap.String("", targetScriptPath))
}

// RunAll 方法按顺序执行所有检查
func (cl *Checklist) RunAll() {
	config.GlobalLogger.Info("Begin performing pre-deployment checks...")
	cl.checkDirectory()
	cl.checkScript()
	cl.checkScriptInDirectory()
	config.GlobalLogger.Info("All inspections have been completed.")
}

func Exec() {
	scriptPath := projectDir + scriptName
	targetScriptPath := targetDir + "/" + scriptName
	file, err := os.ReadFile(scriptPath)
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
	err = os.WriteFile(targetScriptPath, file, 0644)
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
}
