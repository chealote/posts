const USERNAME_SPAN = document.getElementById("username-span");
const BASE_URL = "http://127.0.0.1:8080";
const BODIES = document.getElementsByTagName("body");
if (BODIES.length == 0) {
  console.error("No body found");
  // TODO throw exception or something?
}
const BODY = BODIES[0];
console.log("bodies:", BODIES, "and body:", BODY);

function authRedirect(path) {
  BODY.style.display = "none";
  window.location.replace(path);
}

function userinfo() {
  const userinfo = sessionStorage.getItem("userinfo");
  console.log("this is the userinfo:", userinfo, window.location);
  return JSON.parse(userinfo);
}

function redirectIfInvalidSession() {
  // TODO remove this, not working as expected, there should be
  // like a middleware that renders only if auth allows it
  BODY.style.display = "none";

  const redirectPath = "/sign/sign.html";
  const user = userinfo();
  if (!user) {
    console.log("no user found in session");
    authRedirect(redirectPath);
  }

  isLoggedIn()
    .then(
      ok => {
        if (!ok) {
          console.log("isLoggedIn(): not ok!");
          authRedirect(redirectPath);
        } else {
          setTimeout(function() {
            fillHeader(user);
            BODY.style.display = "block";
          }, 300);
        }
      },
      err => {
        console.log("isLoggedIn(): error!", err);
        authRedirect(redirectPath);
      });
}

async function logout() {
  console.log("trying to logout...");
  const response = await fetch(`${BASE_URL}/logout`, {
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    }
  });

  if (response.ok) {
    authRedirect("/");
  } else {
    console.error("error loggin out");
  }
}

async function isLoggedIn() {
  const response = await fetch(`${BASE_URL}/token`, {
    method: "GET",
    headers: {
      "Authorization": sessionStorage.getItem("token"),
    },
  });

  return new Promise((resolve, reject) => {
    return resolve(response.ok);
  });
}

function createNewPost() {
  authRedirect("/posts/create.html");
}

async function fillHeader(user) {
  header.innerHTML = `(${user.name})
  <button onclick="logout()">logout</button>
  <button onclick="createNewPost()">new post</button>
`;
}

redirectIfInvalidSession();
