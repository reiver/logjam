
const DEBUG = true;

export const useLogger = () => {
    function getTime(){
        const d = new Date();
        return d.getMinutes().toLocaleString(undefined, {minimumIntegerDigits: 2}) + ':' +
        d.getSeconds().toLocaleString(undefined, {minimumIntegerDigits: 2}) + '.' +
            d.getMilliseconds();
    }
    
    function log(category: string, ...message: any) {
        if (DEBUG){
            console.log(`[${getTime()}] [${category}]`, ...message);
        }
    }

    function error(category: string, ...message: any) {
        if (DEBUG){
            console.error(`[${getTime()}] [${category}]`, ...message);
        }
    }

    return {log, error};
}
