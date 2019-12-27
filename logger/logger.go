package logger

import (
	"io"
	"log"
	"os"
)

//Logger defines and internal structure that contains all configuration for the log package
type Logger struct {
	log *log.Logger
}

//NewLogger creates a new instance of the Logger wrapper.
func NewLogger(file string) *Logger {
	out, _ := os.Create(file)
	wrt := io.MultiWriter(os.Stdout, out)
	newLog := log.New(wrt, "", log.LstdFlags)
	return &Logger{
		log: newLog,
	}
}

func (z *Logger) Println(v interface{}) {
	z.log.Println(v)
}

func (z *Logger) Printf(format string, v ...interface{}) {
	z.log.Printf(format, v...)
}

func (z *Logger) Panic(v interface{}) {
	z.log.Panic(v)
}

func (z *Logger) Panicf(format string, v ...interface{}) {
	z.log.Panicf(format, v...)
}

func (z *Logger) Panicln(v interface{}) {
	z.log.Panicln(v)
}

func (z *Logger) Fatal(v interface{}) {
	z.log.Fatal(v)
}

func (z *Logger) Fatalf(format string, v ...interface{}) {
	z.log.Fatalf(format, v...)
}

func (z *Logger) Fatalln(v interface{}) {
	z.log.Fatalln(v)
}
