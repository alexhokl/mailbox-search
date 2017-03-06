# mailbox-search [![Build Status](https://travis-ci.org/alexhokl/mailbox-search.svg?branch=master)](https://travis-ci.org/alexhokl/mailbox-search)
CLI tool to filter mail message files by email addresses and dates

### Usage

To dump paths to all files of mail having targets in To, Cc, or Bcc. The
following example assumes `dovecot` and `maildir` are used.

```sh
export MAILBOX_SEARCH_IS_SENT=false
export MAILBOX_SEARCH_TARGETS=user.to.be.recovered.1@test.com,user.to.be.recovered2@test.com
export MAILBOX_SEARCH_DOMAIN=test.com
export MAILBOX_START_DATE=2016-01-01T00:00:00Z
export MAILBOX_END_DATE=2016-01-01T00:00:00Z

mailbox-search $(find . -type d -name "cur" -not -path "*/.Restored/*" -not -path "*/.spam/*" -not -path "*/.Sent/*" -not -path "*/.Trash/*" -not -path "*/.Junk*" -not -path "*/.Drafts/*" -not -path "*/.Archive/*" -not -path "*/.Infected*")
```

To dump paths to all files of mail having only one of the targets in To, Cc, Bcc. The following example assumes `dovecot` and `maildir` are used.

```sh
export MAILBOX_SEARCH_IS_SENT=true
export MAILBOX_SEARCH_TARGETS=user.to.be.recovered.1@test.com,user.to.be.recovered2@test.com
export MAILBOX_SEARCH_DOMAIN=test.com
export MAILBOX_START_DATE=2016-01-01T00:00:00Z
export MAILBOX_END_DATE=2016-01-01T00:00:00Z

mailbox-search $(find . -type d -name "cur" -path "*/.Sent/*")
```

