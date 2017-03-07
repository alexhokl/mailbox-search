package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type mode int

const (
	NORMAL mode = iota
	SENT
	NORMAL_MALFORM
)

func main() {
	targets, errTargets := getTargets()
	if errTargets != nil {
		fmt.Println(errTargets)
		return
	}
	processMode, errMode := getMode()
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
			err := processMailFile(d, f, processMode, domain, targets, startDate, endDate)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func processMailFile(directory string, file os.FileInfo, processMode mode, domain string, targets []string, startDate time.Time, endDate time.Time) error {
	fileReader, errOpen := os.Open(path.Join(directory, file.Name()))
	if errOpen != nil {
		return errOpen
	}
	defer fileReader.Close()

	mailMsg, err := mail.ReadMessage(fileReader)
	if err != nil {
		return err
	}

	switch processMode {
	case NORMAL:
		if isInAddressList(mailMsg, targets) && isInDateRange(mailMsg, startDate, endDate) {
			printPath(directory, file)
		}
	case SENT:
		if isTargetsTheOnlyAddress(mailMsg, domain, targets) && isInDateRange(mailMsg, startDate, endDate) {
			printPath(directory, file)
		}
	case NORMAL_MALFORM:
		if !isInAddressList(mailMsg, targets) && isInDateRange(mailMsg, startDate, endDate) && isContainAddress(mailMsg, targets) {
			printPath(directory, file)
		}
	default:
		return errors.New(fmt.Sprintf("Mode [%s] is not supported", processMode))
	}

	return nil
}

func isContainAddress(mailMessage *mail.Message, targets []string) bool {
	senders := mailMessage.Header.Get("To")
	copies := mailMessage.Header.Get("Cc")
	blindCopies := mailMessage.Header.Get("Bcc")
	if hasTargetsInString(senders, targets) {
		return true
	}
	if hasTargetsInString(copies, targets) {
		return true
	}
	if hasTargetsInString(blindCopies, targets) {
		return true
	}
	return false
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

func hasTargetsInString(strList string, targets []string) bool {
	for _, t := range targets {
		if strings.Contains(strList, t) {
			return true
		}
	}
	return false
}

func isInDateRange(mailMessage *mail.Message, startDate time.Time, endDate time.Time) bool {
	date, _ := mailMessage.Header.Date()
	return date.After(startDate) && date.Before(endDate)
}

func getMode() (mode, error) {
	valStr, err := getEnvironmentVariable("MAILBOX_SEARCH_MODE")
	if err != nil {
		return NORMAL, err
	}
	return parseMode(valStr)
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
		return "", errors.New(fmt.Sprintf("Environment variable %s is not set", name))
	}
	return val, nil
}

func getDate(valStr string) (time.Time, error) {
	val, err := time.ParseInLocation(time.RFC3339, valStr, nil)
	if err != nil {
		return time.Now(), errors.New(fmt.Sprintf("Unable to parse date. Expected format %s", time.RFC3339))
	}
	return val, nil
}

func parseMode(valStr string) (mode, error) {
	switch valStr {
	case "normal":
		return NORMAL, nil
	case "sent":
		return SENT, nil
	case "normal_malform":
		return NORMAL_MALFORM, nil
	default:
		return NORMAL, errors.New(fmt.Sprintf("Mode [%s] is not supported", valStr))
	}
}

func printPath(directory string, file os.FileInfo) {
	fullPath, _ := filepath.Abs(path.Join(directory, file.Name()))
	fmt.Println(fullPath)
}
