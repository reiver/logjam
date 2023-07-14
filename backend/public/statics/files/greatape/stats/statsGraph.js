export class GraphGenerator {
    constructor() {
        this.fileContents = null;
        this.jsonElements = [];
        this.audioJitter = [];
        this.audioPacketsLost = [];
        this.audioAvailableOutgoingBitrate = [];
        this.videoJitter = [];
        this.videoPacketsLost = [];
        this.videoAvailableOutgoingBitrate = [];
        this.videoFramesPerSecond = [];
        this.fig = null;
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
                const jitter = jsonData.jitter;
                const packetsLost = jsonData.packetsLost;
                const availableOutgoingBitrate =
                    jsonData.availableOutgoingBitrate;
                const framesPerSecond = jsonData.framesPerSecond;

                if (kind === 'audio') {
                    this.audioJitter.push(jitter);
                    this.audioPacketsLost.push(packetsLost);
                    this.audioAvailableOutgoingBitrate.push(
                        availableOutgoingBitrate
                    );
                } else if (kind === 'video') {
                    this.videoJitter.push(jitter);
                    this.videoPacketsLost.push(packetsLost);
                    this.videoAvailableOutgoingBitrate.push(
                        availableOutgoingBitrate
                    );
                    this.videoFramesPerSecond.push(framesPerSecond);
                }
            } catch (error) {
                console.log(`Error parsing JSON: ${error}`);
            }
        });
    }

    createGraph() {
        this.fig = {
            data: [],
            layout: {
                grid: { rows: 1, columns: 4, pattern: 'free' },
                showlegend: false,
                subplotTitles: [
                    'Jitter',
                    'Packets Lost',
                    'Available Outgoing Bitrate',
                    'Frames Per Second',
                ],
            },
        };

        this.fig.data.push({
            y: this.audioJitter,
            name: 'Audio Jitter',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y',
            subplot: 'xy',
        });

        this.fig.data.push({
            y: this.audioPacketsLost,
            name: 'Audio Packets Lost',
            type: 'box',
            xaxis: 'x2',
            yaxis: 'y',
            subplot: 'xy2',
        });

        this.fig.data.push({
            y: this.audioAvailableOutgoingBitrate,
            name: 'Audio Bitrate',
            type: 'box',
            xaxis: 'x3',
            yaxis: 'y2',
            subplot: 'xy3',
        });

        this.fig.data.push({
            y: this.videoAvailableOutgoingBitrate,
            name: 'Video Bitrate',
            type: 'box',
            xaxis: 'x3',
            yaxis: 'y2',
            subplot: 'xy3',
        });

        this.fig.data.push({
            y: this.videoJitter,
            name: 'Video Jitter',
            type: 'box',
            xaxis: 'x',
            yaxis: 'y2',
            subplot: 'xy4',
        });

        this.fig.data.push({
            y: this.videoPacketsLost,
            name: 'Video Packets Lost',
            type: 'box',
            xaxis: 'x2',
            yaxis: 'y2',
            subplot: 'xy5',
        });

        this.fig.data.push({
            y: this.videoFramesPerSecond,
            name: 'Video Frames Per Second',
            type: 'box',
            xaxis: 'x4',
            yaxis: 'y2',
            subplot: 'xy7',
        });
    }

    displayGraph() {
        const config = { responsive: true };
        Plotly.newPlot(
            'graph-container',
            this.fig.data,
            this.fig.layout,
            config
        );
    }
}
