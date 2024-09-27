const BASE_URL = 'http://127.0.0.1:8080';
const LOGIN_FORM = document.getElementById('login-form');

async function login(event) {
  event.preventDefault();

  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  const response = await fetch(`${BASE_URL}/signin`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': 'index.html',
    },
    body: JSON.stringify({
      'name': username,
      'password': password,
    })
  });

  console.log(response);

  if (response.ok) {
    /*
    const token = await response.text;
    document.cookie = `token=${token}; path=/; secure`;
    console.log(`token: ${token}`);
    localStorage.setItem('token', token);
    */

    // cookie should be set here, check in index if the cookie is set
    // TODO redirect to index.html
    window.location.replace('../index.html');
  } else {
    const errorText = document.getElementById('error-text');
    errorText.textContent = 'Login failed. Please check your credentials.';
  }
}
