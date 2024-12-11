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
	CommandName string = "Orders"
	FileName    string = "./import/data.csv"

	Tenant string = "diku"

	LoginUri      string = "authn/login"
	LoginBody     string = `{"username":"%s","password":"%s"}`
	OkapiTokenKey string = "okapiToken"
	Username      string = "diku_admin"
	Password      string = "admin"

	OrderUri        string = "orders/composite-orders"
	OrderStatusKey  string = "workflowStatus"
	OrderOpenStatus string = "Open"
)

func ParseCsvAndOpenOrdersInBulk(gatewayHostname string, enableDebug bool, threadCount int) {
	start := time.Now()
	slog.Info(CommandName, "Started, time", time.Now().Format(time.UnixDate))

	file, err := os.Open(FileName)
	if err != nil {
		slog.Error(CommandName, "os.Open error", "")
		panic(err)
	}
	defer file.Close()

	okapiToken := GetOkapiToken(gatewayHostname, enableDebug)

	var blockingChannel = make(chan int, threadCount)
	reader := csv.NewReader(file)
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

		go GetAndOpenOrder(gatewayHostname, enableDebug, okapiToken, orderId, &blockingChannel)

		slog.Info(CommandName, "Blocked count", len(blockingChannel))
	}

	BlockSelectMain(blockingChannel, start)
}

func BlockSelectMain(blockingChannel chan int, start time.Time) {
	select {
	case <-blockingChannel:
		slog.Info(CommandName, "Stopped, time", time.Now().Format(time.UnixDate))
		elapsed := time.Since(start)
		slog.Info(CommandName, "Elapsed, duration", elapsed)
	default:
	}
}

func GetOkapiToken(gatewayHostname string, enableDebug bool) string {
	loginUrl := fmt.Sprintf("%s/%s", gatewayHostname, LoginUri)

	loginHeaders := map[string]string{
		util.ContentTypeHeader: util.JsonContentType,
		util.XOkapiTenant:      Tenant,
	}

	loginBytes := []byte(fmt.Sprintf(LoginBody, Username, Password))

	return util.DoPostReturnMapStringInteface(CommandName, loginUrl, enableDebug, loginBytes, loginHeaders)[OkapiTokenKey].(string)
}

func GetAndOpenOrder(gatewayHostname string, enableDebug bool, okapiToken string, orderId string, blockingChannel *chan int) {
	openOrderUrl := fmt.Sprintf("%s/%s/%s", gatewayHostname, OrderUri, orderId)

	openOrderHeaders := map[string]string{
		util.ContentTypeHeader: util.JsonContentType,
		util.XOkapiTenant:      Tenant,
		util.XOkapiToken:       okapiToken,
	}

	compositeOrder := util.DoGetDecodeReturnMapStringInteface(CommandName, openOrderUrl, enableDebug, true, openOrderHeaders)
	compositeOrder[OrderStatusKey] = OrderOpenStatus
	openOrderBytes, err := json.Marshal(compositeOrder)

	if err != nil {
		slog.Error(CommandName, "json.Marshal error", "")
		panic(err)
	}

	util.DoPutReturnNoContent(CommandName, openOrderUrl, enableDebug, openOrderBytes, openOrderHeaders)
	<-*blockingChannel
}
