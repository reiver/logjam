function setMyName(){
    try {
        myName = localStorage.getItem('logjam_myName');
        document.getElementById('myName').value = myName;
    } catch (e) {
        console.log(e);
    }
    if (myName === '' || !myName) {
        myName = makeid(20);
        try {
            localStorage.setItem('logjam_myName', myName);
        } catch (e) {
            console.log(e);
        }
    }
}

function handleClick() {
    let newName = document.getElementById("myName").value;
    if (newName) {
        myName = newName;
        localStorage.setItem('logjam_myName', myName);
        console.log('myName=', myName);
    }
    document.getElementById("getName").style.display = "none";
    return false;
}
