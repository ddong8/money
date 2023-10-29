package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

type ImageResizePayload struct {
	SourceURL string
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func getStock() string {
	stock_code := "AG2312"
	url := "http://hq.sinajs.cn/list=" + stock_code
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	req.Header.Add("Referer", "https://finance.sina.com.cn")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	bodyText := string(body)
	start := strings.Index(bodyText, "\"") + 1
	end := strings.LastIndex(bodyText, "\"") - 1
	dataStr := ""
	if end > start {
		dataStr = bodyText[start:end]
	}
	price := strings.Split(dataStr, ",")[7]
	return price
}

func sendPush(bodyStr string) {

	json := []byte(`{"body": "当前白银价格 ` + bodyStr + `","device_key": "StHezvE2w77GLuscNKRw75","title": "bleem", "badge": 1, "icon": "https://day.app/assets/images/avatar.jpg", "group": "test", "url": "https://mritd.com","category": "myNotificationCategory","sound": "minuet.caf"}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://api.day.app/push", body)
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}

func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	for {
		BodyStr := getStock()
		value, err := strconv.ParseFloat(BodyStr, 64)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
		if value >= 5890 {
			sendPush(BodyStr)
		}
		time.Sleep(10 * time.Second)
	}
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for struct
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
