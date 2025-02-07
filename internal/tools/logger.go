package tools

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	logBufferSize       = 100
	defaultLogChannelID = -4714261807
	batchSize           = 15
)

var (
	infoLogger   *log.Logger
	warnLogger   *log.Logger
	errorLogger  *log.Logger
	bot          *tgBotAPI.BotAPI
	logToChannel bool
	logChannelID int64
	logChan      chan string
	logBatch     []string
	lastSendTime time.Time
)

type LogContext struct {
	ChatID   int64
	Username string
}

func InitLogger(botAPI *tgBotAPI.BotAPI) {
	bot = botAPI

	// Инициализация настроек логирования в канал
	//logToChannel = os.Getenv("LOG_TO_CHANNEL") == "true"
	logToChannel = true
	logChannelID = defaultLogChannelID

	if channelID := os.Getenv("LOG_CHANNEL_ID"); channelID != "" {
		if id, err := strconv.ParseInt(channelID, 10, 64); err == nil {
			logChannelID = id
		} else {
			log.Printf("Ошибка парсинга LOG_CHANNEL_ID: %v, используется значение по умолчанию", err)
		}
	}

	file, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)

	if logToChannel {
		// Инициализация канала для логов с буфером
		logChan = make(chan string, logBufferSize)
		// Запуск горутины обработчика логов
		go processLogs()
		LogInfo("Логирование в канал включено (ID: %d)", logChannelID)
	} else {
		LogInfo("Логирование в канал отключено")
	}
}

func processLogs() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case msg := <-logChan:
			logBatch = append(logBatch, msg)
			if len(logBatch) >= batchSize {
				sendBatch()
			}
		case <-ticker.C:
			if len(logBatch) > 0 && time.Since(lastSendTime) > time.Second*9 {
				sendBatch()
			}
		}
	}
}

func sendBatch() {
	if len(logBatch) == 0 {
		return
	}

	text := "```\n" + strings.Join(logBatch, "\n") + "\n```"
	message := tgBotAPI.NewMessage(logChannelID, text)
	message.ParseMode = "MarkdownV2"

	if _, err := bot.Send(message); err != nil {
		log.Printf("Ошибка отправки батча логов: %v", err)
	}

	lastSendTime = time.Now()
	logBatch = logBatch[:0]
}

func formatMessage(level, msg string, ctx *LogContext) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if ctx != nil {
		return fmt.Sprintf("[%s] [%s] [chat:%d user:%s] %s", timestamp, level, ctx.ChatID, ctx.Username, msg)
	}
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, msg)
}

func sendToChannel(msg string) {
	if logToChannel && logChan != nil {
		select {
		case logChan <- msg:
			// Сообщение успешно добавлено в буфер
		default:
			// Буфер полон, пропускаем сообщение
			log.Printf("Буфер логов переполнен, сообщение пропущено")
		}
	}
}

func LogInfo(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("INFO", msg, nil)
	infoLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}

func LogInfoWithContext(ctx *LogContext, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("INFO", msg, ctx)
	infoLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}

func LogWarn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("WARN", msg, nil)
	warnLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}

func LogWarnWithContext(ctx *LogContext, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("WARN", msg, ctx)
	warnLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}

func LogError(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("ERROR", msg, nil)
	errorLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}

func LogErrorWithContext(ctx *LogContext, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	formattedMsg := formatMessage("ERROR", msg, ctx)
	errorLogger.Println(msg)
	fmt.Println(formattedMsg)
	sendToChannel(formattedMsg)
}
