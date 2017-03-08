# mailbox-search [![Build Status](https://travis-ci.org/alexhokl/mailbox-search.svg?branch=master)](https://travis-ci.org/alexhokl/mailbox-search)
CLI tool to filter mail message files by email addresses and dates

### Usage

To dump paths to all files of mail having targets in To, Cc, or Bcc. The
following example assumes `dovecot` and `maildir` are used.

```sh
export MAILBOX_SEARCH_MODE=normal
export MAILBOX_SEARCH_TARGETS=user.to.be.recovered.1@test.com,user.to.be.recovered2@test.com
export MAILBOX_SEARCH_DOMAIN=test.com
export MAILBOX_SEARCH_START_DATE=2016-01-01T00:00:00Z
export MAILBOX_SEARCH_END_DATE=2017-01-01T00:00:00Z

mailbox-search $(find . -type d -name "cur" -not -path "*/.Restored/*" -not -path "*/.spam/*" -not -path "*/.Sent/*" -not -path "*/.Trash/*" -not -path "*/.Junk*" -not -path "*/.Drafts/*" -not -path "*/.Archive/*" -not -path "*/.Infected*")
```

To dump paths to all files of mail having only one of the targets in To, Cc, Bcc. The following example assumes `dovecot` and `maildir` are used.

```sh
export MAILBOX_SEARCH_MODE=sent
export MAILBOX_SEARCH_TARGETS=user.to.be.recovered.1@test.com,user.to.be.recovered2@test.com
export MAILBOX_SEARCH_DOMAIN=test.com
export MAILBOX_SEARCH_START_DATE=2016-01-01T00:00:00Z
export MAILBOX_SEARCH_END_DATE=2017-01-01T00:00:00Z

mailbox-search $(find . -type d -name "cur" -path "*/.Sent/*")
```

To dump paths to all files of mail having targets in To, Cc, or Bcc and the
target address is malformed, for instance `hello@test.com` instead of
`<hello@test.com>` or `Hello <hello@test.com>`. The following example assumes `dovecot` and `maildir` are used.

```sh
export MAILBOX_SEARCH_MODE=normal_malform
export MAILBOX_SEARCH_TARGETS=user.to.be.recovered.1@test.com,user.to.be.recovered2@test.com
export MAILBOX_SEARCH_DOMAIN=test.com
export MAILBOX_SEARCH_START_DATE=2016-01-01T00:00:00Z
export MAILBOX_SEARCH_END_DATE=2017-01-01T00:00:00Z

mailbox-search $(find . -type d -name "cur" -not -path "*/.Restored/*" -not -path "*/.spam/*" -not -path "*/.Sent/*" -not -path "*/.Trash/*" -not -path "*/.Junk*" -not -path "*/.Drafts/*" -not -path "*/.Archive/*" -not -path "*/.Infected*")
```

To dump paths to all files of mail containing a specified subject.

```sh
export MAILBOX_SEARCH_MODE=subject
export MAILBOX_SEARCH_SUBJECT="Some subject of interest"

mailbox-search $(find . -type d -name "cur" -not -path "*/.Restored/*" -not -path "*/.spam/*" -not -path "*/.Sent/*" -not -path "*/.Trash/*" -not -path "*/.Junk*" -not -path "*/.Drafts/*" -not -path "*/.Archive/*" -not -path "*/.Infected*")
```

### Installation

##### Option 1

If you have Go installed, all you need is `go get github.com/alexhokl/mailbox-search`.

##### Option 2

Download binary from release page and put the binary in one of the directories
specified in `PATH` enviornment variable.

