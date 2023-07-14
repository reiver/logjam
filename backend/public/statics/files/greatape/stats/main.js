import { GraphGenerator } from './statsGraph.js';

const graphGenerator = new GraphGenerator();
graphGenerator.readFile('data.txt', (fileContents) => {
    graphGenerator.parseJSONData(fileContents);
    graphGenerator.createGraph();
    graphGenerator.displayGraph();
});
