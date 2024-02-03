```mermaid
flowchart
subgraph bluetooth device
    subgraph Application MCU
    TG[TinyGo]
    end
    subgraph Wireless MCU
    ESP32
    end
    TG<--HCI protocol over UART-->ESP32[NINA-FW]
end
```
