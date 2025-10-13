package conf

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/knadh/koanf/providers/confmap"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	k "github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Koan    = k.New(".")
	RunMode string
	once    sync.Once
	Log     *zerolog.Logger

	//go:embed config.toml
	confgo []byte
)

const (
	OS_WINDOWS      = 1
	OS_MAC          = 2
	OS_LINUX        = 3
	OS_BSD          = 4
	OS_UNIX_SOLARIS = 5
	OS_UNIX_IBM     = 6
	OS_MAINFRAME    = 7
)

func StartConfig() {
	if err := Koan.Load(rawbytes.Provider(confgo), toml.Parser()); err != nil {
		fmt.Printf("ERROR=%e\n", err)
		return
	}
	RunMode = Koan.String("runmode")
	Log = GetLogger()

}

func GetLogger() *zerolog.Logger {
	var zl *zerolog.Logger
	once.Do(func() {
		zl = configure()
	})
	return zl
}

func DetectOs() int {
	os := runtime.GOOS
	switch os {
	case "windows":
		return OS_WINDOWS
	case "darwin":
		return OS_MAC
	case "linux": // linux family
	case "android":
		return OS_LINUX
	case "netbsd": // BSD family
	case "freebsd":
	case "openbsd":
	case "dragonfly":
		return OS_BSD
	case "solaris": // unix solaris
	case "illumos": // unix solaris derivatives
		return OS_UNIX_SOLARIS
	case "aix": //unix ibm
		return OS_UNIX_IBM
	case "zos": //mainframe
		return OS_MAINFRAME
	}
	return 0
}

func configure() *zerolog.Logger {
	extractFilePath()
	createFileFolder()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	level := Koan.String(RunMode + ".logging_level")
	switch level {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	var writers []io.Writer
	if Koan.Bool(RunMode + ".console_logging") {
		writers = append(writers, newConsoleWriter())
	}
	if Koan.Bool(RunMode + ".file_logging") {
		writers = append(writers, newRollingFile())
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Caller().Logger()

	return &logger
}

func createFileFolder() {
	if _, err := os.Stat(Koan.String("filepath")); os.IsNotExist(err) {
		if err := os.MkdirAll(Koan.String("filepath"),
			0777); err != nil {
			panic(err)
			return
		}
	}
	if _, err := os.Stat(Koan.String("filestore")); os.IsNotExist(err) {
		if err := os.MkdirAll(Koan.String("filestore"),
			0774); err != nil {
			panic(err)
			return
		}
	}

	if _, err := os.Stat(Koan.String("completeFilename")); err == nil {
		fmt.Println("file exist")
		return
	}
	_, err := os.OpenFile(Koan.String("completeFilename"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	//err = os.Chmod(Koan.String("completeFilename"), 0774)
	//if err != nil {
	//	panic(err)
	//}
}

func newRollingFile() io.Writer {

	lumberjackwriter := &lumberjack.Logger{
		Filename:   Koan.String("completeFilename"),
		MaxBackups: Koan.Int(RunMode + ".log_maxbackup"), // files
		MaxSize:    Koan.Int(RunMode + ".log_maxsize"),   // megabytes
		MaxAge:     Koan.Int(RunMode + ".log_maxage"),    // days
		Compress:   true,
	}
	//fmt.Printf("rolling completeFilename %s\n", Koan.String("completeFilename"))
	//writer := diode.NewWriter(lumberjackwriter, 1000, 10*time.Millisecond, func(missed int) {
	//	fmt.Printf("Logger Dropped %d messages", missed)
	//})
	return lumberjackwriter
}

func newConsoleWriter() io.Writer {
	output := zerolog.ConsoleWriter{Out: os.Stdout,
		TimeFormat: time.RFC3339Nano, PartsOrder: []string{zerolog.
				LevelFieldName, zerolog.TimestampFieldName,
			zerolog.CallerFieldName, zerolog.MessageFieldName}}
	output.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf(" %s | ", i)
	}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("---> %s", i)
	}
	output.FormatCaller = func(i interface{}) string {
		return fmt.Sprintf(" %s ", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf(" %s: ", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToLower(fmt.Sprintf(" %s ", i))
	}
	return output
}

func extractFilePath() {
	Koan.Load(confmap.Provider(map[string]interface{}{"os": DetectOs()}, "."), nil)
	Koan.Load(confmap.Provider(map[string]interface{}{"portapi": Koan.String(RunMode + ".port_api")}, "."), nil)
	var filepath, filename, filestore string

	if Koan.Int("os") == OS_WINDOWS {
		filepath = Koan.String(RunMode + ".log_path_windows")
		filename = Koan.String(RunMode + ".log_filename")
		filestore = Koan.String(RunMode + ".file_path_windows")
	} else {
		filepath = Koan.String(RunMode + ".log_path")
		filename = Koan.String(RunMode + ".log_filename")
		filestore = Koan.String(RunMode + ".file_path")
	}

	completeFilename := filepath + string(os.PathSeparator) + filename

	Koan.Load(confmap.Provider(map[string]interface{}{"completeFilename": strings.ToLower(completeFilename)}, "."), nil)
	Koan.Load(confmap.Provider(map[string]interface{}{"filestore": filestore}, "."), nil)
	Koan.Load(confmap.Provider(map[string]interface{}{"filepath": filepath}, "."), nil)
}

func CheckRunMode() {
	StartConfig()
	//createLogFile()
	//file, err := os.OpenFile(Koan.String("completeFile"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	Log.Info().Msgf("logging to %s with RunMode: %s", Koan.String("completeFilename"), RunMode)
}
