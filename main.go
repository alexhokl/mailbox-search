package main

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	targets := []string{}
	isSentMode := true
	domain := ""
	startDate := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
	directories := os.Args[1:]
	for _, d := range directories {
		files, errDir := ioutil.ReadDir(d)
		if errDir != nil {
			continue
		}
		for _, f := range files {
			file, errOpen := os.Open(path.Join(d, f.Name()))
			if errOpen != nil {
				fmt.Println(errOpen)
				return
			}

			mailMsg, err := mail.ReadMessage(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			if isSentMode {
				if isTargetsTheOnlyAddress(mailMsg, domain, targets) && isInDateRange(mailMsg, startDate, endDate) {
					fullPath, _ := filepath.Abs(path.Join(d, f.Name()))
					fmt.Println(fullPath)
				}
			} else {
				if isInAddressList(mailMsg, targets) && isInDateRange(mailMsg, startDate, endDate) {
					fullPath, _ := filepath.Abs(path.Join(d, f.Name()))
					fmt.Println(fullPath)
				}
			}
			file.Close()
		}
	}
}

func isInAddressList(mailMessage *mail.Message, targets []string) bool {
	senders, _ := mailMessage.Header.AddressList("To")
	copies, _ := mailMessage.Header.AddressList("Cc")
	blindCopies, _ := mailMessage.Header.AddressList("Bcc")
	if hasTargets(senders, targets) {
		return true
	}
	if hasTargets(copies, targets) {
		return true
	}
	if hasTargets(blindCopies, targets) {
		return true
	}
	return false
}

func isTargetsTheOnlyAddress(mailMessage *mail.Message, domain string, targets []string) bool {
	for _, t := range targets {
		if isTargetTheOnlyAddress(mailMessage, domain, t) {
			return true
		}
	}
	return false
}

func isTargetTheOnlyAddress(mailMessage *mail.Message, domain string, target string) bool {
	senders, _ := mailMessage.Header.AddressList("To")
	copies, _ := mailMessage.Header.AddressList("Cc")
	blindCopies, _ := mailMessage.Header.AddressList("Bcc")

	addresses := []*mail.Address{}
	addresses = append(addresses, senders...)
	addresses = append(addresses, copies...)
	addresses = append(addresses, blindCopies...)

	domainAddresses := map[string]bool{}
	for _, a := range addresses {
		if strings.HasSuffix(a.Address, domain) {
			if _, exists := domainAddresses[a.Address]; !exists {
				domainAddresses[a.Address] = true
			}
		}
	}

	if len(domainAddresses) == 1 {
		for a, _ := range domainAddresses {
			if a == target {
				return true
			}
		}
	}

	return false
}

func hasTargets(list []*mail.Address, targets []string) bool {
	for _, a := range list {
		for _, t := range targets {
			if t == a.Address {
				return true
			}
		}
	}
	return false
}

func isInDateRange(mailMessage *mail.Message, startDate time.Time, endDate time.Time) bool {
	date, _ := mailMessage.Header.Date()
	return date.After(startDate) && date.Before(endDate)
}
