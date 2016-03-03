package logging

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log/syslog"
)

//Logging levels
type SRDebugLevel uint8

const (
	CRIT   SRDebugLevel = 0
	ERR    SRDebugLevel = 1
	WARN   SRDebugLevel = 2
	ALERT  SRDebugLevel = 3
	EMERG  SRDebugLevel = 4
	NOTICE SRDebugLevel = 5
	INFO   SRDebugLevel = 6
	DEBUG  SRDebugLevel = 7
	TRACE  SRDebugLevel = 8
)

func ConvertLevelStrToVal(str string) SRDebugLevel {
	var val SRDebugLevel
	switch str {
	case "crit":
		val = CRIT
	case "err":
		val = ERR
	case "warn":
		val = WARN
	case "alert":
		val = ALERT
	case "emerg":
		val = EMERG
	case "notice":
		val = NOTICE
	case "info":
		val = INFO
	case "debug":
		val = DEBUG
	case "trace":
		val = TRACE
	}
	return val
}

type ComponentsJson struct {
	Name  string       `json:Name`
	Level SRDebugLevel `json:Level`
}

type LoggingJson struct {
	SystemLogging string           `json:SystemLogging`
	Components    []ComponentsJson `json:Components`
}

type Writer struct {
	componentName string
	globalLogging bool
	myLogLevel    SRDebugLevel
	initialized   bool
	sysLogger     *syslog.Writer
}

func NewLogger(name string, tag string) (*Writer, error) {
	var loggingConfig LoggingJson
	var err error
	srLogger := &Writer{}
	srLogger.componentName = name

	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	fileName := *paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}
	fileName = fileName + "logging.json"

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read logging config file")
		return nil, err
	}
	json.Unmarshal(data, &loggingConfig)

	if loggingConfig.SystemLogging == "on" {
		srLogger.globalLogging = true
	}
	for _, module := range loggingConfig.Components {
		if module.Name == name {
			srLogger.myLogLevel = module.Level
			break
		}
	}

	srLogger.sysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err == nil {
		srLogger.initialized = true
		fmt.Println("Logging level ", srLogger.myLogLevel, " set for ", srLogger.componentName)
	}
	return srLogger, err
}

func (logger *Writer) SetGlobal(logging string) error {
	if logging == "on" {
		logger.globalLogging = true
	} else {
		logger.globalLogging = false
	}
	return nil
}

func (logger *Writer) SetLevel(level SRDebugLevel) error {
	logger.myLogLevel = level
	fmt.Println("Changed logging level to: ", logger.myLogLevel, " for ", logger.componentName)
	return nil
}

func (logger *Writer) Crit(message string) error {
	if logger.initialized {
		return logger.sysLogger.Crit(message)
	}
	return nil
}

func (logger *Writer) Err(message string) error {
	if logger.initialized {
		return logger.sysLogger.Err(message)
	}
	return nil
}

func (logger *Writer) Warning(message string) error {
	if logger.initialized {
		return logger.sysLogger.Warning(message)
	}
	return nil
}

func (logger *Writer) Alert(message string) error {
	if logger.initialized {
		return logger.sysLogger.Alert(message)
	}
	return nil
}

func (logger *Writer) Emerg(message string) error {
	if logger.initialized {
		return logger.sysLogger.Emerg(message)
	}
	return nil
}

func (logger *Writer) Notice(message string) error {
	if logger.initialized && logger.globalLogging && logger.myLogLevel >= NOTICE {
		return logger.sysLogger.Notice(message)
	}
	return nil
}

func (logger *Writer) Info(message string) error {
	if logger.initialized && logger.globalLogging && logger.myLogLevel >= INFO {
		return logger.sysLogger.Info(message)
	}
	return nil
}

func (logger *Writer) Debug(message string) error {
	if logger.initialized && logger.globalLogging && logger.myLogLevel >= DEBUG {
		return logger.sysLogger.Debug(message)
	}
	return nil
}

func (logger *Writer) Write(message string) (int, error) {
	if logger.initialized && logger.globalLogging && logger.myLogLevel >= TRACE {
		n, err := logger.sysLogger.Write([]byte(message))
		return n, err
	}
	return 0, nil
}

func (logger *Writer) Close() error {
	var err error
	if logger.initialized {
		err = logger.sysLogger.Close()
	}
	logger = nil
	return err
}
