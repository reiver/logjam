const emailInput = document.querySelector('#email');
const submitBtn = document.getElementById("submit-btn");
const emailError = document.getElementById('email-error');

submitBtn.addEventListener('click',validateEmail);

function validateEmail(){

    const email = emailInput.value.trim();
    if (email === '') {
        console.log("empty");
        emailError.textContent = 'No email';
        emailError.classList.add('error');
      } else {
        console.log("not eempty");
        emailError.textContent = '';
        emailError.classList.remove('error');
      }
}

