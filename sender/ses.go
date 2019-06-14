package sender

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html
// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-email-addresses-procedure.html

// https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.VerifyEmailIdentity
// https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.WaitUntilIdentityExists
// https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.ListIdentities
// https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.DeleteIdentity

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/send-personalized-email-api.html
// https://docs.aws.amazon.com/ses/latest/APIReference/API_CreateCustomVerificationEmailTemplate.html
// https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.CreateCustomVerificationEmailTemplate

// https://us-west-2.console.aws.amazon.com/ses/home?region=us-west-2#smtp-settings:

import (
	"bufio"
	"bytes"
	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/go-gomail/gomail"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"io"
	_ "log"
)

type SESSender struct {
	gomail.Sender
	service *ses.SES
}

func NewSESSender(dsn string) (gomail.Sender, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)

	s := SESSender{
		service: svc,
	}

	return &s, nil
}

func (s *SESSender) Send(from string, to []string, msg io.WriterTo) error {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	_, err := msg.WriteTo(wr)

	if err != nil {
		return err
	}

	wr.Flush()

	raw_msg := &ses.RawMessage{
		Data: buf.Bytes(),
	}
	
	for _, recipient := range to {

		err := s.sendMessage(from, recipient, raw_msg)

		if err != nil {
			return nil
		}
	}

	return nil
}

func (s *SESSender) sendMessage(sender string, recipient string, msg *ses.RawMessage) error {
	
	req := &ses.SendRawEmailInput{
		RawMessage: msg,
	}
	
	_, err := s.service.SendRawEmail(req)

	if err != nil {
		return err
	}

	return nil
}
