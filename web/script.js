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

  return new Promise((resolve, reject) => {
    return resolve(response.ok);
  });
}

function userinfo() {
  const userinfo = sessionStorage.getItem("userinfo");
  return JSON.parse(userinfo);
}

function redirectIfInvalidSession() {
  const redirectPath = "sign/sign.html";
  const user = userinfo();
  if (!user) {
    window.location.replace(redirectPath);
  }

  USERNAME_SPAN.innerHTML = user.name;
  isLoggedIn()
    .then(
      ok => {
        if (!ok) {
          window.location.replace(redirectPath);
        }
      },
      err => {
        window.location.replace(redirectPath);
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

async function listPosts() {
  const response = await fetch(`${BASE_URL}/posts`, {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    }
  });

  if (! response.ok) {
    return;
  }
  console.log("response from listPosts:", response);

  const postListDiv = document.getElementById("post-list");
  const postList = await response.json();
  for (const post of postList) {
    const link = document.createElement("a");
    link.href = "#";
    link.innerHTML = post.title;

    postListDiv.appendChild(link);
  }
}

redirectIfInvalidSession();
listPosts();
