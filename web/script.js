const BASE_URL = 'http://127.0.0.1:8080';
const CONTENT_DIV = document.getElementById('content');

async function fetchContent() {
  const response = await fetch(BASE_URL, {
    method: 'GET',
  });

  return new Promise((resolve, reject) => {
    if (response.ok) {
      return resolve(response.text);
    }
    reject('auth failed?');
  });
}

function loadContent() {
  // handle error or something
  const content = fetchContent()
    .then(content => {
      CONTENT_DIV.innerHTML = content;
    },
    err => {
      console.error("some error:", err);
    });
}

async function isLoggedIn() {
  // TODO: check if token in localStorage is valid
  const response = await fetch(`${BASE_URL}`, {
    method: 'GET',
  });

  console.log(response);
  return new Promise((resolve, reject) => {
    return resolve(response.ok);
  });
}

function redirectIfInvalidSession() {
  CONTENT_DIV.innerHTML = "loading...";
  isLoggedIn()
  .then(ok => {
    if (ok) {
      loadContent();
    } else {
      window.location.replace('login/login.html');
    }
  },
  err => {
    window.location.replace('login/login.html');
  });
}

redirectIfInvalidSession();
