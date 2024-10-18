const BASE_URL = 'http://127.0.0.1:8080';
const CONTENT_DIV = document.getElementById('content');
const USERNAME_SPAN = document.getElementById('username');

async function fetchContent(path) {
  const response = await fetch(`${BASE_URL}${path}`, {
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

function loadContent(path) {
  // handle error or something
  if (!path || path === '') {
    path = "/";
  }
  CONTENT_DIV.innerHTML = "loading...";
  const content = fetchContent(path)
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
  isLoggedIn()
    .then(ok => {
      if (ok) {
        // TODO this is the username encoded
        const token = sessionStorage.getItem("token");
        USERNAME_SPAN.innerHTML = atob(token);

        loadContent();
      } else {
        window.location.replace('sign/sign.html');
      }
    },
      err => {
        window.location.replace('sign/sign.html');
      });
}

async function logout() {
  const response = await fetch(`${BASE_URL}/logout`, {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    }
  });

  if (response.ok) {
    response.text()
    .then(function(response) {
      CONTENT_DIV.innerHTML = response;
      window.location.replace('sign/sign.html');
    });
  } else {
    CONTENT_DIV.innerHTML = "error loggin out";
  }
}

redirectIfInvalidSession();
