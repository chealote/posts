const USERNAME_SPAN = document.getElementById("username-span");
const BASE_URL = 'http://127.0.0.1:8080';

function userinfo() {
  const userinfo = sessionStorage.getItem("userinfo");
  console.log("this is the userinfo:", userinfo, window.location);
  return JSON.parse(userinfo);
}

function redirectIfInvalidSession() {
  const redirectPath = "/sign/sign.html";
  const user = userinfo();
  if (!user) {
    console.log("no user found in session");
    window.location.replace(redirectPath);
  }

  USERNAME_SPAN.innerHTML = user.name;
  isLoggedIn()
    .then(
      ok => {
        if (!ok) {
          console.log("isLoggedIn(): not ok!");
          window.location.replace(redirectPath);
        }
      },
      err => {
        console.log("isLoggedIn(): error!", err);
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
    window.location.replace('/');
    console.log("this is the replace:", window.location);
  } else {
    console.log("error loggin out");
  }
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


redirectIfInvalidSession();
