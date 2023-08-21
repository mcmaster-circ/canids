# Install zeek
echo 'deb http://download.opensuse.org/repositories/security:/zeek/Debian_11/ /' | sudo tee /etc/apt/sources.list.d/security:zeek.list
curl -fsSL https://download.opensuse.org/repositories/security:zeek/Debian_11/Release.key | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/security_zeek.gpg > /dev/null
sudo apt update
sudo apt install zeek

 

# Install custom zeek config, specifying 1hr log rotation interval, and 1 day log purge interval.
cat > /opt/zeek/etc/zeekctl.cfg << EOF
## Global ZeekControl configuration file.

 

###############################################
# Mail Options

 

# Recipient address for all emails sent out by Zeek and ZeekControl.
MailTo = root@localhost

 

# Mail connection summary reports each log rotation interval.  A value of 1
# means mail connection summaries, and a value of 0 means do not mail
# connection summaries.  This option has no effect if the trace-summary
# script is not available.
MailConnectionSummary = 1

 

# Lower threshold (in percentage of disk space) for space available on the
# disk that holds SpoolDir. If less space is available, "zeekctl cron" starts
# sending out warning emails.  A value of 0 disables this feature.
MinDiskSpace = 5

 

# Send mail when "zeekctl cron" notices the availability of a host in the
# cluster to have changed.  A value of 1 means send mail when a host status
# changes, and a value of 0 means do not send mail.
MailHostUpDown = 1

 

###############################################
# Logging Options

 

# Rotation interval in seconds for log files on manager (or standalone) node.
# A value of 0 disables log rotation.
LogRotationInterval = 3600

 

# Expiration interval for archived log files in LogDir.  Files older than this
# will be deleted by "zeekctl cron".  The interval is an integer followed by
# one of these time units:  day, hr, min.  A value of 0 means that logs
# never expire.
LogExpireInterval = 1day

 

# Enable ZeekControl to write statistics to the stats.log file.  A value of 1
# means write to stats.log, and a value of 0 means do not write to stats.log.
StatsLogEnable = 1

 

# Number of days that entries in the stats.log file are kept.  Entries older
# than this many days will be removed by "zeekctl cron".  A value of 0 means
# that entries never expire.
StatsLogExpireInterval = 0

 

###############################################
# Other Options

 

# Show all output of the zeekctl status command.  If set to 1, then all output
# is shown.  If set to 0, then zeekctl status will not collect or show the peer
# information (and the command will run faster).
StatusCmdShowAll = 0

 

# Number of days that crash directories are kept.  Crash directories older
# than this many days will be removed by "zeekctl cron".  A value of 0 means
# that crash directories never expire.
CrashExpireInterval = 0

 

# Site-specific policy script to load. Zeek will look for this in
# $PREFIX/share/zeek/site. A default local.zeek comes preinstalled
# and can be customized as desired.
SitePolicyScripts = local.zeek

 

# Location of the log directory where log files will be archived each rotation
# interval.
LogDir = /opt/zeek/logs

 

# Location of the spool directory where files and data that are currently being
# written are stored.
SpoolDir = /opt/zeek/spool

 

# Location of the directory in which the databases for Broker datastore backed
# Zeek tables are stored.
BrokerDBDir = /opt/zeek/spool/brokerstore

 

# Location of other configuration files that can be used to customize
# ZeekControl operation (e.g. local networks, nodes).
CfgDir = /opt/zeek/etc
EOF

 

# Create the systemd zeek service
cat > /etc/systemd/system/canids-zeek.service << EOF
[Unit]
Description=Zeek daemon to provide logs for canids ingestion
After=network.target

 

[Service]
Type=oneshot
RemainAfterExit=yes
User=root
Group=root
WorkingDirectory=/opt/zeek/
Environment="PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/opt/zeek/bin"
ExecStart=/opt/zeek/bin/zeekctl cron enable
ExecStart=/opt/zeek/bin/zeekctl deploy
ExecStop=/opt/zeek/bin/zeekctl stop
ExecStop=/opt/zeek/bin/zeekctl cron disable

 

[Install]
WantedBy=multi-user.target
EOF

 

systemctl daemon-reload # Load the zeek service
systemctl enable --now canids-zeek # Enable and start the zeek service

 

# Create the crontab for zeek
cat > /etc/cron.d/zeek << EOF
# crontab entries for canids/zeek
SHELL=/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/opt/zeek/bin

 

*/5 * * * * root zeekctl cron
EOF