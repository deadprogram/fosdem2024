```mermaid
flowchart LR
subgraph flightbadge
    subgraph main processor
        M4[Pybadge]
    end
    subgraph wifi processor
        M4<--SPI-->ESP32
    end
end
subgraph drone
    ESP32<--WiFi--->D[DJI Tello]
end
```
