const BASE_URL = "http://127.0.0.1:8080";
const SIGN_FORM = document.getElementById('sign-form');
const USER_INPUT = document.getElementById('username');
const PASSWORD_INPUT = document.getElementById('password');
const ERROR_TEXT = document.getElementById('error-text');

let REMOTE_URL = "";
let CURRENT_TYPE = "in";
let NEXT_TYPE = undefined;

function render() {
  if (NEXT_TYPE) {
    CURRENT_TYPE = NEXT_TYPE;
  }
  const alt = CURRENT_TYPE === "in" ? "up" : "in";
  const altText = `Sign ${alt}`;

  const switchLink = document.getElementById("switch-link");
  switchLink.innerText = altText;

  const text = `Sign ${CURRENT_TYPE}`;
  document.title = text;

  const header = document.getElementById("h1");
  header.innerHTML = text;

  const button = document.getElementById("button");
  button.innerText = text;

  const endpoint = `sign${CURRENT_TYPE}`;
  REMOTE_URL = `${BASE_URL}/${endpoint}`;

  USER_INPUT.value = "";
  PASSWORD_INPUT.value = "";

  NEXT_TYPE = alt;
}

async function sign(event) {
  event.preventDefault();

  const username = USER_INPUT.value;
  const password = PASSWORD_INPUT.value;
  const user = JSON.stringify({
    'name': username,
    'password': password,
  });

  try {
    const response = await fetch(`${REMOTE_URL}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': 'index.html',
      },
      body: user,
    });

    if (response.ok) {
      response.text()
        .then(token => {
          sessionStorage.setItem('token', token);
          sessionStorage.setItem('userinfo', user);
          if (CURRENT_TYPE == "in") {
            window.location.replace('../index.html');
          } else {
            CURRENT_TYPE = "in"
            window.location.replace('sign.html');
          }
        });
    } else {
      ERROR_TEXT.textContent = `Sign${CURRENT_TYPE} failed. Please check your credentials.`;
    }
  } catch (err) {
    const ERROR_TEXT = document.getElementById('error-text');
    ERROR_TEXT.textContent = `Sign${CURRENT_TYPE} failed. Please try again later.`;
  }
}

function validateInput(regex, textInput, errorDiv) {
  const value = textInput.value;

  if (value === "" || value.match(regex)) {
    textInput.classList.remove("invalid-input");
    errorDiv.innerHTML = '';
  } else {
    textInput.classList.add("invalid-input");
    errorDiv.innerHTML = `input should match the following: ${regex}`;
  }
}

function validatePassword() {
  const reValidPassword = /^[a-z0-9]{3,10}$/;
  const value = PASSWORD_INPUT.value;
  const errorDiv = document.getElementById("error-password-div");

  validateInput(reValidPassword, PASSWORD_INPUT, errorDiv);
  isValidForm();
}

function validateUsername() {
  const reValidUsername = /^[a-z0-9]{3,10}$/;
  const value = USER_INPUT.value;
  const errorDiv = document.getElementById("error-username-div");

  validateInput(reValidUsername, USER_INPUT, errorDiv);
  isValidForm();
}

function isValidForm() {
  const userValue = USER_INPUT.value;
  const passValue = PASSWORD_INPUT.value;
  const eitherEmpty = userValue === "" || passValue === "";

  const invalidUser = USER_INPUT.classList.contains("invalid-input");
  const invalidPass = PASSWORD_INPUT.classList.contains("invalid-input");

  const button = document.getElementById("button");
  button.disabled = invalidUser || invalidPass || eitherEmpty;
}

USER_INPUT.value = "";
PASSWORD_INPUT.value = "";
isValidForm();
render();
