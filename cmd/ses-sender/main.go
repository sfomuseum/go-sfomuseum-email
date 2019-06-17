package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/gomail"
	"github.com/aaronland/gomail-ses"
	"log"
	"os"
	"path/filepath"
)

func main() {

	sender_addr := flag.String("sender", "", "...")
	recipient_addr := flag.String("recipient", "", "...")
	// subject := flag.String("subject", "", "...")

	ses_dsn := flag.String("ses-dsn", "credentials=session region=us-west-2", "...")

	flag.Parse()

	if *recipient_addr == "" {
		*recipient_addr = *sender_addr
	}

	s, err := ses.NewSESSender(*ses_dsn)

	if err != nil {
		log.Fatal(err)
	}

	m := gomail.NewMessage()

	m.SetHeader("From", *sender_addr)
	m.SetHeader("To", *recipient_addr)

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		fname := filepath.Base(abs_path)

		fh, err := os.Open(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		m.EmbedReader(fname, fh)

		img := fmt.Sprintf(`<p>Hello world</p><img src="cid:%s" alt="My image" /><p>WUB WUB WUB</p>`, fname)
		m.SetBody("text/html", img)
	}

	err = gomail.Send(s, m)

	if err != nil {
		log.Fatal(err)
	}

}
