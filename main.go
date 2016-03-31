package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

var failed, passed uint64

func main() {

	watchs := flag.String("watch", "8.8.8.8,233.5.5.5", "nsserver that need to be watched")
	domains := flag.String("domain", "google.com,baidu.com", "domains that need to be watched")
	timeout := flag.Int("timeout", 5, "set query timeout in seconds")
	help := flag.Bool("help", false, "usages")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}

	logrus.SetFormatter(&logrus.TextFormatter{})

	if len(*watchs) <= 0 {
		logrus.Fatalln("need set -watch")
	}

	if len(*domains) <= 0 {
		logrus.Fatalln("need set -domain")
	}

	ws := strings.Split(*watchs, ",")
	ds := strings.Split(*domains, ",")

	for _, ns := range ws {
		for _, d := range ds {
			go func(d0, ns0 string, t0 int) {
				tick := time.Tick(time.Second)
				for {
					<-tick
					go watch(d0, ns0, time.Duration(t0))
				}
			}(d, ns, *timeout)
		}
	}

	go func() {
		tick := time.Tick(time.Second * 60)
		for {
			logrus.WithFields(logrus.Fields{
				"failed": atomic.LoadUint64(&failed),
				"passed": atomic.LoadUint64(&passed),
			}).Infoln("status")
			<-tick
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c
	// cleanup
	os.Exit(0)
}

func watch(target, server string, timeout time.Duration) {
	c := dns.Client{}
	c.DialTimeout = time.Second * timeout
	c.ReadTimeout = time.Second * timeout
	c.WriteTimeout = time.Second * timeout
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	r, _, err := c.Exchange(&m, server+":53")
	if err != nil || len(r.Answer) == 0 {
		fields := logrus.Fields{
			"error":  err,
			"target": target,
			"server": server,
		}
		if r != nil {
			fields["answers"] = len(r.Answer)
		}
		logrus.WithFields(fields).Warnln("dns failed to reponse")
		atomic.AddUint64(&failed, 1)
	} else {
		atomic.AddUint64(&passed, 1)
	}
}
