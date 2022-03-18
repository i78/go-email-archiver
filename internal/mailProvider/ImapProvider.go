package mailProvider

import (
	"context"
	"email-archiver-cli/internal/util"
	"email-archiver-cli/proto/emailarchver/config"
	"email-archiver-cli/proto/emailarchver/email"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/mail"
	"time"
)

type IMAPProvider struct {
	context         context.Context
	config          config.FetchMailConfig
	newMailCallback func(mail email.EMail) error
	log             *log.Entry
}

func NewIMAPProvider(ctx context.Context, config config.FetchMailConfig, newEMailCallback func(mail email.EMail) error) {
	provider := IMAPProvider{
		context:         ctx,
		config:          config,
		newMailCallback: newEMailCallback,
		log:             util.PrepareLogging(log.Fields{"source": "IMAPProvider"}),
	}

	provider.run()
}

func (i *IMAPProvider) run() {
	log.Info("Starting IMAP Provider")
	for {
		select {
		case <-time.After(2 * time.Second):
			log.Info("Checking mail")
			i.fetch()
		case <-i.context.Done():
			log.Info("Shutting down IMAPProvider")
		}
	}
}

func (i *IMAPProvider) fetch() {
	i.log.Debug("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(i.config.Server, nil)
	if err != nil {
		log.Fatal(err)
	}
	i.log.Debug("Connected")

	defer c.Logout()

	// Login
	if err := c.Login(i.config.Username, i.config.Password); err != nil {
		i.log.Fatal(err)
	}
	i.log.Debug("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	i.log.Debug("Mailboxes:")
	for m := range mailboxes {
		i.log.Debug("* " + m.Name)
	}

	if err := <-done; err != nil {
		i.log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		i.log.Fatal(err)
	}
	i.log.Debug("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 8 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 8
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	section := &imap.BodySectionName{}

	messages := make(chan *imap.Message, 10)
	// done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	for msg := range messages {

		r := msg.GetBody(section)
		if r == nil {
			i.log.Fatal("Server didn't returned message body")
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			i.log.Fatal(err)
		}

		email := mapEMail(m, func() ([]byte, error) {
			return ioutil.ReadAll(m.Body)
		})

		if archiveError := i.newMailCallback(email); archiveError != nil {
			i.log.WithFields(log.Fields{
				"msgId":         email.MessageId,
				"innerError":    archiveError,
				"derivedAction": "retry,ignore"}).Error("archiver reported error")
		} else {
			i.log.WithFields(log.Fields{
				"msgId":         email.MessageId,
				"derivedAction": "delete-from-server"}).Info("archiver completed")
		}

	}

	if err := <-done; err != nil {
		i.log.Fatal(err)
	}
}
