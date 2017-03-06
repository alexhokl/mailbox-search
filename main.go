package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	targets, errTargets := getTargets()
	if errTargets != nil {
		fmt.Println(errTargets)
		return
	}
	isSentMode, errMode := isSentMode()
	if errMode != nil {
		fmt.Println(errMode)
		return
	}
	domain, errDomain := getDomain()
	if errDomain != nil {
		fmt.Println(errDomain)
		return
	}
	startDate, errStartDate := getStartDate()
	if errStartDate != nil {
		fmt.Println(errStartDate)
		return
	}
	endDate, errEndDate := getEndDate()
	if errEndDate != nil {
		fmt.Println(errEndDate)
		return
	}
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

func isSentMode() (bool, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_IS_SENT")
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(valStr)
}

func getTargets() ([]string, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_TARGETS")
	if err != nil {
		return nil, err
	}
	return strings.Split(valStr, ","), nil
}

func getDomain() (string, error) {
	return getEnvironmentVariable("MAILBOX_SEARCH_DOMAIN")
}

func getStartDate() (time.Time, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_START_DATE")
	if err != nil {
		return time.Now(), err
	}
	return getDate(valStr)
}

func getEndDate() (time.Time, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_END_DATE")
	if err != nil {
		return time.Now(), err
	}
	return getDate(valStr)
}

func getEnvironmentVariable(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", errors.New(fmt.Sprintf("Environment variable %s is not set"))
	}
	return val, nil
}

func getDate(valStr string) (time.Time, error) {
	const format = "2016-12-31T12:00:00Z"
	val, err := time.ParseInLocation(format, valStr, nil)
	if err != nil {
		return time.Now(), errors.New(fmt.Sprintf("Unable to parse date. Expected format %s", format))
	}
	return val, nil
}
