```mermaid
flowchart LR
subgraph flightbadge
    subgraph main processor
        M4[Pybadge]
    end
    subgraph bluetooth processor
        M4<--UART-->ESP32
    end
end
subgraph drone
    ESP32<--BLE--->D[Parrot Minidrone]
end
```
