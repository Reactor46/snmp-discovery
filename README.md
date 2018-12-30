# SNMP Scan - Python scripts

This script is an adaptation for Observium. Because it uses the `add_device.php` script, this 'tool' is REAAAAALLLLLLYYYYYYYY slow (yep, slower than adding manally device.... but it's required for autodiscovery on network).

For the future, I'll change this script with a Golang script (and with manual checking of SNMP port to reduce the number of `add_device` calls)
