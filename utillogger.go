package utilaio

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DEFAULT_BUFFER_SIZE = 1024
)

/*
log level
*/
const (
	NOLOG = iota
	INFO
	WARNING
	DEBUG
	TRACE
)

/*
log message prefix
*/
const (
	MSG_PREFIX_INFO    = "INFO"
	MSG_PREFIX_WARNING = "WARNING"
	MSG_PREFIX_DEBUG   = "DEBUG"
	MSG_PREFIX_TRACE   = "DEV"
)

type Log struct {
	receive           bool
	loglevel          int
	logsize           int
	buffersize        int
	filename          string
	location          string
	done              chan struct{}
	rotate            chan struct{}
	buffer            chan string
	logfiledescriptor *os.File
	wg                sync.WaitGroup
	mu                sync.Mutex
}

/*
syntax:

	Init(log_filename string, log_location string, log_level int, log_size int, buffer_size int)

mandatory:

	log_filename, log_location, log_level

optional:

	log_size, buffer_size
*/
func Init(log_filename string, log_location string, log_level int, args ...interface{}) *Log {

	var log_size int
	var buffer_size = DEFAULT_BUFFER_SIZE

	for index, arg := range args {
		switch index {
		case 0:
			if arg != nil {
				log_size = arg.(int)
			}
		case 1:
			if arg != nil {
				buffer_size = arg.(int)
			}
		}
	}

	logger := &Log{
		loglevel:   log_level,
		filename:   log_filename,
		location:   log_location,
		logsize:    log_size,
		buffersize: buffer_size,
		buffer:     make(chan string, buffer_size),
		done:       make(chan struct{}, 1),
		rotate:     make(chan struct{}, 1),
	}

	return logger
}

/*
syntax:

	Update(log_filename string, log_location string, log_level int, log_size int, buffer_size int)

mandatory:

	log_filename, log_location, log_level

optional:

	log_size, buffer_size
*/
func (logger *Log) Update(log_filename string, log_location string, log_level int, args ...interface{}) error {
	var err error
	var buffer_size = DEFAULT_BUFFER_SIZE

	log_size := logger.logsize

	for index, arg := range args {
		switch index {
		case 0:
			if arg != nil {
				log_size = arg.(int)
			}
		case 1:
			if arg != nil {
				buffer_size = arg.(int)
			}
		}
	}

	if (logger.filename != log_filename) || (logger.location != log_location) || (logger.loglevel != log_level) || (logger.logsize != log_size) || (logger.buffersize != buffer_size) {
		logger.Stop()

		logger.filename = log_filename
		logger.location = log_location
		logger.loglevel = log_level
		logger.logsize = log_size
		logger.buffersize = buffer_size
		logger.buffer = make(chan string, buffer_size)

		err = logger.Start()
	}
	return err
}

func (logger *Log) open() error {
	var err error

	err = os.MkdirAll(logger.location, 0777)
	if err != nil {
		return err
	}

	logfilename := filepath.Join(logger.location, logger.filename)

	logger.logfiledescriptor, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return err
}

func (logger *Log) close() error {
	return logger.logfiledescriptor.Close()
}

func (logger *Log) Start() error {

	logger.mu.Lock()

	if logger.receive {
		return fmt.Errorf("logger already in running state")
	}

	err := logger.open()

	if err == nil {
		logger.logfiledescriptor.Write([]byte(formatMessage(MSG_PREFIX_INFO, "STARTED")))
		logger.receive = true
		logger.wg.Add(1)
		go func() {
			logger.flusher()
			logger.wg.Done()
		}()
	}

	logger.mu.Unlock()

	return err
}

func (logger *Log) Stop() {

	logger.mu.Lock()

	if logger.receive {
		logger.receive = false

		if len(logger.done) < 1 {
			logger.done <- struct{}{}
		}

		logger.wg.Wait()
		logger.logfiledescriptor.Write([]byte(formatMessage(MSG_PREFIX_INFO, "STOPPED")))
		logger.close()
	}

	logger.mu.Unlock()
}

func formatMessage(level string, msg string) string {
	now := time.Now()
	return fmt.Sprintf("[%s]\t[%s]\t%s\n", level, now.Format(time.ANSIC), msg)
}

