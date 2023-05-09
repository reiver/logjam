const emailInput = document.querySelector('#email');
const submitBtn = document.getElementById("submit-btn");
const emailError = document.getElementById('email-error');
const subText = document.getElementById("subtext-message");
const instructions = document.getElementById("instructions");
const emailBox = document.getElementById("email");
const emailIcon = document.querySelector('#email-icon');
const emailLabel = document.getElementById("email-label");

const errorText = "Wrong Email Address";
const errorIcon = '<img src="AttentionCircle.svg" alt="Error icon">';

//remove br
if(window.innerWidth<480){
  subText.innerHTML = subText.innerHTML.replace("<br>", " ");
}


emailInput.addEventListener('input', function() {
  // code to handle the change event
  const email = emailInput.value.trim();
 
  if(email!== '' && !validateEmailFormat(email)){
    displayError(true);
  }else{
    displayError(false);
  }
});

submitBtn.addEventListener('click',function(event) {
  event.preventDefault();
  validateEmail();
});

function validateEmail(){

    const email = emailInput.value.trim();

    if (email === '') {
        console.log("empty");
        displayError(true);
      } else if(!validateEmailFormat(email)){
        displayError(true);
      } 
      else {

        let metaData = "Testing Data";

        //save Email Address
        if(saveEmailAddress(email,metaData)){

            submitBtn.disabled = true;

            //hide error msg
            displayError(false);

            //remove extra elements
            hideExtraElements();

            //display success msg
            displayMessage();
        }else{
            alert("saving failed");
        }
      
      }
}

function validateEmailFormat(email) {
  // Regular expression to validate email format
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}


function displayError(display){
  if(display){
    //dislay it
    emailError.textContent = errorText;
    emailError.classList.add('error');
    emailInput.classList.add('error');
    emailIcon.innerHTML = errorIcon;
    emailIcon.classList.add('error');
    emailBox.style.borderColor = "#C00006";
    emailLabel.style.display = "block";
    emailLabel.classList.add('error');
  }else{
    //hide it

    emailError.textContent = '';
    emailError.classList.remove('error');
    emailInput.classList.remove('error');
    emailIcon.innerHTML = '';
    emailIcon.classList.remove('error');
    emailBox.style.borderColor = "#a8a8a8";
    emailLabel.style.display = "none";
    emailLabel.classList.remove('error');


  }
}

function displayMessage(){
  //display message
  if(window.innerWidth>480){
      subText.innerHTML = "Thanks For Being Interested in testing the Beta Version<br>of GreatApe Application, we'll notify you via email."
  }else{
      subText.innerHTML = "Thanks For Being Interested in testing the Beta Version of GreatApe Application, we'll notify you via email."
  }
}

function hideExtraElements(){
    submitBtn.textContent = "Submitted";
    submitBtn.style.backgroundColor = "#a8a8a8";
    submitBtn.style.color = "#FFFFFF";
    emailBox.style.display = "none";
    instructions.style.display = "none";
}


function saveEmailAddress(email,metaData){

    let data = {
      email: email,
      data: metaData
    };

    fetch('http://logjamlayoutdemo.group.video:8090/subscribe', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    })
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(data => {
      console.log('Response:', data);
    })
    .catch(error => {
      console.error('Error:', error);
    });

    return true;
}

