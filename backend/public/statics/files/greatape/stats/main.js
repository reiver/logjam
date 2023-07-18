import { GraphGenerator } from './statsGraph.js';

window.addEventListener('DOMContentLoaded', () => {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const tabKey = urlParams.get('tabKey');

    if (tabKey) {
        setInterval(() => {
            const name = localStorage.getItem(`${tabKey}_name`);
            const data = localStorage.getItem(`${tabKey}_data`);

            const graphGenerator = new GraphGenerator();
            graphGenerator.parseJSONData(data);
            graphGenerator.createGraph();
            graphGenerator.displayGraph();
        }, 1000);
    }
});