func (logger *Log) LogMessage(level int, format string, args ...interface{}) {

	if logger.loglevel >= level {

		var prefix string

		switch level {
		case INFO:
			prefix = MSG_PREFIX_INFO
		case WARNING:
			prefix = MSG_PREFIX_WARNING
		case DEBUG:
			prefix = MSG_PREFIX_DEBUG
		case TRACE:
			prefix = MSG_PREFIX_TRACE
		default:
			return
		}

		logMessage := formatMessage(prefix, fmt.Sprintf(format, args...))

		logger.bufferMessage(logMessage)
	}
}

func (logger *Log) RenameOldLogFiles() {

	newest_log_filename := filepath.Join(logger.location, logger.filename)
	log_filename_1 := filepath.Join(logger.location, logger.filename+".1")
	log_filename_2 := filepath.Join(logger.location, logger.filename+".2")
	log_filename_3 := filepath.Join(logger.location, logger.filename+".3")
	log_filename_4 := filepath.Join(logger.location, logger.filename+".4")
	log_filename_5 := filepath.Join(logger.location, logger.filename+".5")
	log_filename_6 := filepath.Join(logger.location, logger.filename+".6")
	log_filename_7 := filepath.Join(logger.location, logger.filename+".7")
	log_filename_8 := filepath.Join(logger.location, logger.filename+".8")
	log_filename_9 := filepath.Join(logger.location, logger.filename+".9")
	oldest_log_filename := filepath.Join(logger.location, logger.filename+".10")

	if _, err := os.Stat(log_filename_9); err == nil {
		os.Rename(log_filename_9, oldest_log_filename)
	}
	if _, err := os.Stat(log_filename_8); err == nil {
		os.Rename(log_filename_8, log_filename_9)
	}
	if _, err := os.Stat(log_filename_7); err == nil {
		os.Rename(log_filename_7, log_filename_8)
	}
	if _, err := os.Stat(log_filename_6); err == nil {
		os.Rename(log_filename_6, log_filename_7)
	}
	if _, err := os.Stat(log_filename_5); err == nil {
		os.Rename(log_filename_5, log_filename_6)
	}
	if _, err := os.Stat(log_filename_4); err == nil {
		os.Rename(log_filename_4, log_filename_5)
	}
	if _, err := os.Stat(log_filename_3); err == nil {
		os.Rename(log_filename_3, log_filename_4)
	}
	if _, err := os.Stat(log_filename_2); err == nil {
		os.Rename(log_filename_2, log_filename_3)
	}
	if _, err := os.Stat(log_filename_1); err == nil {
		os.Rename(log_filename_1, log_filename_2)
	}
	if _, err := os.Stat(newest_log_filename); err == nil {
		os.Rename(newest_log_filename, log_filename_1)
	}
}

func (logger *Log) Rotate() {

	logger.mu.Lock()

	logfilename := filepath.Join(logger.location, logger.filename)

	if file, err := os.Stat(logfilename); err == nil {
		size := file.Size()
		if size > int64(logger.logsize) {
			if len(logger.rotate) < 1 {
				logger.rotate <- struct{}{}
			}
		}
	}

	logger.mu.Unlock()
}

func (logger *Log) bufferMessage(message string) {

	logger.mu.Lock()

	if len(logger.buffer) < logger.buffersize && logger.receive {
		logger.buffer <- message
	}

	logger.mu.Unlock()
}

func (logger *Log) flusher() {
	for {
		select {
		case <-logger.done:
			for {
				select {
				case msg := <-logger.buffer:
					logger.logfiledescriptor.Write([]byte(msg))
				default:
					return
				}
			}
		case <-logger.rotate:
			logger.close()
			logger.RenameOldLogFiles()
			logger.open()
		case msg := <-logger.buffer:
			logger.logfiledescriptor.Write([]byte(msg))
		}
	}
}

func (logger *Log) QuickLog(level int, format string, args ...interface{}) {

	logger.mu.Lock()
	if logger.loglevel >= level {

		var prefix string

		switch level {
		case INFO:
			prefix = MSG_PREFIX_INFO
		case WARNING:
			prefix = MSG_PREFIX_WARNING
		case DEBUG:
			prefix = MSG_PREFIX_DEBUG
		case TRACE:
			prefix = MSG_PREFIX_TRACE
		default:
			return
		}

		logMessage := formatMessage(prefix, fmt.Sprintf(format, args...))

		err := logger.open()

		if err == nil {
			logger.logfiledescriptor.Write([]byte(logMessage))
		}

		logger.close()
	}

	logger.mu.Unlock()
}
