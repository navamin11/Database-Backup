package handlers

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger       *zap.Logger
	datenow      = time.Now().Format("2006-01-02")
	logPath      = "./logs"
	logExtension = ".log"
	file         = path.Join(logPath, datenow+logExtension)
)

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func InitLogger() {
	// Set up lumberjack as a logger:
	logFile := &lumberjack.Logger{
		Filename:   file, // Or any other path
		MaxSize:    500,  // MB; after this size, a new log file is created
		MaxBackups: 1,    // Number of backups to keep
		MaxAge:     2,    // Days
		Compress:   true, // Compress the backups using gzip
	}

	// Define config for the console output
	prod := zapcore.EncoderConfig{
		TimeKey:    "timestamp",
		LevelKey:   "level",
		EncodeTime: SyslogTimeEncoder,
		MessageKey: "message",
		// EncodeLevel:  CustomEncodeLevel,
		// NameKey: "logger",
		// CallerKey:     "caller",
		// StacktraceKey: "stacktrace",
		// LineEnding: zapcore.DefaultLineEnding,
		// EncodeDuration: milliSecondsDurationEncoder,
		// EncodeCaller:   zapcore.ShortCallerEncoder,
		// EncodeName: zapcore.FullNameEncoder,
	}
	prod.EncodeLevel = zapcore.CapitalLevelEncoder

	fileEncoder := zapcore.NewJSONEncoder(prod)

	// Define config for the console output
	dev := zapcore.EncoderConfig{
		TimeKey:    "timestamp",
		LevelKey:   "level",
		EncodeTime: SyslogTimeEncoder,
		MessageKey: "message",
		// EncodeLevel:  CustomEncodeLevel,
		// NameKey: "logger",
		// CallerKey:     "caller",
		// StacktraceKey: "stacktrace",
		// LineEnding: zapcore.DefaultLineEnding,
		// EncodeDuration: milliSecondsDurationEncoder,
		// EncodeCaller:   zapcore.ShortCallerEncoder,
		// EncodeName: zapcore.FullNameEncoder,
	}

	dev.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(dev)

	core := zapcore.NewTee(
		// file log
		// zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.WarnLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.ErrorLevel),
		// zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DPanicLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.PanicLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.FatalLevel),
		// console log
		// zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.WarnLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel),
		// zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DPanicLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.PanicLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.FatalLevel),
	)

	Logger = zap.New(core, zap.AddCaller())

	defer Logger.Sync()
}

func CreateDirectory() (*string, error) {
	dir := time.Now().Format("2006-01-02")

	if err := os.RemoveAll(dir); err != nil {
		return &dir, err
	}

	// create a directory and give it required permissions
	if err := os.Mkdir(dir, 0o755); err != nil {
		return &dir, err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &dir, err
	}
	return &dir, nil
}

func CreateSubDirectory(mainDir, subDirHost, subDirDb string) (*string, error) {
	path := fmt.Sprintf("%s/%s/%s", mainDir, subDirHost, subDirDb)

	// create a directory and give it required permissions
	if err := os.MkdirAll(path, 0o755); err != nil {
		return &path, err
	}
	return &path, nil
}

func ConnectPostgreSQL(driver, user, pw, host, port, name string) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", host, port, user, pw, name)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}
	return db, err
}

func ConnectMSSQL(driver, user, pw, host, port, name string) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("%s://%s:%s@%s:%v?database=%s", driver, user, pw, host, port, name)

	// Create connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return db, err
	}
	return db, err
}

func ConnectMySQL(driver, user, pw, host, port, name string) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, pw, host, port, name)

	// Create connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, err
	}

	if err := db.Ping(); err != nil {
		return db, err
	}
	return db, err
}

func ZipSource(source, target string) error {
	dir := fmt.Sprintf("./%s", source)
	filename := target + ".zip"

	// Create a ZIP file and zip.Writer
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// Go through all the files of the source
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(dir), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func DeleteDirectory(mainDir string) error {
	dir := fmt.Sprintf("./%s", mainDir)

	// check directory file
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		} else {
			return err
		}
	}

	// delete all
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}
