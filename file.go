package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger information
type FileLogger struct {
	level    Level
	fileName string
	filePath string
	file     *os.File
	errFile  *os.File
	maxSize  int64
}

// NewFileLogger Constructor for the file log structure
func NewFileLogger(level, fileName, filePath string) *FileLogger {
	logLevel := parseLogLevel(level)
	fl := &FileLogger{
		level:    logLevel,
		fileName: fileName,
		filePath: filePath,
		maxSize:  10 * 1024 * 1024,
	}
	fl.initFile()
	return fl
}

func (f *FileLogger) initFile() {
	logName := path.Join(f.filePath, f.fileName)
	// open log file
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("Open file %s failed", logName, err))
	}
	f.file = fileObj
	// open errlog file
	errLogName := fmt.Sprintf("%s.err", logName)
	errFileObj, err := os.OpenFile(errLogName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Errorf("Open file %s failed", errLogName, err))
	}
	f.errFile = errFileObj
}

func (f *FileLogger) checkSplit(file *os.File) bool {
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	return fileSize >= f.maxSize
}

// Segmentation log
func (f *FileLogger) splitLogFile(file *os.File) *os.File {
	fileName := file.Name() // get full path
	backupName := fmt.Sprintf("%s_%v.bak", fileName, time.Now().Unix())
	file.Close()
	os.Rename(fileName, backupName)
	fileObj, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		panic(fmt.Errorf("OpenFileError"))
	}
	return fileObj
}

func (f *FileLogger) log(level Level, format string, args ...interface{}) {
	if f.level == level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")				
	fileName, line, funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s]%s",
		nowStr, fileName, line, funcName, logLevelStr, msg)
	if f.checkSplit(f.file) {
		f.file = f.splitLogFile(f.errFile)
	}
	fmt.Fprintln(f.errFile, logMsg)

	if level >= ErrorLevel {
		if f.checkSplit(f.errFile) {
			f.errFile = f.splitLogFile(f.errFile)
		}
		fmt.Fprintln(f.errFile, logMsg)
	}
}

// Debug debug方法
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.log(DebugLevel, format, args...)
}

// Info info方法
func (f *FileLogger) Info(format string, args ...interface{}) {
	f.log(InfoLevel, format, args...)
}

// Warn warn方法
func (f *FileLogger) Warn(format string, args ...interface{}) {
	f.log(WarningLevel, format, args...)
}

// Error error方法
func (f *FileLogger) Error(format string, args ...interface{}) {
	f.log(ErrorLevel, format, args...)
}

// Fatal fatal方法
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.log(FatalLevel, format, args...)
}

// Close 关闭日志文件句柄
func (f *FileLogger) Close() {
	f.file.Close()
	f.errFile.Close()
}
