#!/usr/bin/env python
#
# Plot in-flight data over time, given the output of test-tool

import sys
import json
import arrow
from pprint import pprint
import matplotlib.pyplot as plt

raw = json.loads(open(sys.argv[1]).read())

events = []
i = 0
for record in raw['Records']:
    events.append({
        'packet_id': i,
        'type': 'request',
        'time': arrow.get(record['RequestAt']),
    })
    events.append({
        'packet_id': i,
        'type': 'proceed',
        'time': arrow.get(record['ProceedAt']),
    })
    events.append({
        'packet_id': i,
        'type': 'response',
        'time': arrow.get(record['ResponseAt']),
    })
    i += 1

events.sort(key=lambda e: e['time'])

start_time = events[0]['time']
end_time = events[-1]['time']

packet_size = raw['Records'][0]['Size']
bandwidth = raw['Bandwidth']
# Work out rtProp:
rtProp = 1e100
for record in raw['Records']:
    rtt = (arrow.get(record['ResponseAt']) - arrow.get(record['ProceedAt'])).total_seconds()
    if rtt < rtProp:
        rtProp = rtt

bdp_in_bytes = bandwidth * rtProp
bdp_in_packets = bdp_in_bytes / packet_size

pprint({
    "RtProp (seconds)": rtProp,
    "Bandwidth (bytes)": bandwidth,
    "BDP (bytes)": bdp_in_bytes,
    "BDP (packets)": bdp_in_packets,
})

# plot current packets on wire over time.
on_wire = 0
x = []
y = []
for e in events:
    if e['type'] == 'proceed':
        on_wire += 1
    elif e['type'] == 'response':
        on_wire -= 1
    x.append((e['time'] - start_time).total_seconds())
    y.append(on_wire)

plt.title("packets on wire over time")
plt.xlabel("time (seconds)")
plt.ylabel("packets")
plt.plot(x, y)

plt.show();
