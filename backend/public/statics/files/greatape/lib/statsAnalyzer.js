export class StatsAnalyzer {
    constructor() {
        this.networkMeasurementData = {
            rtt: [],
            jitter: [],
            packetsLost: [],
            outgoingBitrate: [],
        };
        // this.packetLossThreshold = 50; // 10% packet loss
        // this.lowBandwidthThreshold = 300000; // in bits per second
        // this.latencyThreshold = 50; // in milliseconds
        // this.jitterThreshold = 10; // in milliseconds (example threshold, adjust as needed)
        this.logsCallback = null;
        this.decisionCallback = null;
        this.time = 0;

        const queryParams = new URLSearchParams(window.location.search);
        this.packetLossThreshold = queryParams.get('pcl');
        this.lowBandwidthThreshold = queryParams.get('btr');
        this.jitterThreshold = queryParams.get('jtr');
        this.latencyThreshold = queryParams.get('rtt');

        console.log('pcl: ', this.packetLossThreshold);
        console.log('rtt: ', this.latencyThreshold);
        console.log('jtr: ', this.jitterThreshold);
        console.log('brt: ', this.lowBandwidthThreshold);
    }

    registerCallbacks(logsCallback, decisionCallback) {
        this.logsCallback = logsCallback;
        this.decisionCallback = decisionCallback;
    }

    analyzeStatsReport(report) {
        const jsonData = report[1];
        const kind = jsonData.kind;
        const type = jsonData.type;

        let rtt = 0;
        let jitter = 0;
        let packetsLost = 0;
        let bitrate = jsonData.availableOutgoingBitrate || 0;

        if (type === 'candidate-pair') {
            rtt = jsonData.totalRoundTripTime;
        }

        if (kind === 'video') {
            jitter = jsonData.jitter;
            packetsLost = jsonData.packetsLost;
        }

        // Save the parameters to the networkMeasurementData object
        if (rtt > 0) {
            this.networkMeasurementData.rtt.push(rtt);
        }
        if (jitter > 0) {
            this.networkMeasurementData.jitter.push(jitter);
        }
        if (packetsLost > 0) {
            this.networkMeasurementData.packetsLost.push(packetsLost);
        }
        if (bitrate > 0) {
            this.networkMeasurementData.outgoingBitrate.push(bitrate);
        }

        // const arrayLimit = 20;
        // Perform reconnection check after a few minutes
        // if (
        //     (this.networkMeasurementData.rtt.length >= arrayLimit &&
        //         this.networkMeasurementData.jitter.length >= arrayLimit &&
        //         this.networkMeasurementData.outgoingBitrate.length >=
        //             arrayLimit) ||
        //     this.networkMeasurementData.packetsLost.length >= arrayLimit
        // )
        if (this.time === 20) {
            this.time = 0;

            const avgRTT = this.findMedian(this.networkMeasurementData.rtt);
            this.logsCallback(`AvgRTT: `, avgRTT);

            const avgJitter = this.findMedian(
                this.networkMeasurementData.jitter
            );
            this.logsCallback(`AvgJitter: `, avgJitter);

            const avgPacketsLost = this.findMedian(
                this.networkMeasurementData.packetsLost
            );
            this.logsCallback(`AvgPCLOss: `, avgPacketsLost);

            const avgOutgoingBitrate = this.findMedian(
                this.networkMeasurementData.outgoingBitrate
            );
            this.logsCallback(`AvgBitrate: `, avgOutgoingBitrate);

            if (
                avgRTT > this.latencyThreshold ||
                avgPacketsLost > this.packetLossThreshold ||
                avgJitter > this.jitterThreshold ||
                avgOutgoingBitrate < this.lowBandwidthThreshold
            ) {
                this.logsCallback(
                    `rtt: ${avgRTT} & thresh: ${this.latencyThreshold} \n 
                    pcLoss: ${avgPacketsLost} & thresh: ${this.packetLossThreshold} \n
                    jitter: ${avgJitter} & thresh: ${this.jitterThreshold} \n
                    bitrate: ${avgOutgoingBitrate} & thresh: ${this.lowBandwidthThreshold}`
                );

                if (this.decisionCallback) {
                    this.decisionCallback();
                }
            }

            // Clear the stored data after the reconnection check
            //this.resetNetworkMeasurementData();
        }
    }

    findMedian(arr) {
        const sortedArr = arr.slice().sort((a, b) => a - b);
        const mid = Math.floor(sortedArr.length / 2);
        return sortedArr.length % 2 === 0
            ? (sortedArr[mid - 1] + sortedArr[mid]) / 2
            : sortedArr[mid];
    }

    // Function to calculate the average of an array of numbers
    calculateAverage(arr) {
        // this.updateTheStatus(`calculateAvg `, arr);
        // Filter out null and undefined values from the array
        const validValues = arr.filter(
            (value) => value !== null && value !== undefined && value !== 0
        );

        // If there are no valid values, return 0 to avoid division by zero
        if (validValues.length === 0) return 0;

        // Calculate the sum of valid values and return the average
        const sum = validValues.reduce((total, value) => total + value, 0);
        return sum / validValues.length;
    }
    // Helper method to reset network measurement data
    resetNetworkMeasurementData() {
        this.networkMeasurementData = {
            rtt: [],
            jitter: [],
            packetsLost: [],
            outgoingBitrate: [],
        };
    }
}
