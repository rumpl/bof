package daemon // import "github.com/rumpl/bof/daemon"

import (
	// Importing packages here only to make sure their init gets called and
	// therefore they register themselves to the logdriver factory.
	_ "github.com/rumpl/bof/daemon/logger/awslogs"
	_ "github.com/rumpl/bof/daemon/logger/fluentd"
	_ "github.com/rumpl/bof/daemon/logger/gcplogs"
	_ "github.com/rumpl/bof/daemon/logger/gelf"
	_ "github.com/rumpl/bof/daemon/logger/journald"
	_ "github.com/rumpl/bof/daemon/logger/jsonfilelog"
	_ "github.com/rumpl/bof/daemon/logger/local"
	_ "github.com/rumpl/bof/daemon/logger/logentries"
	_ "github.com/rumpl/bof/daemon/logger/loggerutils/cache"
	_ "github.com/rumpl/bof/daemon/logger/splunk"
	_ "github.com/rumpl/bof/daemon/logger/syslog"
)
