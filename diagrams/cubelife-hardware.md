```mermaid
flowchart
subgraph cubelife
    subgraph main processor
        M4[Metro M4 Airlift]
    end
    subgraph bluetooth
        M4<--UART-->ESP32
    end
    subgraph panels
        Panel1[HUB75 #1]-->Panel2[HUB75 #2]
        Panel2-->Panel3[...]
        Panel3-->Panel6[HUB75 #6]
    end
    M4--SPI-->Panel1
    Panel6--SPI-->M4
end
```
