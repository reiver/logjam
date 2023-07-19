import { GraphGenerator } from './statsGraph.js';

window.addEventListener('message', function (event) {
    var receivedData = event.data;

    const graphGenerator = new GraphGenerator();
    graphGenerator.parseJSONData(receivedData);
    graphGenerator.createGraph();
    graphGenerator.displayGraph();
});
