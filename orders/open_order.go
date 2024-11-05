package orders

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/BKadirkhodjaev/request-cli/util"
)

const (
	CommandName           string = "main"
	EnableDebug           bool   = true
	PersistentHttpThreads int    = 100
	FileName              string = "./import/data.csv"

	DikuTenant string = "diku"
	Hostname   string = "http://localhost:9130"

	LoginUri      string = "authn/login"
	LoginBody     string = `{"username":"%s","password":"%s"}`
	OkapiTokenKey string = "okapiToken"
	Username      string = "diku_admin"
	Password      string = "admin"

	OrderUri        string = "orders/composite-orders"
	OrderStatusKey  string = "workflowStatus"
	OrderOpenStatus string = "Open"
)

func ParseCsvAndOpenOrdersInBulk() {
	start := time.Now()
	slog.Info(CommandName, "Started, time", time.Now().Format(time.UnixDate))
	file, err := os.Open(FileName)
	if err != nil {
		slog.Error(CommandName, "os.Open error", "")
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var blockingChannel = make(chan int, PersistentHttpThreads)
	okapiToken := GetOkapiToken()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			slog.Info(CommandName, "Found an EOF", "")
			blockingChannel <- 1
			break
		}
		if err != nil {
			slog.Error(CommandName, "reader.Read error", "")
			panic(err)
		}
		if len(record) < 3 {
			slog.Error(CommandName, "reader.Read error", "")
			panic(fmt.Errorf("invalid record: %v", record))
		}
		orderId := record[0]
		if orderId == "po_id" {
			slog.Info(CommandName, "Ignoring header line", "")
			continue
		}
		blockingChannel <- 1
		go GetAndOpenOrder(okapiToken, orderId, &blockingChannel)
		slog.Info(CommandName, "Blocked count", len(blockingChannel))
	}
	select {
	case <-blockingChannel:
		slog.Info(CommandName, "Stopped, time", time.Now().Format(time.UnixDate))
		elapsed := time.Since(start)
		slog.Info(CommandName, "Elapsed, elapsed", elapsed)
	default:
	}
}

func GetOkapiToken() string {
	loginUrl := fmt.Sprintf("%s/%s", Hostname, LoginUri)
	loginHeaders := map[string]string{
		util.ContentTypeHeader: util.JsonContentType,
		util.XOkapiTenant:      DikuTenant,
	}
	loginBytes := []byte(fmt.Sprintf(LoginBody, Username, Password))
	return util.DoPostReturnMapStringInteface(CommandName, loginUrl, EnableDebug, loginBytes, loginHeaders)[OkapiTokenKey].(string)
}

func GetAndOpenOrder(okapiToken string, orderId string, requestBlockingChannel *chan int) {
	openOrderUrl := fmt.Sprintf("%s/%s/%s", Hostname, OrderUri, orderId)
	openOrderHeaders := map[string]string{
		util.ContentTypeHeader: util.JsonContentType,
		util.XOkapiTenant:      DikuTenant,
		util.XOkapiToken:       okapiToken,
	}
	compositeOrder := util.DoGetDecodeReturnMapStringInteface(CommandName, openOrderUrl, EnableDebug, true, openOrderHeaders)
	compositeOrder[OrderStatusKey] = OrderOpenStatus
	openOrderBytes, err := json.Marshal(compositeOrder)
	if err != nil {
		slog.Error(CommandName, "json.Marshal error", "")
		panic(err)
	}
	util.DoPutReturnNoContent(CommandName, openOrderUrl, EnableDebug, openOrderBytes, openOrderHeaders)
	<-*requestBlockingChannel
}
