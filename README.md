# SNMP Scan - Python scripts

This script is an adaptation for Observium. Because it use also the `add_device.php` script, this 'tool' is REAAAAALLLLLLYYYYYYYY slow (yep, slowest than adding manually device.... but required for autodiscovery on network).

For the future, I'll change this script with a Golang script (and with mannualy checking the SNMP port to reduce the number of `add_device` calls)