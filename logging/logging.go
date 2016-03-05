package logging

import (
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/sysdCommonDefs"
	"io/ioutil"
	"log/syslog"
)

func ConvertLevelStrToVal(str string) sysdCommonDefs.SRDebugLevel {
	var val sysdCommonDefs.SRDebugLevel
	switch str {
	case "crit":
		val = sysdCommonDefs.CRIT
	case "err":
		val = sysdCommonDefs.ERR
	case "warn":
		val = sysdCommonDefs.WARN
	case "alert":
		val = sysdCommonDefs.ALERT
	case "emerg":
		val = sysdCommonDefs.EMERG
	case "notice":
		val = sysdCommonDefs.NOTICE
	case "info":
		val = sysdCommonDefs.INFO
	case "debug":
		val = sysdCommonDefs.DEBUG
	case "trace":
		val = sysdCommonDefs.TRACE
	}
	return val
}

type ComponentJson struct {
	Module string `json:Module`
	Level  string `json:Level`
}

type LoggingJson struct {
	SystemLogging string          `json:SystemLogging`
	Components    []ComponentJson `json:Components`
}

type Writer struct {
	GlobalLogging   bool
	MyComponentName string
	MyLogLevel      sysdCommonDefs.SRDebugLevel
	initialized     bool
	sysLogger       *syslog.Writer
	subSocket       *nanomsg.SubSocket
	socketCh        chan []byte
}

func NewLogger(paramsDir string, name string, tag string) (*Writer, error) {
	var loggingConfig LoggingJson
	var err error
	srLogger := &Writer{}
	srLogger.MyComponentName = name

	fileName := paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}
	fileName = fileName + "logging.json"

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Failed to read logging config file")
		return nil, err
	}
	err = json.Unmarshal(data, &loggingConfig)
	if err != nil {
		fmt.Println("Failed to unmarshal logging config: ", err)
		return nil, err
	}

	if loggingConfig.SystemLogging == "on" {
		srLogger.GlobalLogging = true
	}
	for _, component := range loggingConfig.Components {
		if component.Module == name {
			srLogger.MyLogLevel = ConvertLevelStrToVal(component.Level)
			break
		}
	}

	srLogger.sysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err == nil {
		srLogger.initialized = true
		fmt.Println("Logging level ", srLogger.MyLogLevel, " set for ", srLogger.MyComponentName)
	}
	return srLogger, err
}

func (logger *Writer) SetGlobal(Enable bool) error {
	logger.GlobalLogging = Enable
	fmt.Println("Changed global logging to: ", logger.GlobalLogging, " for ", logger.MyComponentName)
	return nil
}

func (logger *Writer) SetLevel(level sysdCommonDefs.SRDebugLevel) error {
	logger.MyLogLevel = level
	fmt.Println("Changed logging level to: ", logger.MyLogLevel, " for ", logger.MyComponentName)
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
	if logger.initialized && logger.GlobalLogging && logger.MyLogLevel >= sysdCommonDefs.NOTICE {
		return logger.sysLogger.Notice(message)
	}
	return nil
}

func (logger *Writer) Info(message string) error {
	if logger.initialized && logger.GlobalLogging && logger.MyLogLevel >= sysdCommonDefs.INFO {
		return logger.sysLogger.Info(message)
	}
	return nil
}

func (logger *Writer) Debug(message string) error {
	if logger.initialized && logger.GlobalLogging && logger.MyLogLevel >= sysdCommonDefs.DEBUG {
		return logger.sysLogger.Debug(message)
	}
	return nil
}

func (logger *Writer) Write(message string) (int, error) {
	if logger.initialized && logger.GlobalLogging && logger.MyLogLevel >= sysdCommonDefs.TRACE {
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

func (logger *Writer) SetupSubSocket() error {
	var err error
	var socket *nanomsg.SubSocket
	if socket, err = nanomsg.NewSubSocket(); err != nil {
		logger.Err(fmt.Sprintf("Failed to create subscribe socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	if err = socket.Subscribe(""); err != nil {
		logger.Err(fmt.Sprintf("Failed to subscribe to \"\" on subscribe socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	if _, err = socket.Connect(sysdCommonDefs.PUB_SOCKET_ADDR); err != nil {
		logger.Err(fmt.Sprintf("Failed to connect to publisher socket %s, error:%s", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}

	logger.Info(fmt.Sprintf("Connected to publisher socker %s", sysdCommonDefs.PUB_SOCKET_ADDR))
	if err = socket.SetRecvBuffer(1024 * 1024); err != nil {
		logger.Err(fmt.Sprintln("Failed to set the buffer size for subsriber socket %s, error:", sysdCommonDefs.PUB_SOCKET_ADDR, err))
		return err
	}
	logger.subSocket = socket
	logger.socketCh = make(chan []byte)
	return nil
}

func (logger *Writer) ProcessSysdNotification(rxBuf []byte) error {
	var msg sysdCommonDefs.Notification
	err := json.Unmarshal(rxBuf, &msg)
	if err != nil {
		logger.Err(fmt.Sprintln("Unable to unmarshal sysd notification: ", rxBuf))
		return err
	}
	if msg.Type == sysdCommonDefs.G_LOG {
		var gLog sysdCommonDefs.GlobalLogging
		err = json.Unmarshal(msg.Payload, &gLog)
		if err != nil {
			logger.Err(fmt.Sprintln("Unable to unmarshal sysd global logging notification: ", msg.Payload))
			return err
		}
		logger.SetGlobal(gLog.Enable)
	}
	if msg.Type == sysdCommonDefs.C_LOG {
		var cLog sysdCommonDefs.ComponentLogging
		err = json.Unmarshal(msg.Payload, &cLog)
		if err != nil {
			logger.Err(fmt.Sprintln("Unable to unmarshal sysd component logging notification: ", msg.Payload))
			return err
		}
		if cLog.Name == logger.MyComponentName {
			logger.SetLevel(cLog.Level)
		}
	}
	return nil
}

func (logger *Writer) ProcessSysdNotifications() error {
	for {
		select {
		case rxBuf := <-logger.socketCh:
			if rxBuf != nil {
				logger.ProcessSysdNotification(rxBuf)
			}
		}
	}
	return nil
}

func (logger *Writer) ListenForSysdNotifications() error {
	err := logger.SetupSubSocket()
	if err != nil {
		logger.Err(fmt.Sprintln("Failed to subscribe to sysd notifications"))
		return err
	}
	go logger.ProcessSysdNotifications()
	for {
		rxBuf, err := logger.subSocket.Recv(0)
		if err != nil {
			logger.Err(fmt.Sprintln("Recv on BFD subscriber socket failed with error:", err))
			continue
		}
		logger.socketCh <- rxBuf
	}
}
