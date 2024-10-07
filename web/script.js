const BASE_URL = 'http://127.0.0.1:8080';
const CONTENT_DIV = document.getElementById('content');

async function fetchContent() {
  const response = await fetch(BASE_URL, {
    method: 'GET',
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
  });

  return new Promise((resolve, reject) => {
    if (response.ok) {
      return resolve(response.text());
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
  // TODO: check if token in sessionStorage is valid
  const response = await fetch(`${BASE_URL}`, {
    method: 'GET',
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
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
        window.location.replace('signin/signin.html');
      }
    },
      err => {
        window.location.replace('signin/signin.html');
      });
}

async function logout() {
  const response = await fetch("http://localhost:8080/logout", {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    }
  });

  if (response.ok) {
    response.text()
    .then(function(response) {
      CONTENT_DIV.innerHTML = response;
      window.location.replace('signin/signin.html');
    });
  } else {
    CONTENT_DIV.innerHTML = "error loggin out";
  }
}

redirectIfInvalidSession();
