export class GraphGenerator {
    constructor() {
        this.fileContents = null;
        this.jsonElements = [];
        this.audioJitter = [];
        this.audioPacketsLost = [];
        this.videoJitter = [];
        this.videoPacketsLost = [];
        this.availableOutgoingBitrate = [];
        this.videoFramesPerSecond = [];
        this.outbound_fps = [];
        this.inbound_fps = [];
        this.remote_inbound_fps = [];
        this.source_fps = [];

        this.jitterFig = null;
        this.fpsFig = null;
        this.packetLossFig = null;
        this.bitrateFig = null;
    }

    readFile(filePath, callback) {
        const xhr = new XMLHttpRequest();
        xhr.open('GET', filePath, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {
                callback(xhr.responseText);
            }
        };
        xhr.send();
    }

    parseJSONData(fileContents) {
        this.fileContents = fileContents;
        this.jsonElements = this.fileContents.match(/\{.*?\}/gms);

        console.log('JSONELEMNT:', this.jsonElements);

        this.jsonElements.forEach((element) => {
            try {
                const jsonData = JSON.parse(element);
                const kind = jsonData.kind;
                const type = jsonData.type;
                const jitter = jsonData.jitter;
                const packetsLost = jsonData.packetsLost;
                const framesPerSecond = jsonData.framesPerSecond;

                this.availableOutgoingBitrate.push(
                    jsonData.availableOutgoingBitrate
                );

                if (kind === 'audio') {
                    this.audioJitter.push(jitter);
                    this.audioPacketsLost.push(packetsLost);
                } else if (kind === 'video') {
                    this.videoJitter.push(jitter);
                    this.videoPacketsLost.push(packetsLost);

                    if (type === 'outbound-rtp') {
                        this.outbound_fps.push(framesPerSecond);
                    } else if (type === 'media-source') {
                        this.source_fps.push(framesPerSecond);
                    } else if (type === 'inbound-rtp') {
                        this.inbound_fps.push(framesPerSecond);
                    } else if (type === 'remote-inbound-rtp') {
                        this.remote_inbound_fps.push(framesPerSecond);
                    }
                }
            } catch (error) {
                console.log(`Error parsing JSON: ${error}`);
            }
        });
    }

    createJitterGraph() {
        this.jitterFig = {
            data: [],
            layout: {
                grid: { rows: 1, columns: 2, pattern: 'free' },
                showlegend: false,
            },
        };

        this.jitterFig.data.push({
            y: this.audioJitter,
            name: 'Audio Jitter',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y',
            subplot: 'xy',
        });

        this.jitterFig.data.push({
            y: this.videoJitter,
            name: 'Video Jitter',
            type: 'box',
            xaxis: 'x2',
            yaxis: 'y',
            subplot: 'xy2',
        });
    }

    createFPSGraph() {
        this.fpsFig = {
            data: [],
            layout: {
                grid: { rows: 1, columns: 2, pattern: 'free' },
                showlegend: false,
            },
        };

        this.fpsFig.data.push({
            y: this.inbound_fps,
            name: 'Inbound FPS',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y',
            subplot: 'xy',
        });

        this.fpsFig.data.push({
            y: this.source_fps,
            name: 'Source FPS',
            type: 'box',
            xaxis: 'x2',
            yaxis: 'y',
            subplot: 'xy2',
        });
    }

    createPacketLossGraph() {
        this.packetLossFig = {
            data: [],
            layout: {
                grid: { rows: 1, columns: 2, pattern: 'free' },
                showlegend: false,
            },
        };

        this.packetLossFig.data.push({
            y: this.audioPacketsLost,
            name: 'Audio Packets Loss',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y',
            subplot: 'xy',
        });

        this.packetLossFig.data.push({
            y: this.videoPacketsLost,
            name: 'Video Packets Loss',
            type: 'box',
            xaxis: 'x2',
            yaxis: 'y',
            subplot: 'xy2',
        });
    }

    createBitrateGraph() {
        this.bitrateFig = {
            data: [],
            layout: {
                grid: { rows: 1, columns: 2, pattern: 'free' },
                showlegend: false,
            },
        };

        this.bitrateFig.data.push({
            y: this.availableOutgoingBitrate,
            name: 'Available Outgoing Bitrate',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y',
            subplot: 'xy',
        });
    }

    createGraph() {
        this.createJitterGraph();
        this.createFPSGraph();
        this.createPacketLossGraph();
        this.createBitrateGraph();
    }

    displayGraph() {
        const config = { responsive: true };
        Plotly.newPlot(
            'jitter',
            this.jitterFig.data,
            this.jitterFig.layout,
            config
        );

        Plotly.newPlot('fps', this.fpsFig.data, this.fpsFig.layout, config);

        Plotly.newPlot(
            'packetLoss',
            this.packetLossFig.data,
            this.packetLossFig.layout,
            config
        );

        Plotly.newPlot(
            'bitrate',
            this.bitrateFig.data,
            this.bitrateFig.layout,
            config
        );
    }
}
