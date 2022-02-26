export function getTime(){
    const d = new Date();
    return d.getMinutes().toLocaleString(undefined, {minimumIntegerDigits: 2}) + ':' +
        d.getSeconds().toLocaleString(undefined, {minimumIntegerDigits: 2}) + '.' +
        d.getMilliseconds();
}