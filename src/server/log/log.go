﻿package log

import (
	"pnas/setting"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugarLog *zap.SugaredLogger

func Init() {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   setting.GS.Log.FileName, // 日志文件路径，默认 os.TempDir()
		MaxSize:    10,                      // 每个日志文件保存10M，默认 100M
		MaxBackups: 100,                     // 保留100个备份，默认不限
		MaxAge:     7,                       // 保留7天，默认不限
		Compress:   true,                    // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	// 设置日志级别
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	switch setting.GS.Log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "linenum",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapEncodeLevel,                 // 小写编码器
		EncodeTime:       zapEncodeTime,                  // ISO8601 UTC 时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder, //
		EncodeCaller:     zapEncodeCaller,                // 全路径编码器
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewConsoleEncoder(encoderConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&write)), // 打印到控制台和文件
		write,
		level,
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	logger = zap.New(core, caller, development)

	opts := []zap.Option{
		zap.AddCaller(),      // 添加调用者信息
		zap.AddCallerSkip(1), // 跳过包装函数
		//zap.Fields(zap.String("key", "value")), // 自定义字段
	}

	logger.WithOptions(opts...)
	sugarLog = logger.Sugar().Desugar().WithOptions(opts...).Sugar()

	sugarLog.Info("DefaultLogger init success")
}

func zapEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func zapEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000Z0700") + "]")
}

func zapEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

func Debug(args ...interface{}) {
	sugarLog.Debug(args...)
}

func Info(args ...interface{}) {
	sugarLog.Info(args...)
}

func Warn(args ...interface{}) {
	sugarLog.Warn(args...)
}

func Error(args ...interface{}) {
	sugarLog.Error(args...)
}

func Panic(args ...interface{}) {
	sugarLog.Panic(args...)
}

func Debugf(template string, args ...interface{}) {
	sugarLog.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugarLog.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLog.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLog.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugarLog.Panicf(template, args...)
}
