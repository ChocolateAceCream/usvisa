package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"visa/global"
	"visa/step"
	"visa/utils"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

type USVisaJob struct {
	Config []func(*chromedp.ExecAllocator)
	Ctx    context.Context
	Queue  *ExecutionQueue
	cancel context.CancelFunc
}

func StartListener(l Listener) {
	l.Start()
}

func StopListener(l Listener) {
	l.Stop()
}

func (u *USVisaJob) Init() {
	dir, err := ioutil.TempDir("", "temp")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	u.Config = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", global.CONFIG.Chromedp.Headless),
		chromedp.Flag("ignore-certificate-errors", global.CONFIG.Chromedp.IgnoreCertificateErrors),
		// chromedp.Flag("window-size", "50,400"),
		chromedp.UserDataDir(dir),
	)

	ctx, _ := chromedp.NewExecAllocator(context.Background(), u.Config...)
	u.Ctx, u.cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	deadline := time.Now().Add(time.Duration(global.CONFIG.Visa.TimeInterval) * time.Minute)
	u.Queue = NewExecutionQueue(u.Ctx, u.cancel, deadline)
}

func (u *USVisaJob) Run(ctx context.Context) {
	global.LOGGER.Info("---new run-----")
	go u.Queue.Run()
	global.LOGGER.Info("---step1-----")
	u.Queue.Enqueue(StepFunction{Function: step.Login})
	global.LOGGER.Info("---step2-----")
	u.Queue.Enqueue(StepFunction{Function: step.Step2})
	global.LOGGER.Info("---step3-----")
	u.Queue.Enqueue(StepFunction{Function: step.Step3})
	global.LOGGER.Info("---step4-----")
	u.Queue.Enqueue(StepFunction{Function: step.Step4})
	global.LOGGER.Info("---step5-----")

	//add response listener
	l := RespListenerInit("https://ais.usvisa-info.com/en-ca/niv/schedule/53766549/appointment/days/95.json?appointments[expedite]=false", u.Ctx)

	StartListener(l)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				utils.RandomSleep(65, 75)
				u.Queue.Enqueue(StepFunction{Function: step.Step5})
			}
		}
	}()
	var targetDate string
	var targetTime string
	//thing get ugly from here
outerLoop:
	for {
		select {
		case <-ctx.Done():
			global.LOGGER.Info("this run deadline reached, exiting")
			return
		case <-u.Queue.Ctx.Done():
			global.LOGGER.Info("deadline reached, exiting")
			return
		case requestID := <-l.ReqIDChan:
			c := chromedp.FromContext(l.Ctx)
			newCtx := cdp.WithExecutor(l.Ctx, c.Target)
			resp, err := network.GetResponseBody(requestID).Do(newCtx)
			log.Println("resp?----", string(resp))
			global.LOGGER.Info("resp：", zap.String("resp", string(resp)))

			if err != nil {
				log.Println(requestID, err)
				return
			}
			var temp []AvailableDate
			if err := json.Unmarshal(resp, &temp); err == nil {
				if len(temp) > 0 {
					firstDate := temp[0].Date
					global.LOGGER.Info("first available date", zap.String("date", firstDate))
					if CheckDate(firstDate) {
						//continue next step
						targetDate = firstDate
						StopListener(l)
						break outerLoop
					}
				}
			}

			global.LOGGER.Info("---repeat step5-----")
		}
	}
	global.LOGGER.Info("targetDate", zap.String("targetDate", targetDate))

	//step6: get time

	//add response listener
	url := fmt.Sprintf(`https://ais.usvisa-info.com/en-ca/niv/schedule/53766549/appointment/times/95.json?date=%s&appointments[expedite]=false`, targetDate)
	l = RespListenerInit(url, u.Ctx)

	StartListener(l)

	u.Queue.Enqueue(StepFunction{FunctionWithArgs: step.Step6, Args: []interface{}{targetDate}})

	for requestID := range l.ReqIDChan {
		c := chromedp.FromContext(l.Ctx)
		newCtx := cdp.WithExecutor(l.Ctx, c.Target)
		resp, err := network.GetResponseBody(requestID).Do(newCtx)
		global.LOGGER.Info("time query resp?", zap.String("resp", string(resp)))
		if err != nil {
			log.Println(requestID, err)
			return
		}
		var temp AvailableTime
		if err := json.Unmarshal(resp, &temp); err == nil {
			if len(temp.AvailableTimes) > 0 {
				firstTime := temp.AvailableTimes[0]
				global.LOGGER.Info("first available time is", zap.String("time ", string(firstTime)))
				//continue next step
				targetTime = firstTime
				StopListener(l)
				break
			}
		}
	}
	global.LOGGER.Info("targetTime", zap.String("targetTime", targetTime))
	u.Queue.Enqueue(StepFunction{FunctionWithArgs: step.Step7, Args: []interface{}{targetTime}})
	u.Queue.Wait()
	global.LOGGER.Info(fmt.Sprint("appointment made at: %s %s\n", targetDate, targetTime))
}

func CheckDate(date string) bool {
	layout := "2006-01-02"
	r, err := time.Parse(layout, date)
	if err != nil {
		global.LOGGER.Error("time parse error", zap.Error(err))
	}
	deadline := time.Date(2025, time.January, 24, 0, 0, 0, 0, time.UTC)
	return r.Before(deadline)
}
