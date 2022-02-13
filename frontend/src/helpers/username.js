export function getName() {
    let myName = localStorage.getItem('logjam_myName');
}

if (confirm(`Do you want to change ${myName} as your name?`)) {
    const newName = prompt('What is you name?');
    if (newName) {
        myName = newName;
        localStorage.setItem('logjam_myName', myName);
    }
}
if (myName === '' || !myName) {
    myName = makeid(20);
    try {
        localStorage.setItem('logjam_myName', myName);
    } catch (e) {
        console.log(e);
    }
}
